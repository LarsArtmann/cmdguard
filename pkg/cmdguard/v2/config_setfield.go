package v2

import (
	"fmt"
	"reflect"
	"time"
)

// SetField sets a field value on a config struct using reflection.
// This is the backward-compatible wrapper that uses the global type registry.
func SetField(cfg any, fieldName string, value any) error {
	return setField(cfg, fieldName, value, globalTypeRegistry)
}

// setField sets a field value on a config struct using reflection with the given type registry.
func setField(cfg any, fieldName string, value any, tr *typeRegistry) error {
	field, err := getField(cfg, fieldName)
	if err != nil {
		return fmt.Errorf(
			"SetField: cfg=%T, fieldName=%q, value=%T: %w",
			cfg,
			fieldName,
			value,
			err,
		)
	}

	val := reflect.ValueOf(value)

	// Handle type conversions
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))

		return nil
	}

	// Handle string to custom type conversions
	if val.Kind() == reflect.String {
		err := setStringField(field, val.String(), tr)
		if err != nil {
			return fmt.Errorf(
				"SetField: cfg=%T, fieldName=%q, value=%q: %w",
				cfg,
				fieldName,
				value,
				err,
			)
		}

		return nil
	}

	// Handle time.Duration to Duration conversion
	if val.Type() == reflect.TypeFor[time.Duration]() &&
		field.Type() == reflect.TypeFor[Duration]() {
		duration, ok := val.Interface().(time.Duration)
		if !ok {
			return fmt.Errorf(
				"SetField: type assertion failed for time.Duration, cfg=%T, fieldName=%q, value=%v: %w",
				cfg,
				fieldName,
				value,
				ErrTypeConversion,
			)
		}

		field.Set(reflect.ValueOf(FromDuration(duration)))

		return nil
	}

	return fmt.Errorf(
		"SetField: cannot convert cfg=%T, fieldName=%q, value=%T to %s: %w",
		cfg,
		fieldName,
		value,
		field.Type(),
		ErrTypeConversion,
	)
}

// getField retrieves a field from config by name.
func getField(cfg any, fieldName string) (reflect.Value, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf(
			"getField: cfg=%T, fieldName=%q: %w",
			cfg,
			fieldName,
			ErrConfigNotPointer,
		)
	}

	v = v.Elem()

	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf(
			"getField: cfg=%T, fieldName=%q: %w",
			cfg,
			fieldName,
			ErrFieldNotFound,
		)
	}

	if !field.CanSet() {
		return reflect.Value{}, fmt.Errorf(
			"getField: cfg=%T, fieldName=%q: %w",
			cfg,
			fieldName,
			ErrFieldNotSettable,
		)
	}

	return field, nil
}

// setStringFieldError wraps an error with field and string context.
func setStringFieldError(field reflect.Value, str string, err error) error {
	return fmt.Errorf("setStringField: field=%s, str=%q: %w", field.Type(), str, err)
}

// setStringField handles string to custom type conversions via the TypeHandler registry.
func setStringField(field reflect.Value, str string, tr *typeRegistry) error {
	// Special handling for Enum (needs current allowed values from the field)
	if field.Type() == reflect.TypeFor[Enum]() {
		current, ok := field.Interface().(Enum)
		if !ok {
			return NewConfigError(
				field.Type().String(),
				ErrInvalidFlagType,
			)
		}

		allowed := current.Allowed()

		if len(allowed) == 0 {
			field.Set(reflect.ValueOf(Enum{value: str}))

			return nil
		}

		parsed, err := ParseEnum(str, allowed)
		if err != nil {
			return setStringFieldError(field, str, err)
		}

		field.Set(reflect.ValueOf(parsed))

		return nil
	}

	// Try the TypeHandler registry for all other types
	if handledByTypeRegistry(tr, field.Type()) {
		parsed, err := dispatchParse(tr, str, FlagTag{Type: field.Type()})
		if err != nil {
			return setStringFieldError(field, str, err)
		}

		parsedVal := reflect.ValueOf(parsed)
		if parsedVal.Type().ConvertibleTo(field.Type()) {
			field.Set(parsedVal.Convert(field.Type()))

			return nil
		}

		field.Set(parsedVal)

		return nil
	}

	return fmt.Errorf("setStringField: field=%s, str=%q: %w",
		field.Type(), str, ErrUnsupportedConversion)
}
