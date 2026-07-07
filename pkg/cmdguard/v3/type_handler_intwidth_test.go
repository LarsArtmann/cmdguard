package v3

import (
	"reflect"
	"testing"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

func TestIntKindHandler_ParseInRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		typ   reflect.Type
		value string
		want  any
	}{
		{"int8 max", reflect.TypeFor[int8](), "127", int64(127)},
		{"int8 min", reflect.TypeFor[int8](), "-128", int64(-128)},
		{"int16 max", reflect.TypeFor[int16](), "32767", int64(32767)},
		{"int32 max", reflect.TypeFor[int32](), "2147483647", int64(2147483647)},
		{"int64 large", reflect.TypeFor[int64](), "9000000000", int64(9000000000)},
		{"int platform", reflect.TypeFor[int](), "42", int64(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := FlagTag{Type: tt.typ}
			result, err := dispatchParse(globalTypeRegistry, tt.value, tag)
			testutil.AssertNoError(t, err)
			if result != tt.want {
				t.Errorf("dispatchParse(%s, %s) = %v (%T), want %v", tt.typ, tt.value, result, result, tt.want)
			}
		})
	}
}

func TestIntKindHandler_ParseOverflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		typ   reflect.Type
		value string
	}{
		{"int8 overflow positive", reflect.TypeFor[int8](), "128"},
		{"int8 overflow negative", reflect.TypeFor[int8](), "-129"},
		{"int8 big", reflect.TypeFor[int8](), "999"},
		{"int16 overflow", reflect.TypeFor[int16](), "32768"},
		{"int32 overflow", reflect.TypeFor[int32](), "2147483648"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := FlagTag{Type: tt.typ}
			_, err := dispatchParse(globalTypeRegistry, tt.value, tag)
			testutil.AssertExpectedError(t, err)
		})
	}
}

func TestUintKindHandler_ParseInRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		typ   reflect.Type
		value string
		want  any
	}{
		{"uint8 max", reflect.TypeFor[uint8](), "255", uint64(255)},
		{"uint16 max", reflect.TypeFor[uint16](), "65535", uint64(65535)},
		{"uint32 max", reflect.TypeFor[uint32](), "4294967295", uint64(4294967295)},
		{"uint64 large", reflect.TypeFor[uint64](), "18000000000", uint64(18000000000)},
		{"uint platform", reflect.TypeFor[uint](), "42", uint64(42)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := FlagTag{Type: tt.typ}
			result, err := dispatchParse(globalTypeRegistry, tt.value, tag)
			testutil.AssertNoError(t, err)
			if result != tt.want {
				t.Errorf("dispatchParse(%s, %s) = %v (%T), want %v", tt.typ, tt.value, result, result, tt.want)
			}
		})
	}
}

func TestUintKindHandler_ParseOverflow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		typ   reflect.Type
		value string
	}{
		{"uint8 overflow", reflect.TypeFor[uint8](), "256"},
		{"uint8 big", reflect.TypeFor[uint8](), "999"},
		{"uint16 overflow", reflect.TypeFor[uint16](), "65536"},
		{"uint negative", reflect.TypeFor[uint8](), "-1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tag := FlagTag{Type: tt.typ}
			_, err := dispatchParse(globalTypeRegistry, tt.value, tag)
			testutil.AssertExpectedError(t, err)
		})
	}
}

func TestCheckIntRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   int64
		bitSize int
		wantErr bool
	}{
		{"int8 in range", 100, 8, false},
		{"int8 at max", 127, 8, false},
		{"int8 at min", -128, 8, false},
		{"int8 over max", 128, 8, true},
		{"int8 under min", -129, 8, true},
		{"int16 in range", 30000, 16, false},
		{"int16 over", 32768, 16, true},
		{"int32 over", 2147483648, 32, true},
		{"int64 huge", 9000000000, 64, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkIntRange(tt.value, tt.bitSize)
			if tt.wantErr {
				testutil.AssertErrorIs(t, err, ErrValueOutOfRange)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

func TestCheckUintRange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   uint64
		bitSize int
		wantErr bool
	}{
		{"uint8 in range", 200, 8, false},
		{"uint8 at max", 255, 8, false},
		{"uint8 over", 256, 8, true},
		{"uint16 over", 65536, 16, true},
		{"uint32 over", 4294967296, 32, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := checkUintRange(tt.value, tt.bitSize)
			if tt.wantErr {
				testutil.AssertErrorIs(t, err, ErrValueOutOfRange)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

func TestIntBitSizeMapping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		kind reflect.Kind
		want int
	}{
		{reflect.Int8, 8},
		{reflect.Int16, 16},
		{reflect.Int32, 32},
		{reflect.Int64, 64},
		{reflect.Int, 0},
		{reflect.Uint8, 8},
		{reflect.Uint16, 16},
		{reflect.Uint32, 32},
		{reflect.Uint64, 64},
		{reflect.Uint, 0},
	}

	for _, tt := range tests {
		if tt.kind >= reflect.Uint && tt.kind <= reflect.Uintptr {
			continue // handled by uintBitSize below
		}

		if got := intBitSize(tt.kind); got != tt.want {
			t.Errorf("intBitSize(%s) = %d, want %d", tt.kind, got, tt.want)
		}
	}

	uintTests := []struct {
		kind reflect.Kind
		want int
	}{
		{reflect.Uint8, 8},
		{reflect.Uint16, 16},
		{reflect.Uint32, 32},
		{reflect.Uint64, 64},
		{reflect.Uint, 0},
	}

	for _, tt := range uintTests {
		if got := uintBitSize(tt.kind); got != tt.want {
			t.Errorf("uintBitSize(%s) = %d, want %d", tt.kind, got, tt.want)
		}
	}
}
