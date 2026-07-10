package v3

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/pflag"

	"github.com/larsartmann/cmdguard/v3/pkg/testutil"
)

// newCountFlagTag builds a FlagTag for the verbose int-count flag used by
// type-handler tests. Pass short="" to omit the short alias.
func newCountFlagTag(short, help string) FlagTag {
	return FlagTag{
		Name: "verbose", Short: short, Help: help,
		Type: reflect.TypeFor[int](), Count: true,
	}
}

// registerCountFlagAndAssert builds a count flag in a fresh flag set, registers
// it via the given register function, and asserts the lookup succeeded.
// Used by Register/Dispatch test pairs that only differ in the registration entry point.
func registerCountFlagAndAssert(t *testing.T, register func(fs *pflag.FlagSet, tag FlagTag) error, short, help string) {
	t.Helper()

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	tag := newCountFlagTag(short, help)

	testutil.AssertNoError(t, register(fs, tag))

	if f := fs.Lookup("verbose"); f == nil {
		t.Error("flag 'verbose' should be registered")
	}
}

// countFlagTestCases runs the with-short and without-short count-flag subtests
// sharing a single help string. It exists so callers don't rewrite the same
// t.Run pair (with/without short) per registration entry point.
func countFlagTestCases(t *testing.T, body func(short, help string), help string) {
	t.Helper()

	cases := []struct {
		name  string
		short string
	}{
		{"with short", "v"},
		{"without short", ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			body(tc.short, help)
		})
	}
}

func TestTypeHandler_RegisterCustomTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     FlagTag
		wantErr bool
	}{
		{
			name: "Duration type registers string flag",
			tag: FlagTag{
				Name: "timeout", Short: "t", Default: "5s",
				Help: "timeout duration", Type: reflect.TypeFor[Duration](),
			},
		},
		{
			name: "Enum type registers string flag with values",
			tag: FlagTag{
				Name: "mode", Default: "fast",
				Help: "run mode", Type: reflect.TypeFor[Enum](),
				Values: []string{"fast", "slow"},
			},
		},
		{
			name: "LogLevel type registers string flag",
			tag: FlagTag{
				Name: "log-level", Default: "info",
				Help: "log level", Type: reflect.TypeFor[LogLevel](),
			},
		},
		{
			name: "LogFormat type registers string flag",
			tag: FlagTag{
				Name: "log-format", Default: "text",
				Help: "log format", Type: reflect.TypeFor[LogFormat](),
			},
		},
		{
			name: "URL type registers string flag",
			tag: FlagTag{
				Name: "endpoint", Default: "http://localhost",
				Help: "API endpoint", Type: reflect.TypeFor[URL](),
			},
		},
		{
			name: "Email type registers string flag",
			tag: FlagTag{
				Name: "email", Default: "test@example.com",
				Help: "contact email", Type: reflect.TypeFor[Email](),
			},
		},
		{
			name: "Port type registers string flag",
			tag: FlagTag{
				Name: "port", Default: "8080",
				Help: "listen port", Type: reflect.TypeFor[Port](),
			},
		},
		{
			name: "FilePath type registers string flag",
			tag: FlagTag{
				Name: "path", Default: "/tmp",
				Help: "file path", Type: reflect.TypeFor[FilePath](),
			},
		},
		{
			name: "HostPort type registers string flag",
			tag: FlagTag{
				Name: "addr", Default: "localhost:8080",
				Help: "host:port", Type: reflect.TypeFor[HostPort](),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
			h, ok := globalTypeRegistry.lookupHandler(tt.tag.Type)
			testutil.AssertBoolTrue(t, ok, "handler should be registered for "+tt.tag.Type.String())

			err := h.Register(fs, tt.tag)
			if tt.wantErr {
				testutil.AssertExpectedError(t, err)
			} else {
				testutil.AssertNoError(t, err)
				testutil.AssertBoolTrue(
					t,
					fs.HasFlags(),
					"flagset should have flags after register",
				)
			}
		})
	}
}

func TestTypeHandler_RegisterCustomTypes_Parse(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tag     FlagTag
		value   string
		wantErr bool
	}{
		{
			name:  "Duration parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[Duration]()},
			value: "5m",
		},
		{
			name:    "Duration parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[Duration]()},
			value:   "not-a-duration",
			wantErr: true,
		},
		{
			name:  "LogLevel parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[LogLevel]()},
			value: "info",
		},
		{
			name:    "LogLevel parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[LogLevel]()},
			value:   "verbose",
			wantErr: true,
		},
		{
			name:  "LogFormat parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[LogFormat]()},
			value: "json",
		},
		{
			name:  "URL parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[URL]()},
			value: "https://example.com",
		},
		{
			name:    "URL parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[URL]()},
			value:   "://bad",
			wantErr: true,
		},
		{
			name:  "Email parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[Email]()},
			value: "user@example.com",
		},
		{
			name:    "Email parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[Email]()},
			value:   "not-an-email",
			wantErr: true,
		},
		{
			name:  "Port parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[Port]()},
			value: "8080",
		},
		{
			name:    "Port parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[Port]()},
			value:   "99999",
			wantErr: true,
		},
		{
			name:  "FilePath parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[FilePath]()},
			value: "/tmp/test",
		},
		{
			name:  "HostPort parse valid",
			tag:   FlagTag{Type: reflect.TypeFor[HostPort]()},
			value: "localhost:8080",
		},
		{
			name:    "HostPort parse invalid",
			tag:     FlagTag{Type: reflect.TypeFor[HostPort]()},
			value:   "localhost:99999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h, ok := globalTypeRegistry.lookupHandler(tt.tag.Type)
			testutil.AssertBoolTrue(t, ok, "handler should exist")

			_, err := h.Parse(tt.value, tt.tag)
			if tt.wantErr {
				testutil.AssertExpectedError(t, err)
			} else {
				testutil.AssertNoError(t, err)
			}
		})
	}
}

func TestTypeHandler_RegisterCustomTypes_Default(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		tag      FlagTag
		expected any
	}{
		{
			name:     "Duration default returns parsed value",
			tag:      FlagTag{Type: reflect.TypeFor[Duration](), Default: "10s"},
			expected: FromDuration(10 * time.Second),
		},
		{
			name:     "Enum default",
			tag:      FlagTag{Type: reflect.TypeFor[Enum](), Default: "fast"},
			expected: "fast",
		},
		{
			name:     "LogLevel default",
			tag:      FlagTag{Type: reflect.TypeFor[LogLevel](), Default: "warn"},
			expected: "warn",
		},
		{
			name:     "URL default",
			tag:      FlagTag{Type: reflect.TypeFor[URL](), Default: "http://example.com"},
			expected: "http://example.com",
		},
		{
			name:     "Port default",
			tag:      FlagTag{Type: reflect.TypeFor[Port](), Default: "443"},
			expected: "443",
		},
		{
			name:     "empty default returns empty string",
			tag:      FlagTag{Type: reflect.TypeFor[URL](), Default: ""},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h, ok := globalTypeRegistry.lookupHandler(tt.tag.Type)
			testutil.AssertBoolTrue(t, ok, "handler should exist")

			result := h.Default(tt.tag)
			if result != tt.expected {
				t.Errorf("Default() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypeHandler_RegisterCountHandler(t *testing.T) {
	t.Parallel()

	t.Run("count handler is registered", func(t *testing.T) {
		t.Parallel()
		if globalTypeRegistry.countHandler == nil {
			t.Error("count handler should be registered")
		}
	})

	countFlagTestCases(t, func(short, help string) {
		t.Helper()

		registerCountFlagAndAssert(
			t,
			func(fs *pflag.FlagSet, tag FlagTag) error {
				return globalTypeRegistry.countHandler.Register(fs, tag)
			},
			short, help,
		)
	}, "verbosity level")

	t.Run("count handler parses integer value", func(t *testing.T) {
		t.Parallel()

		tag := FlagTag{Type: reflect.TypeFor[int](), Count: true}
		result, err := globalTypeRegistry.countHandler.Parse("3", tag)
		testutil.AssertNoError(t, err)
		if result != 3 {
			t.Errorf("Parse(\"3\") = %v, want 3", result)
		}
	})

	t.Run("count handler default is 0", func(t *testing.T) {
		t.Parallel()

		tag := FlagTag{Type: reflect.TypeFor[int](), Count: true}
		result := globalTypeRegistry.countHandler.Default(tag)
		if result != 0 {
			t.Errorf("Default() = %v, want 0", result)
		}
	})
}

func TestTypeHandler_DispatchRegister_WithCount(t *testing.T) {
	t.Parallel()

	countFlagTestCases(t, func(short, help string) {
		t.Helper()

		registerCountFlagAndAssert(
			t,
			func(fs *pflag.FlagSet, tag FlagTag) error {
				return dispatchRegister(globalTypeRegistry, fs, tag)
			},
			short, help,
		)
	}, "verbosity")
}

func TestTypeHandler_LookupHandler_Fallback(t *testing.T) {
	t.Parallel()

	t.Run("exact type match takes priority", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[Duration]())
		testutil.AssertBoolTrue(t, ok, "should find Duration by exact type")
	})

	t.Run("kind fallback for int", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[int]())
		testutil.AssertBoolTrue(t, ok, "should find int by kind")
	})

	t.Run("kind fallback for string", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[string]())
		testutil.AssertBoolTrue(t, ok, "should find string by kind")
	})

	t.Run("kind fallback for bool", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[bool]())
		testutil.AssertBoolTrue(t, ok, "should find bool by kind")
	})

	t.Run("kind fallback for float64", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[float64]())
		testutil.AssertBoolTrue(t, ok, "should find float64 by kind")
	})

	t.Run("kind fallback for uint", func(t *testing.T) {
		t.Parallel()

		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[uint]())
		testutil.AssertBoolTrue(t, ok, "should find uint by kind")
	})

	t.Run("unregistered type returns false", func(t *testing.T) {
		t.Parallel()

		type Custom struct{}
		_, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[Custom]())
		testutil.AssertBoolFalse(t, ok, "unregistered struct type should not have handler")
	})
}

func TestTypeHandler_RegisterTypeHandler_PublicAPI(t *testing.T) {
	t.Run("custom type can be registered and looked up", func(t *testing.T) {
		type Widget struct{ Name string }

		custom := TypeHandlerFunc{
			ParseFunc: func(value string, _ FlagTag) (any, error) {
				return Widget{Name: value}, nil
			},
			DefaultFunc: func(_ FlagTag) any {
				return Widget{Name: "default"}
			},
		}

		testutil.AssertNoError(t, RegisterTypeHandler(reflect.TypeFor[Widget](), custom))

		h, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[Widget]())
		testutil.AssertBoolTrue(t, ok, "custom type should be found after registration")

		result, err := h.Parse("test", FlagTag{})
		testutil.AssertNoError(t, err)

		w, ok := result.(Widget)
		testutil.AssertBoolTrue(t, ok, "result should be Widget")
		testutil.AssertEqual(t, w.Name, "test")
	})
}

func TestTypeHandler_DispatchDefault(t *testing.T) {
	t.Parallel()

	t.Run("dispatchDefault table", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name     string
			tag      FlagTag
			expected any
		}{
			{
				"empty default returns zero value",
				FlagTag{Type: reflect.TypeFor[int](), Default: ""},
				0,
			},
			{
				"int default converts to exact type",
				FlagTag{Type: reflect.TypeFor[int](), Default: "42"},
				42,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := dispatchDefault(globalTypeRegistry, tt.tag)
				testutil.AssertEqual(t, result, tt.expected)
			})
		}
	})

	t.Run("uint32 default converts to exact type", func(t *testing.T) {
		t.Parallel()

		tag := FlagTag{Type: reflect.TypeFor[uint32](), Default: "100"}
		result := dispatchDefault(globalTypeRegistry, tag)
		if result != uint32(100) {
			t.Errorf("dispatchDefault(100) = %v (%T), want uint32(100)", result, result)
		}
	})

	t.Run("unknown type returns raw default string", func(t *testing.T) {
		t.Parallel()

		type Foo struct{}
		tag := FlagTag{Type: reflect.TypeFor[Foo](), Default: "bar"}
		result := dispatchDefault(globalTypeRegistry, tag)
		testutil.AssertEqual(t, result, "bar")
	})
}

func TestTypeHandler_DispatchParse(t *testing.T) {
	t.Parallel()

	t.Run("int parse returns int64", func(t *testing.T) {
		t.Parallel()

		tag := FlagTag{Type: reflect.TypeFor[int]()}
		result, err := dispatchParse(globalTypeRegistry, "42", tag)
		testutil.AssertNoError(t, err)
		if result != int64(42) {
			t.Errorf("dispatchParse(42) = %v, want int64(42)", result)
		}
	})

	t.Run("unknown type returns raw string", func(t *testing.T) {
		t.Parallel()

		type Custom struct{}
		tag := FlagTag{Type: reflect.TypeFor[Custom]()}
		result, err := dispatchParse(globalTypeRegistry, "hello", tag)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, result, "hello")
	})
}

func TestHandledByTypeRegistry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		typ     reflect.Type
		handled bool
	}{
		{"string is handled", reflect.TypeFor[string](), true},
		{"int is handled", reflect.TypeFor[int](), true},
		{"bool is handled", reflect.TypeFor[bool](), true},
		{"Duration is handled", reflect.TypeFor[Duration](), true},
		{"URL is handled", reflect.TypeFor[URL](), true},
		{"unregistered struct is not handled", reflect.TypeFor[struct{ X int }](), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := handledByTypeRegistry(globalTypeRegistry, tt.typ)
			if result != tt.handled {
				t.Errorf("handledByTypeRegistry(%v) = %v, want %v", tt.typ, result, tt.handled)
			}
		})
	}
}

func TestTypeHandler_KindHandlers(t *testing.T) {
	t.Parallel()

	t.Run("string kind register and parse", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "name",
			Short:   "n",
			Default: "world",
			Help:    "name",
			Type:    reflect.TypeFor[string](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertNoError(t, err)

		result, err := dispatchParse(globalTypeRegistry, "hello", tag)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, result, "hello")
	})

	t.Run("bool kind register and parse", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "debug",
			Short:   "d",
			Default: "false",
			Help:    "debug mode",
			Type:    reflect.TypeFor[bool](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertNoError(t, err)

		result, err := dispatchParse(globalTypeRegistry, "true", tag)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, result, true)
	})

	t.Run("float64 kind register and parse", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "ratio",
			Default: "1.5",
			Help:    "ratio",
			Type:    reflect.TypeFor[float64](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertNoError(t, err)

		result, err := dispatchParse(globalTypeRegistry, "2.7", tag)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, result, 2.7)
	})

	t.Run("slice kind register and parse", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "tags",
			Default: "a,b",
			Help:    "tags",
			Type:    reflect.TypeFor[[]string](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertNoError(t, err)

		result, err := dispatchParse(globalTypeRegistry, "x,y,z", tag)
		testutil.AssertNoError(t, err)
		s, ok := result.([]string)
		testutil.AssertBoolTrue(t, ok, "result should be []string")
		if len(s) != 3 || s[0] != "x" || s[2] != "z" {
			t.Errorf("Parse(x,y,z) = %v, want [x y z]", s)
		}
	})

	t.Run("invalid bool default returns error", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "bad",
			Default: "notbool",
			Help:    "bad bool",
			Type:    reflect.TypeFor[bool](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertExpectedError(t, err)
	})

	t.Run("invalid int default returns error", func(t *testing.T) {
		t.Parallel()

		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		tag := FlagTag{
			Name:    "bad",
			Default: "notint",
			Help:    "bad int",
			Type:    reflect.TypeFor[int](),
		}
		err := dispatchRegister(globalTypeRegistry, fs, tag)
		testutil.AssertExpectedError(t, err)
	})
}

func TestTypeHandlerFunc_NilRegisterFunc(t *testing.T) {
	t.Parallel()

	t.Run("nil RegisterFunc returns nil", func(t *testing.T) {
		t.Parallel()

		h := TypeHandlerFunc{
			ParseFunc:   func(_ string, _ FlagTag) (any, error) { return "x", nil },
			DefaultFunc: func(_ FlagTag) any { return "y" },
		}

		err := h.Register(nil, FlagTag{})
		testutil.AssertNoError(t, err)
	})
}

func TestRegisterGoDurationHandler(t *testing.T) {
	t.Parallel()

	RegisterGoDurationHandler()

	t.Run("registers time.Duration handler", func(t *testing.T) {
		h, ok := globalTypeRegistry.lookupHandler(reflect.TypeFor[time.Duration]())
		testutil.AssertBoolTrue(t, ok, "time.Duration should be registered")

		result, err := h.Parse("5s", FlagTag{})
		testutil.AssertNoError(t, err)
		if result != 5*time.Second {
			t.Errorf("Parse(5s) = %v, want 5s", result)
		}
	})

	t.Run("time.Duration default", func(t *testing.T) {
		h, _ := globalTypeRegistry.lookupHandler(reflect.TypeFor[time.Duration]())
		result := h.Default(FlagTag{Default: "10m"})
		if result != 10*time.Minute {
			t.Errorf("Default(10m) = %v, want 10m", result)
		}
	})

	t.Run("time.Duration default table", func(t *testing.T) {
		t.Parallel()

		h, _ := globalTypeRegistry.lookupHandler(reflect.TypeFor[time.Duration]())

		tests := []struct {
			name     string
			defVal   string
			expected time.Duration
		}{
			{"empty default returns zero", "", time.Duration(0)},
			{"invalid default returns zero", "not-a-duration", time.Duration(0)},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				result := h.Default(FlagTag{Default: tt.defVal})
				if result != tt.expected {
					t.Errorf("Default(%s) = %v, want %v", tt.defVal, result, tt.expected)
				}
			})
		}
	})

	t.Run("time.Duration register with short", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		h, _ := globalTypeRegistry.lookupHandler(reflect.TypeFor[time.Duration]())
		err := h.Register(fs, FlagTag{Name: "timeout", Short: "t", Default: "30s", Help: "timeout"})
		testutil.AssertNoError(t, err)

		f := fs.Lookup("timeout")
		if f == nil {
			t.Error("flag 'timeout' should be registered")
		}
	})

	t.Run("time.Duration invalid default returns error", func(t *testing.T) {
		fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
		h, _ := globalTypeRegistry.lookupHandler(reflect.TypeFor[time.Duration]())
		err := h.Register(fs, FlagTag{Name: "bad", Default: "not-a-duration", Help: "bad duration"})
		testutil.AssertExpectedError(t, err)
	})
}
