package v4

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

	return applyFieldValue(field, value, tr)
}

// setFieldByIndex sets a field resolved via its reflect index path (nested structs).
func setFieldByIndex(cfg any, index []int, value any, tr *typeRegistry) error {
	field, err := getFieldByIndex(cfg, index)
	if err != nil {
		return fmt.Errorf("SetField: cfg=%T, index=%v, value=%T: %w", cfg, index, value, err)
	}

	return applyFieldValue(field, value, tr)
}

// setFieldByTag resolves the field via tag.Index (nested) or tag.Field (flat) and sets it.
func setFieldByTag(cfg any, tag FlagTag, value any, tr *typeRegistry) error {
	if len(tag.Index) > 0 {
		return setFieldByIndex(cfg, tag.Index, value, tr)
	}

	return setField(cfg, tag.Field, value, tr)
}

// fieldByTag resolves a reflect.Value for the field described by tag.
// Uses Index (nested path) when available, falling back to FieldByName.
func fieldByTag(v reflect.Value, tag FlagTag) reflect.Value {
	if len(tag.Index) > 0 {
		return v.FieldByIndex(tag.Index)
	}

	return v.FieldByName(tag.Field)
}

// applyFieldValue performs the type-conversion logic to set value onto field.
func applyFieldValue(field reflect.Value, value any, tr *typeRegistry) error {
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
				"SetField: value=%q, field=%s: %w",
				value,
				field.Type(),
				err,
			)
		}

		return nil
	}

	// Handle time.Duration to Duration conversion
	if val.Type() == reflect.TypeFor[time.Duration]() &&
		field.Type() == reflect.TypeFor[Duration]() {
		duration, ok := reflect.TypeAssert[time.Duration](val)
		if !ok {
			return fmt.Errorf(
				"SetField: type assertion failed for time.Duration, value=%v, field=%s: %w",
				value,
				field.Type(),
				ErrTypeConversion,
			)
		}

		field.Set(reflect.ValueOf(FromDuration(duration)))

		return nil
	}

	return fmt.Errorf(
		"SetField: cannot convert value=%T to %s: %w",
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

// getFieldByIndex resolves a field via its reflect index path (supports nested structs).
func getFieldByIndex(cfg any, index []int) (reflect.Value, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf(
			"getFieldByIndex: cfg=%T, index=%v: %w",
			cfg,
			index,
			ErrConfigNotPointer,
		)
	}

	field := v.Elem().FieldByIndex(index)
	if !field.CanSet() {
		return reflect.Value{}, fmt.Errorf(
			"getFieldByIndex: cfg=%T, index=%v: %w",
			cfg,
			index,
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
		current, ok := reflect.TypeAssert[Enum](field)
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

		if !parsedVal.Type().AssignableTo(field.Type()) {
			return fmt.Errorf("setStringField: field=%s, str=%q, parsed=%s: %w",
				field.Type(), str, parsedVal.Type(), ErrUnsupportedConversion)
		}

		field.Set(parsedVal)

		return nil
	}

	return fmt.Errorf("setStringField: field=%s, str=%q: %w",
		field.Type(), str, ErrUnsupportedConversion)
}
