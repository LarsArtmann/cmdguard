package v2

import (
	"maps"
	"reflect"
	"sync"

	"github.com/spf13/pflag"
)

// TypeHandler defines how a flag type is registered, parsed, defaulted, and set.
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
}

// globalTypeRegistry is the default registry with all built-in types.
var globalTypeRegistry = newTypeRegistry()

func newTypeRegistry() *typeRegistry {
	r := &typeRegistry{
		byType: make(map[reflect.Type]TypeHandler),
		byKind: make(map[reflect.Kind]TypeHandler),
	}

	r.registerKinds()
	r.registerCustomTypes()

	return r
}

// lookupHandler finds the TypeHandler for a given reflect.Type.
func (r *typeRegistry) lookupHandler(typ reflect.Type) (TypeHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if h, ok := r.byType[typ]; ok {
		return h, true
	}

	if h, ok := r.byKind[typ.Kind()]; ok {
		return h, true
	}

	return nil, false
}

// clone returns an independent copy of the typeRegistry.
// The maps are shallow-copied; TypeHandler values are shared (they are stateless).
func (r *typeRegistry) clone() *typeRegistry {
	r.mu.RLock()
	defer r.mu.RUnlock()

	c := &typeRegistry{
		byType:       make(map[reflect.Type]TypeHandler, len(r.byType)),
		byKind:       make(map[reflect.Kind]TypeHandler, len(r.byKind)),
		countHandler: r.countHandler,
	}

	maps.Copy(c.byType, r.byType)
	maps.Copy(c.byKind, r.byKind)

	return c
}

// register adds or replaces a TypeHandler for a specific reflect.Type.
func (r *typeRegistry) register(typ reflect.Type, handler TypeHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.byType[typ] = handler
}

// RegisterTypeHandler registers a custom TypeHandler for a specific reflect.Type.
// This writes to the global defaults template; new FlagRegistries will include
// this handler. For per-instance registration, use FlagRegistry.RegisterTypeHandler.
func RegisterTypeHandler(typ reflect.Type, handler TypeHandler) {
	globalTypeRegistry.register(typ, handler)
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
			return tr.countHandler.Register(flags, tag)
		}
	}

	h, ok := tr.lookupHandler(tag.Type)
	if !ok {
		registerStringFlag(flags, tag.Name, tag.Short, tag.Default, tag.Help)

		return nil
	}

	return h.Register(flags, tag)
}

// dispatchParse dispatches value parsing to the TypeHandler registry.
func dispatchParse(tr *typeRegistry, value string, tag FlagTag) (any, error) {
	h, ok := tr.lookupHandler(tag.Type)
	if !ok {
		return value, nil
	}

	return h.Parse(value, tag)
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
