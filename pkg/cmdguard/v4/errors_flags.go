package v4

import "errors"

// Flag-related sentinel errors.
var (
	// ErrFlagParseFailed indicates flag parsing failed.
	ErrFlagParseFailed = errors.New("failed to parse flags")

	// ErrInvalidFlagType indicates the flag type F is not a struct or pointer to struct.
	ErrInvalidFlagType = errors.New("invalid flag type")

	// ErrFlagNotFound indicates a flag was not found in the command.
	ErrFlagNotFound = errors.New("flag not found")

	// ErrRequiredFlag indicates a required flag was not set.
	ErrRequiredFlag = errors.New("required flag not set")

	// ErrFlagInstance indicates a flag struct instance could not be created.
	ErrFlagInstance = errors.New("failed to create flag instance")

	// ErrInvalidEnum indicates an invalid enum value.
	ErrInvalidEnum = errors.New("invalid enum value")

	// ErrValueTooShort indicates a string value is shorter than required.
	ErrValueTooShort = errors.New("value too short")

	// ErrValueTooLong indicates a string value is longer than allowed.
	ErrValueTooLong = errors.New("value too long")

	// ErrValueTooSmall indicates a numeric value is smaller than minimum.
	ErrValueTooSmall = errors.New("value too small")

	// ErrValueTooLarge indicates a numeric value is larger than maximum.
	ErrValueTooLarge = errors.New("value too large")

	// ErrValueOutOfRange indicates a numeric value does not fit the field's
	// integer bit-width (e.g. 999 written into an int8).
	ErrValueOutOfRange = errors.New("value out of range")

	// ErrValuePatternMismatch indicates a value does not match the required pattern.
	ErrValuePatternMismatch = errors.New("value does not match pattern")

	// ErrValueEmpty indicates a value is empty but is required to be non-empty.
	ErrValueEmpty = errors.New("value is empty")

	// ErrUnknownValidator indicates a validator name in a validate tag is not registered.
	ErrUnknownValidator = errors.New("unknown validator")

	// ErrInvalidValidatorParam indicates a validator parameter could not be parsed.
	ErrInvalidValidatorParam = errors.New("invalid validator parameter")

	// ErrFieldNotFound indicates a struct field was not found by name.
	ErrFieldNotFound = errors.New("field not found")

	// ErrFieldNotSettable indicates a struct field is unexported or otherwise not settable.
	ErrFieldNotSettable = errors.New("field is not settable")

	// ErrTypeConversion indicates a type conversion or assertion failed.
	ErrTypeConversion = errors.New("type conversion failed")

	// ErrUnsupportedConversion indicates a string-to-type conversion is not supported.
	ErrUnsupportedConversion = errors.New("unsupported string conversion")

	// ErrNilValue indicates a required value was nil.
	ErrNilValue = errors.New("value must not be nil")
)
