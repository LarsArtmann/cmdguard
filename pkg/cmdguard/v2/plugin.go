package v2

import (
	"fmt"
	"reflect"
)

// Plugin bundles custom type handlers and validators for one-step registration.
// Implementations register their extensions in Register, returning an error on
// failure (consistent with the zero-panic guarantee).
type Plugin interface {
	Name() string
	Register(r PluginRegistrar) error
}

// PluginRegistrar is the registration target passed to Plugin.Register.
// It exposes the same type-handler and validator registration that the global
// and per-instance APIs use, but scoped to a single plugin application.
type PluginRegistrar struct {
	types      *typeRegistry
	validators *validatorRegistry
}

// TypeHandler registers a custom type handler for a specific reflect.Type.
func (r PluginRegistrar) TypeHandler(typ reflect.Type, handler TypeHandler) {
	r.types.register(typ, handler)
}

// TypeHandlerFunc is a convenience wrapper for registering a TypeHandlerFunc.
func (r PluginRegistrar) TypeHandlerFunc(typ reflect.Type, handler TypeHandlerFunc) {
	r.types.register(typ, handler)
}

// Validator registers a named flag validator.
func (r PluginRegistrar) Validator(name string, fn FlagValidator) {
	r.validators.register(name, fn)
}

// RegisterPlugin applies a plugin to the global registries.
// This makes the plugin's handlers and validators available to all new
// FlagRegistries (copy-on-write). Returns an error if the plugin's Register
// fails.
func RegisterPlugin(plugin Plugin) error {
	if plugin == nil {
		return fmt.Errorf("%w: plugin is nil", ErrServiceRegistration)
	}

	return plugin.Register(PluginRegistrar{
		types:      globalTypeRegistry,
		validators: globalValidators,
	})
}

// WithPlugin registers a plugin during CLI construction.
// The plugin is applied to the global registries before the CLI initializes its
// flags, so the plugin's handlers and validators are available to all commands.
func WithPlugin[T any](plugin Plugin) CLIOption[T] {
	return func(cli *CLI[T]) {
		_ = RegisterPlugin(plugin)
	}
}
