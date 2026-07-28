package v4

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateFlagPrototype(t *testing.T) {
	t.Parallel()

	type flags struct {
		Name string
	}

	t.Run("returns non-nil pointer as-is", func(t *testing.T) {
		t.Parallel()

		original := &flags{Name: "test"}

		proto := createFlagPrototype(original)
		if proto != original {
			t.Error("expected same pointer for non-nil input")
		}
	})

	t.Run("creates new instance for nil pointer", func(t *testing.T) {
		t.Parallel()

		var f *flags

		proto := createFlagPrototype(f)
		if proto == nil {
			t.Fatal("expected non-nil prototype for nil pointer input")
		}
	})

	t.Run("returns struct as-is", func(t *testing.T) {
		t.Parallel()

		original := flags{Name: "test"}

		proto := createFlagPrototype(original)
		if proto.Name != original.Name {
			t.Errorf("Name = %q, want %q", proto.Name, original.Name)
		}
	})

	t.Run("returns zero value for NoFlags", func(t *testing.T) {
		t.Parallel()

		proto := createFlagPrototype(NoFlags{})

		var want NoFlags
		if proto != want {
			t.Error("expected zero NoFlags")
		}
	})
}

func TestIsNilPointer(t *testing.T) {
	t.Parallel()
	t.Run("nil interface", func(t *testing.T) {
		t.Parallel()

		if !isNilPointer(nil) {
			t.Error("expected true for nil interface")
		}
	})

	t.Run("nil pointer", func(t *testing.T) {
		t.Parallel()

		var p *int
		if !isNilPointer(p) {
			t.Error("expected true for nil *int")
		}
	})

	t.Run("non-nil pointer", func(t *testing.T) {
		t.Parallel()

		x := 42
		if isNilPointer(&x) {
			t.Error("expected false for non-nil *int")
		}
	})

	t.Run("struct", func(t *testing.T) {
		t.Parallel()

		if isNilPointer(struct{}{}) {
			t.Error("expected false for struct{}")
		}
	})

	t.Run("nil slice", func(t *testing.T) {
		t.Parallel()

		var s []string
		if !isNilPointer(s) {
			t.Error("expected true for nil slice")
		}
	})

	t.Run("nil map", func(t *testing.T) {
		t.Parallel()

		var m map[string]string
		if !isNilPointer(m) {
			t.Error("expected true for nil map")
		}
	})

	t.Run("non-empty string", func(t *testing.T) {
		t.Parallel()

		if isNilPointer("hello") {
			t.Error("expected false for string")
		}
	})
}

func TestCreateNilFlags(t *testing.T) {
	t.Parallel()

	type flags struct {
		Name string `default:"default" flag:"name"`
	}

	t.Run("creates instance for pointer type", func(t *testing.T) {
		t.Parallel()

		fc, ptr, err := createNilFlags[*flags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if fc == nil {
			t.Error("expected non-nil flags copy")
		}

		if ptr == nil {
			t.Error("expected non-nil flags pointer")
		}
	})

	t.Run("creates instance for struct type", func(t *testing.T) {
		t.Parallel()

		fc, ptr, err := createNilFlags[flags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if fc.Name != "" {
			t.Errorf("expected zero value, got Name=%q", fc.Name)
		}

		if ptr == nil {
			t.Error("expected non-nil flags pointer")
		}
	})

	t.Run("returns struct pointer for NoFlags", func(t *testing.T) {
		t.Parallel()

		fc, ptr, err := createNilFlags[NoFlags]()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var wantNF NoFlags
		if fc != wantNF {
			t.Error("expected zero NoFlags")
		}

		if ptr == nil {
			t.Error("expected non-nil pointer for struct{} type")
		}
	})
}

func TestFlagsToPtr(t *testing.T) {
	t.Parallel()

	type flags struct {
		Name string
	}

	t.Run("pointer returned as-is", func(t *testing.T) {
		t.Parallel()

		original := &flags{Name: "test"}

		fc, ptr := flagsToPtr(original)
		if fc != original {
			t.Error("expected same pointer")
		}

		if ptr != original {
			t.Error("expected same pointer")
		}
	})

	t.Run("struct creates new pointer", func(t *testing.T) {
		t.Parallel()

		original := flags{Name: "test"}

		fc, ptr := flagsToPtr(original)
		if fc.Name != original.Name {
			t.Errorf("Name = %q, want %q", fc.Name, original.Name)
		}

		if ptr == nil {
			t.Error("expected non-nil pointer")
		}
	})
}

func TestParseAndSyncFlags(t *testing.T) {
	t.Parallel()
	t.Run("returns flags unchanged with nil registry", func(t *testing.T) {
		t.Parallel()

		type flags struct {
			Name string
		}

		f := flags{Name: "test"}
		cmd := &cobra.Command{Use: "test"}

		result, err := parseAndSyncFlags(cmd, f, nil, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Name != "test" {
			t.Errorf("Name = %q, want %q", result.Name, "test")
		}
	})
}

func TestCloneAndParseFlags(t *testing.T) {
	t.Parallel()

	type flags struct {
		Name string
	}

	t.Run("clones and parses pointer flags", func(t *testing.T) {
		t.Parallel()

		original := &flags{Name: "original"}
		cmd := &cobra.Command{Use: "test"}

		result, err := cloneAndParseFlags(cmd, original, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Name != "original" {
			t.Errorf("Name = %q, want %q", result.Name, "original")
		}
	})

	t.Run("clones and parses struct flags", func(t *testing.T) {
		t.Parallel()

		original := flags{Name: "original"}
		cmd := &cobra.Command{Use: "test"}

		result, err := cloneAndParseFlags(cmd, original, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Name != "original" {
			t.Errorf("Name = %q, want %q", result.Name, "original")
		}
	})

	t.Run("creates instance for nil pointer flags", func(t *testing.T) {
		t.Parallel()

		var f *flags

		cmd := &cobra.Command{Use: "test"}

		result, err := cloneAndParseFlags(cmd, f, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result == nil {
			t.Error("expected non-nil result for nil pointer input")
		}
	})

	t.Run("handles NoFlags", func(t *testing.T) {
		t.Parallel()

		cmd := &cobra.Command{Use: "test"}

		result, err := cloneAndParseFlags(cmd, NoFlags{}, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result != (NoFlags{}) {
			t.Error("expected zero NoFlags")
		}
	})
}

func TestFormatFieldValue(t *testing.T) {
	t.Parallel()

	t.Run("invalid value returns empty", func(t *testing.T) {
		t.Parallel()

		result := formatFieldValue(reflect.Value{})
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("string kind", func(t *testing.T) {
		t.Parallel()

		v := reflect.ValueOf("hello")
		if got := formatFieldValue(v); got != "hello" {
			t.Errorf("got %q, want %q", got, "hello")
		}
	})

	t.Run("integer kinds", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name  string
			value any
			want  string
		}{
			{"int", int(42), "42"},
			{"int8", int8(8), "8"},
			{"int16", int16(16), "16"},
			{"int32", int32(32), "32"},
			{"int64", int64(64), "64"},
			{"uint", uint(42), "42"},
			{"uint8", uint8(8), "8"},
			{"uint16", uint16(16), "16"},
			{"uint32", uint32(32), "32"},
			{"uint64", uint64(64), "64"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				got := formatFieldValue(reflect.ValueOf(tt.value))
				if got != tt.want {
					t.Errorf("got %q, want %q", got, tt.want)
				}
			})
		}
	})

	t.Run("float kind", func(t *testing.T) {
		t.Parallel()

		got := formatFieldValue(reflect.ValueOf(float64(3.14)))
		if got != "3.14" {
			t.Errorf("got %q, want %q", got, "3.14")
		}
	})

	t.Run("bool kind", func(t *testing.T) {
		t.Parallel()

		got := formatFieldValue(reflect.ValueOf(true))
		if got != "true" {
			t.Errorf("got %q, want %q", got, "true")
		}
	})

	t.Run("pointer dereferences", func(t *testing.T) {
		t.Parallel()

		s := "pointed"
		got := formatFieldValue(reflect.ValueOf(&s))
		if got != "pointed" {
			t.Errorf("got %q, want %q", got, "pointed")
		}
	})

	t.Run("nil pointer returns empty", func(t *testing.T) {
		t.Parallel()

		var s *string
		got := formatFieldValue(reflect.ValueOf(s))
		if got != "" {
			t.Errorf("expected empty for nil pointer, got %q", got)
		}
	})

	t.Run("fmt.Stringer", func(t *testing.T) {
		t.Parallel()

		s := Duration{duration: 5e9}
		got := formatFieldValue(reflect.ValueOf(s))
		if got != "5s" {
			t.Errorf("got %q, want %q", got, "5s")
		}
	})

	t.Run("struct falls through to Sprintf", func(t *testing.T) {
		t.Parallel()

		type plain struct{ X int }

		got := formatFieldValue(reflect.ValueOf(plain{X: 7}))
		if got != "{7}" {
			t.Errorf("got %q, want %q", got, "{7}")
		}
	})
}
