package v3

import (
	"fmt"
	"maps"
	"reflect"
	"sync"

	"github.com/spf13/pflag"
)

// TypeHandler defines how a flag type is registered, parsed, defaulted, and set.
// TODO(v4): consider renaming to TypeCodec for precision (see docs/reviews/2026-07-18_09-44_naming-review.html).
type TypeHandler interface {
	Register(flags *pflag.FlagSet, tag FlagTag) error
	Parse(value string, tag FlagTag) (any, error)
	Default(tag FlagTag) any
}

// registerStringFlag registers a string flag with optional shorthand.
func registerStringFlag(flags *pflag.FlagSet, name, short, value, usage string) {
	if short != "" {
		_ = flags.StringP(name, short, value, usage)
	} else {
		_ = flags.String(name, value, usage)
	}
}

// registerStringFlagFromTag registers a string flag using FlagTag fields and returns nil.
func registerStringFlagFromTag(flags *pflag.FlagSet, tag FlagTag) error {
	registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)

	return nil
}

// TypeHandlerFunc is a functional adapter for TypeHandler where Register is not needed.
type TypeHandlerFunc struct {
	ParseFunc    func(value string, tag FlagTag) (any, error)
	DefaultFunc  func(tag FlagTag) any
	RegisterFunc func(flags *pflag.FlagSet, tag FlagTag) error
}

func (h TypeHandlerFunc) Register(flags *pflag.FlagSet, tag FlagTag) error {
	if h.RegisterFunc != nil {
		return h.RegisterFunc(flags, tag)
	}

	return nil
}

func (h TypeHandlerFunc) Parse(value string, tag FlagTag) (any, error) {
	return h.ParseFunc(value, tag)
}

func (h TypeHandlerFunc) Default(tag FlagTag) any {
	return h.DefaultFunc(tag)
}

// typeRegistry maps reflect.Type or reflect.Kind to TypeHandlers.
type typeRegistry struct {
	mu           sync.RWMutex
	byType       map[reflect.Type]TypeHandler
	byKind       map[reflect.Kind]TypeHandler
	countHandler TypeHandler
	owned        bool          // true if this instance owns its maps (can mutate)
	parent       *typeRegistry // nil when owned; shared source for COW reads
}

// globalTypeRegistry is the default registry with all built-in types.
var globalTypeRegistry = newTypeRegistry()

func newTypeRegistry() *typeRegistry {
	r := &typeRegistry{
		byType: make(map[reflect.Type]TypeHandler),
		byKind: make(map[reflect.Kind]TypeHandler),
		owned:  true,
	}

	r.registerKinds()
	r.registerCustomTypes()

	return r
}

// lookupHandler finds the TypeHandler for a given reflect.Type.
// For copy-on-write instances (owned=false), delegates to the parent registry.
func (r *typeRegistry) lookupHandler(typ reflect.Type) (TypeHandler, bool) {
	r.mu.RLock()

	if r.owned {
		defer r.mu.RUnlock()

		if h, ok := r.byType[typ]; ok {
			return h, true
		}

		if h, ok := r.byKind[typ.Kind()]; ok {
			return h, true
		}

		return nil, false
	}

	parent := r.parent
	r.mu.RUnlock()

	return parent.lookupHandler(typ)
}

// share returns a copy-on-write view of this registry.
// The returned instance reads from this registry's maps until the first write,
// at which point it clones lazily. This avoids the clone cost for the common
// case where no per-instance customization is used.
func (r *typeRegistry) share() *typeRegistry {
	root := r
	if !r.owned && r.parent != nil {
		root = r.parent
	}

	root.mu.RLock()
	defer root.mu.RUnlock()

	return &typeRegistry{
		countHandler: root.countHandler,
		owned:        false,
		parent:       root,
	}
}

// register adds or replaces a TypeHandler for a specific reflect.Type.
// For copy-on-write instances, triggers a lazy clone on first write.
func (r *typeRegistry) register(typ reflect.Type, handler TypeHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.owned && r.parent != nil {
		r.parent.mu.RLock()
		r.byType = maps.Clone(r.parent.byType)
		r.byKind = maps.Clone(r.parent.byKind)
		r.parent.mu.RUnlock()
		r.parent = nil
		r.owned = true
	}

	r.byType[typ] = handler
}

// RegisterTypeHandler registers a custom TypeHandler for a specific reflect.Type.
// This writes to the global defaults template; new FlagRegistries will include
// this handler. For per-instance registration, use FlagRegistry.RegisterTypeHandler.
// Returns an error if typ or handler is nil.
func RegisterTypeHandler(typ reflect.Type, handler TypeHandler) error {
	if typ == nil {
		return fmt.Errorf("%w: typ is nil", ErrServiceRegistration)
	}

	if handler == nil {
		return fmt.Errorf("%w: handler is nil for type %v", ErrServiceRegistration, typ)
	}

	globalTypeRegistry.register(typ, handler)

	return nil
}

// handledByTypeRegistry checks whether the given type has a handler in the given registry.
func handledByTypeRegistry(tr *typeRegistry, typ reflect.Type) bool {
	_, ok := tr.lookupHandler(typ)

	return ok
}

// dispatchRegister dispatches flag registration to the TypeHandler registry.
func dispatchRegister(tr *typeRegistry, flags *pflag.FlagSet, tag FlagTag) error {
	if tag.Count {
		if tr.countHandler != nil {
			err := tr.countHandler.Register(flags, tag)
			if err != nil {
				return fmt.Errorf("registering count flag %q: %w", tag.Name, err)
			}

			return nil
		}
	}

	h, ok := tr.lookupHandler(tag.Type)
	if !ok {
		registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)

		return nil
	}

	err := h.Register(flags, tag)
	if err != nil {
		return fmt.Errorf("registering flag %q: %w", tag.Name, err)
	}

	return nil
}

// dispatchParse dispatches value parsing to the TypeHandler registry.
// Returns an error if the parsed value type doesn't match the target field type.
func dispatchParse(tr *typeRegistry, value string, tag FlagTag) (any, error) {
	h, ok := tr.lookupHandler(tag.Type)
	if !ok {
		return value, nil
	}

	parsed, err := h.Parse(value, tag)
	if err != nil {
		return nil, fmt.Errorf("value=%q, field=%q: %w", value, tag.Field, err)
	}

	parsedVal := reflect.ValueOf(parsed)
	if parsedVal.IsValid() && !parsedVal.Type().AssignableTo(tag.Type) &&
		!parsedVal.Type().ConvertibleTo(tag.Type) {
		return nil, fmt.Errorf("value=%q: dispatchParse: type handler returned %s, field %q requires %s: %w",
			value, parsedVal.Type(), tag.Field, tag.Type, ErrUnsupportedConversion)
	}

	return parsed, nil
}

// dispatchDefault dispatches default value computation to the TypeHandler registry.
func dispatchDefault(tr *typeRegistry, tag FlagTag) any {
	if tag.Default == "" {
		return reflect.Zero(tag.Type).Interface()
	}

	h, ok := tr.lookupHandler(tag.Type)
	if !ok {
		return tag.Default
	}

	def := h.Default(tag)

	defVal := reflect.ValueOf(def)
	if defVal.IsValid() && defVal.Type().ConvertibleTo(tag.Type) {
		return defVal.Convert(tag.Type).Interface()
	}

	return def
}
