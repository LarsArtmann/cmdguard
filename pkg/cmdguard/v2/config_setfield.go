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
		return err
	}

	val := reflect.ValueOf(value)

	// Handle type conversions
	if val.Type().ConvertibleTo(field.Type()) {
		field.Set(val.Convert(field.Type()))
		return nil
	}

	// Handle string to custom type conversions
	if val.Kind() == reflect.String {
		if err := setStringField(field, val.String()); err != nil {
			return err
		}
		return nil
	}

	// Handle time.Duration to Duration conversion
	if val.Type() == reflect.TypeOf(time.Duration(0)) && field.Type() == reflect.TypeOf(Duration{}) {
		field.Set(reflect.ValueOf(FromDuration(val.Interface().(time.Duration))))
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", value, field.Type())
}

// getField retrieves a field from config by name.
func getField(cfg any, fieldName string) (reflect.Value, error) {
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("config must be a pointer to struct")
	}

	v = v.Elem()
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return reflect.Value{}, fmt.Errorf("field %q not found", fieldName)
	}

	if !field.CanSet() {
		return reflect.Value{}, fmt.Errorf("field %q is not settable", fieldName)
	}

	return field, nil
}

// setStringField handles string to custom type conversions.
func setStringField(field reflect.Value, str string) error {
	switch field.Type() {
	case reflect.TypeOf(LogLevel{}):
		return parseAndSetLogLevel(field, str)
	case reflect.TypeOf(LogFormat{}):
		return parseAndSetLogFormat(field, str)
	case reflect.TypeOf(Duration{}):
		return parseAndSetDuration(field, str)
	case reflect.TypeOf(Enum{}):
		field.Set(reflect.ValueOf(Enum{value: str}))
		return nil
	}
	return fmt.Errorf("unsupported string conversion for %s", field.Type())
}

// parseAndSetLogLevel parses and sets a LogLevel field.
func parseAndSetLogLevel(field reflect.Value, str string) error {
	parsed, err := ParseLogLevel(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
}

// parseAndSetLogFormat parses and sets a LogFormat field.
func parseAndSetLogFormat(field reflect.Value, str string) error {
	parsed, err := ParseLogFormat(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
}

// parseAndSetDuration parses and sets a Duration field.
func parseAndSetDuration(field reflect.Value, str string) error {
	parsed, err := ParseDuration(str)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(parsed))
	return nil
}
