package v2

import (
	"fmt"
	"reflect"
	"time"
)

// SetField sets a field value on a config struct using reflection.
func SetField(cfg any, fieldName string, value any) error {
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
		err := setStringField(field, val.String())
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

// setStringField handles string to custom type conversions.
func setStringField(field reflect.Value, str string) error {
	switch field.Type() {
	case reflect.TypeFor[LogLevel]():
		return wrapErr(parseAndSetLogLevel(field, str), field, str)
	case reflect.TypeFor[LogFormat]():
		return wrapErr(parseAndSetLogFormat(field, str), field, str)
	case reflect.TypeFor[Duration]():
		return wrapErr(parseAndSetDuration(field, str), field, str)
	case reflect.TypeFor[Enum]():
		current := field.Interface().(Enum)
		allowed := current.Allowed()

		if len(allowed) == 0 {
			field.Set(reflect.ValueOf(Enum{value: str}))

			return nil
		}

		parsed, err := ParseEnum(str, allowed)
		if err != nil {
			return fmt.Errorf("setStringField: field=%s, str=%q: %w", field.Type(), str, err)
		}

		field.Set(reflect.ValueOf(parsed))

		return nil
	}

	return fmt.Errorf("setStringField: field=%s, str=%q: %w",
		field.Type(), str, ErrUnsupportedConversion)
}

// wrapErr adds context to a setStringField error.
func wrapErr(err error, field reflect.Value, str string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("setStringField: field=%s, str=%q: %w", field.Type(), str, err)
}

// parseField is a helper for parsing string values into fields.
func parseField[T any](
	field reflect.Value,
	str string,
	typeName string,
	parser func(string) (T, error),
) error {
	parsed, err := parser(str)
	if err != nil {
		return fmt.Errorf("parsing %s field %s with %q: %w", typeName, field.Type(), str, err)
	}

	field.Set(reflect.ValueOf(parsed))

	return nil
}

// parseAndSetLogLevel parses and sets a LogLevel field.
func parseAndSetLogLevel(field reflect.Value, str string) error {
	return parseField(field, str, "log level", ParseLogLevel)
}

// parseAndSetLogFormat parses and sets a LogFormat field.
func parseAndSetLogFormat(field reflect.Value, str string) error {
	return parseField(field, str, "log format", ParseLogFormat)
}

// parseAndSetDuration parses and sets a Duration field.
func parseAndSetDuration(field reflect.Value, str string) error {
	return parseField(field, str, "duration", ParseDuration)
}
