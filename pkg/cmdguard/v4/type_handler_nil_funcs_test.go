package v4

import (
	"reflect"
	"testing"

	"github.com/spf13/pflag"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

// nilHandlerProbeType is a field type no built-in handler claims, so a
// test-registered TypeHandlerFunc is the only handler dispatch can find.
type nilHandlerProbeType string

// nilFuncsHandler returns a TypeHandlerFunc with every optional func nil.
func nilFuncsHandler() TypeHandlerFunc {
	return TypeHandlerFunc{}
}

// TestDispatchRegister_NilRegisterFuncRejected pins the footgun fix: a
// TypeHandlerFunc without RegisterFunc must fail at dispatch (CLI build time)
// instead of silently registering nothing, which used to surface far from the
// cause as "unknown flag" at invocation time.
func TestDispatchRegister_NilRegisterFuncRejected(t *testing.T) {
	t.Parallel()

	registry := newTypeRegistry()
	registry.register(reflect.TypeFor[nilHandlerProbeType](), nilFuncsHandler())

	fs := pflag.NewFlagSet("probe", pflag.ContinueOnError)
	tag := FlagTag{Name: "probe-flag", Type: reflect.TypeFor[nilHandlerProbeType]()}

	err := dispatchRegister(registry, fs, tag)
	testutil.AssertErrorContains(t, err, "RegisterFunc", "nil")

	if f := fs.Lookup("probe-flag"); f != nil {
		t.Errorf("nil RegisterFunc registered the flag anyway: %v", f)
	}
}

// TestTypeHandlerFunc_NilParseFuncErrors pins that Parse reports the missing
// ParseFunc as an error instead of panicking on the nil call.
func TestTypeHandlerFunc_NilParseFuncErrors(t *testing.T) {
	t.Parallel()

	tag := FlagTag{Name: "probe-flag", Type: reflect.TypeFor[nilHandlerProbeType]()}

	var (
		got any
		err error
	)

	testutil.AssertDoesNotPanic(t, func() {
		got, err = nilFuncsHandler().Parse("x", tag)
	})

	if got != nil {
		t.Errorf("Parse with nil ParseFunc returned %v, want nil", got)
	}

	testutil.AssertErrorContains(t, err, "ParseFunc", "nil")
}

// TestTypeHandlerFunc_NilDefaultFuncErrors pins that Default tolerates a nil
// DefaultFunc by yielding a nil default instead of panicking.
func TestTypeHandlerFunc_NilDefaultFuncErrors(t *testing.T) {
	t.Parallel()

	tag := FlagTag{Name: "probe-flag", Type: reflect.TypeFor[nilHandlerProbeType]()}

	var got any

	testutil.AssertDoesNotPanic(t, func() {
		got = nilFuncsHandler().Default(tag)
	})

	if got != nil {
		t.Errorf("Default with nil DefaultFunc returned %v, want nil", got)
	}
}
