package v4

import (
	"reflect"
	"testing"

	"github.com/larsartmann/cmdguard/v4/pkg/testutil"
)

type cowTestType struct{ Value string }

const cowTestValidator = "cow-test-validator"

func TestCOWTypeRegistry_InstanceWriteDoesNotLeakToGlobal(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	registry.RegisterTypeHandler(reflect.TypeFor[cowTestType](), TypeHandlerFunc{
		ParseFunc:   func(value string, _ FlagTag) (any, error) { return cowTestType{Value: value}, nil },
		DefaultFunc: func(_ FlagTag) any { return cowTestType{} },
	})

	_, found := registry.types.lookupHandler(reflect.TypeFor[cowTestType]())
	testutil.AssertBoolTrue(t, found, "instance should have custom handler")

	_, globalHas := globalTypeRegistry.lookupHandler(reflect.TypeFor[cowTestType]())
	testutil.AssertBoolFalse(t, globalHas, "global should NOT have instance-registered handler")
}

func TestCOWTypeRegistry_TwoInstancesAreIsolated(t *testing.T) {
	t.Parallel()

	reg1, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	reg2, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	reg1.RegisterTypeHandler(reflect.TypeFor[cowTestType](), TypeHandlerFunc{
		ParseFunc:   func(value string, _ FlagTag) (any, error) { return cowTestType{Value: "reg1"}, nil },
		DefaultFunc: func(_ FlagTag) any { return cowTestType{} },
	})

	_, reg1Has := reg1.types.lookupHandler(reflect.TypeFor[cowTestType]())
	testutil.AssertBoolTrue(t, reg1Has, "reg1 should have custom handler")

	_, reg2Has := reg2.types.lookupHandler(reflect.TypeFor[cowTestType]())
	testutil.AssertBoolFalse(t, reg2Has, "reg2 should NOT have reg1's handler")
}

func TestCOWValidatorRegistry_InstanceWriteDoesNotLeakToGlobal(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(&struct {
		Name string `flag:"name" validate:"cow-test-validator"`
	}{})
	testutil.AssertNoError(t, err)

	registry.RegisterFlagValidator(cowTestValidator, func(_ string) error { return nil })

	_, found := registry.validators.lookup(cowTestValidator)
	testutil.AssertBoolTrue(t, found, "instance should have custom validator")

	_, globalHas := globalValidators.lookup(cowTestValidator)
	testutil.AssertBoolFalse(t, globalHas, "global should NOT have instance-registered validator")
}

func TestCOWValidatorRegistry_TwoInstancesAreIsolated(t *testing.T) {
	t.Parallel()

	reg1, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	reg2, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	reg1.RegisterFlagValidator(cowTestValidator, func(_ string) error { return nil })

	_, reg1Has := reg1.validators.lookup(cowTestValidator)
	testutil.AssertBoolTrue(t, reg1Has, "reg1 should have custom validator")

	_, reg2Has := reg2.validators.lookup(cowTestValidator)
	testutil.AssertBoolFalse(t, reg2Has, "reg2 should NOT have reg1's validator")
}

func TestCOWTypeRegistry_LookupWorksBeforeWrite(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	handler, found := registry.types.lookupHandler(reflect.TypeFor[string]())
	testutil.AssertBoolTrue(t, found, "should find built-in string handler")
	if handler == nil {
		t.Error("handler should not be nil")
	}
}

func TestCOWValidatorRegistry_LookupWorksBeforeWrite(t *testing.T) {
	t.Parallel()

	registry, err := NewFlagRegistry(&struct {
		Name string `flag:"name"`
	}{})
	testutil.AssertNoError(t, err)

	_, found := registry.validators.lookup("email")
	testutil.AssertBoolTrue(t, found, "should find built-in email validator")
}
