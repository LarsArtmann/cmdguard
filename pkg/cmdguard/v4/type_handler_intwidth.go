package v4

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/pflag"
)

// Bit-widths for the fixed-size integer kinds. bitSizePlatformWidth (0) means
// "platform int/uint" and is what strconv expects for reflect.Int/reflect.Uint.
const (
	bitSizeInt8  = 8
	bitSizeInt16 = 16
	bitSizeInt32 = 32
	bitSizeInt64 = 64
	// bitSizePlatformWidth tells strconv to use the platform int width.
	bitSizePlatformWidth = 0
)

// intBitSize returns the strconv bit-size for a signed integer kind.
func intBitSize(k reflect.Kind) int {
	switch k { //nolint:exhaustive // non-integer kinds fall through to default
	case reflect.Int8:
		return bitSizeInt8
	case reflect.Int16:
		return bitSizeInt16
	case reflect.Int32:
		return bitSizeInt32
	case reflect.Int64:
		return bitSizeInt64
	default:
		return bitSizePlatformWidth
	}
}

// uintBitSize returns the strconv bit-size for an unsigned integer kind.
func uintBitSize(k reflect.Kind) int {
	switch k { //nolint:exhaustive // non-integer kinds fall through to default
	case reflect.Uint8:
		return bitSizeInt8
	case reflect.Uint16:
		return bitSizeInt16
	case reflect.Uint32:
		return bitSizeInt32
	case reflect.Uint64:
		return bitSizeInt64
	default:
		return bitSizePlatformWidth
	}
}

// makeIntKindHandler builds a TypeHandlerFunc for a signed integer kind that
// validates flag values against the field's actual bit-width. Values outside
// the representable range produce an error instead of silently wrapping.
// Previously every signed kind shared one handler parsing with bitSize 64, so a
// value like 999 written into an int8 silently wrapped to -25.
func makeIntKindHandler(bitSize int) TypeHandlerFunc { //nolint:dupl // mirrors makeUintKindHandler (signed vs unsigned)
	return TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseIntDefault(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid int default for flag %q: %w", tag.Name, err)
			}

			err = checkIntRange(def, bitSize)
			if err != nil {
				return fmt.Errorf("int default for flag %q: %w", tag.Name, err)
			}

			if tag.Short != "" {
				flags.IntP(tag.Name, tag.Short, int(def), tag.Help)
			} else {
				flags.Int(tag.Name, int(def), tag.Help)
			}

			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strconv.ParseInt(value, 10, bitSize)
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseIntDefault(tag.Default)

			return v
		},
	}
}

// makeUintKindHandler builds a TypeHandlerFunc for an unsigned integer kind that
// validates flag values against the field's actual bit-width.
func makeUintKindHandler(bitSize int) TypeHandlerFunc { //nolint:dupl // mirrors makeIntKindHandler (signed vs unsigned)
	return TypeHandlerFunc{
		RegisterFunc: func(flags *pflag.FlagSet, tag FlagTag) error {
			def, err := parseUintDefault(tag.Default)
			if err != nil {
				return fmt.Errorf("invalid uint default for flag %q: %w", tag.Name, err)
			}

			err = checkUintRange(def, bitSize)
			if err != nil {
				return fmt.Errorf("uint default for flag %q: %w", tag.Name, err)
			}

			if tag.Short != "" {
				flags.UintP(tag.Name, tag.Short, uint(def), tag.Help)
			} else {
				flags.Uint(tag.Name, uint(def), tag.Help)
			}

			return nil
		},
		ParseFunc: func(value string, _ FlagTag) (any, error) {
			return strconv.ParseUint(value, 10, bitSize)
		},
		DefaultFunc: func(tag FlagTag) any {
			v, _ := parseUintDefault(tag.Default)

			return v
		},
	}
}

// checkIntRange reports an error if v does not fit in a signed integer of the
// given bit-width. bitSizePlatformWidth means the platform int width.
func checkIntRange(v int64, bitSize int) error {
	size := bitSize
	if size == bitSizePlatformWidth {
		size = strconv.IntSize
	}

	half := int64(1) << uint(size-1)
	if v < -half || v > half-1 {
		return fmt.Errorf("%w: %d does not fit in int%d", ErrValueOutOfRange, v, size)
	}

	return nil
}

// checkUintRange reports an error if v does not fit in an unsigned integer of
// the given bit-width. bitSizePlatformWidth means the platform uint width.
func checkUintRange(v uint64, bitSize int) error {
	size := bitSize
	if size == bitSizePlatformWidth {
		size = strconv.IntSize
	}

	// At 64 bits every uint64 fits; shifting by 64 would yield 0 and misfire.
	if size >= bitSizeInt64 {
		return nil
	}

	cutoff := uint64(1) << uint(size)
	if v >= cutoff {
		return fmt.Errorf("%w: %d does not fit in uint%d", ErrValueOutOfRange, v, size)
	}

	return nil
}
