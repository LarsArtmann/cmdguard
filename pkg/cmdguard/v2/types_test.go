package v2

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseEnum(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		allowed   []string
		wantErr   bool
		wantValue string
	}{
		{
			name:      "valid value",
			value:     "debug",
			allowed:   []string{"debug", "info", "warn", "error"},
			wantErr:   false,
			wantValue: "debug",
		},
		{
			name:    "invalid value",
			value:   "invalid",
			allowed: []string{"debug", "info", "warn", "error"},
			wantErr: true,
		},
		{
			name:      "empty allowed list",
			value:     "anything",
			allowed:   []string{},
			wantErr:   true,
			wantValue: "",
		},
		{
			name:      "single allowed value",
			value:     "only",
			allowed:   []string{"only"},
			wantErr:   false,
			wantValue: "only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e, err := ParseEnum(tt.value, tt.allowed)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, IsEnumError(err))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantValue, e.String())
				assert.Equal(t, tt.allowed, e.Allowed())
			}
		})
	}
}

func TestMustEnum(t *testing.T) {
	t.Run("valid value", func(t *testing.T) {
		e := MustEnum("info", []string{"debug", "info", "warn"})
		assert.Equal(t, "info", e.String())
	})

	t.Run("panics on invalid", func(t *testing.T) {
		defer func() {
			r := recover()
			require.NotNil(t, r)
			assert.Contains(t, fmt.Sprintf("%v", r), "invalid value")
		}()
		MustEnum("invalid", []string{"valid"})
	})
}

func TestEnum_Methods(t *testing.T) {
	e := MustEnum("test", []string{"a", "test", "b"})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "test", e.String())
	})

	t.Run("Value", func(t *testing.T) {
		assert.Equal(t, "test", e.Value())
	})

	t.Run("Allowed", func(t *testing.T) {
		assert.Equal(t, []string{"a", "test", "b"}, e.Allowed())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		assert.False(t, e.IsEmpty())
		var empty Enum
		assert.True(t, empty.IsEmpty())
	})
}

func TestEnum_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Level Enum `json:"level"`
	}

	t.Run("marshal", func(t *testing.T) {
		c := config{Level: MustEnum("info", []string{"debug", "info", "warn"})}
		data, err := json.Marshal(c)
		require.NoError(t, err)
		assert.JSONEq(t, `{"level":"info"}`, string(data))
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config
		c.Level = Enum{allowed: []string{"debug", "info", "warn"}}
		err := json.Unmarshal([]byte(`{"level":"info"}`), &c)
		require.NoError(t, err)
		assert.Equal(t, "info", c.Level.String())
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config
		c.Level = Enum{allowed: []string{"debug", "info"}}
		err := json.Unmarshal([]byte(`{"level":"invalid"}`), &c)
		require.Error(t, err)
	})

	t.Run("unmarshal with no allowed", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"level":"any"}`), &c)
		require.NoError(t, err)
		assert.Equal(t, "any", c.Level.String())
	})
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantMs  int64
	}{
		{"seconds", "30s", false, 30000},
		{"minutes", "5m", false, 300000},
		{"hours", "1h", false, 3600000},
		{"complex", "1h30m", false, 5400000},
		{"milliseconds", "500ms", false, 500},
		{"invalid", "not-a-duration", true, 0},
		{"empty", "", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := ParseDuration(tt.input)
			if tt.wantErr {
				require.Error(t, err)
				assert.True(t, IsDurationError(err))
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantMs, d.Milliseconds())
			}
		})
	}
}

func TestMustDuration(t *testing.T) {
	t.Run("valid value", func(t *testing.T) {
		d := MustDuration("30s")
		assert.Equal(t, 30*time.Second, d.Duration())
	})

	t.Run("panics on invalid", func(t *testing.T) {
		defer func() {
			r := recover()
			require.NotNil(t, r)
		}()
		MustDuration("invalid")
	})
}

func TestFromDuration(t *testing.T) {
	td := 5 * time.Minute
	d := FromDuration(td)
	assert.Equal(t, td, d.Duration())
	assert.Equal(t, int64(300000), d.Milliseconds())
	assert.Equal(t, float64(300), d.Seconds())
}

func TestDuration_Methods(t *testing.T) {
	d := MustDuration("2h30m")

	t.Run("Duration", func(t *testing.T) {
		assert.Equal(t, 2*time.Hour+30*time.Minute, d.Duration())
	})

	t.Run("String", func(t *testing.T) {
		assert.Equal(t, "2h30m0s", d.String())
	})

	t.Run("IsZero", func(t *testing.T) {
		assert.False(t, d.IsZero())
		var zero Duration
		assert.True(t, zero.IsZero())
	})

	t.Run("Milliseconds", func(t *testing.T) {
		assert.Equal(t, int64(9000000), d.Milliseconds())
	})

	t.Run("Seconds", func(t *testing.T) {
		assert.Equal(t, float64(9000), d.Seconds())
	})
}

func TestDuration_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Timeout Duration `json:"timeout"`
	}

	t.Run("marshal", func(t *testing.T) {
		c := config{Timeout: MustDuration("30s")}
		data, err := json.Marshal(c)
		require.NoError(t, err)
		assert.JSONEq(t, `{"timeout":"30s"}`, string(data))
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"timeout":"1h"}`), &c)
		require.NoError(t, err)
		assert.Equal(t, time.Hour, c.Timeout.Duration())
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"timeout":"invalid"}`), &c)
		require.Error(t, err)
	})
}

func TestLogLevel(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		assert.Equal(t, "debug", LogLevelDebug.String())
		assert.Equal(t, "info", LogLevelInfo.String())
		assert.Equal(t, "warn", LogLevelWarn.String())
		assert.Equal(t, "error", LogLevelError.String())
	})

	t.Run("ParseLogLevel valid", func(t *testing.T) {
		tests := []string{"debug", "info", "warn", "error"}
		for _, v := range tests {
			l, err := ParseLogLevel(v)
			require.NoError(t, err)
			assert.Equal(t, v, l.String())
		}
	})

	t.Run("ParseLogLevel invalid", func(t *testing.T) {
		_, err := ParseLogLevel("invalid")
		require.Error(t, err)
	})
}

func TestLogLevel_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Level LogLevel `json:"level"`
	}

	t.Run("marshal", func(t *testing.T) {
		c := config{Level: LogLevelInfo}
		data, err := json.Marshal(c)
		require.NoError(t, err)
		assert.JSONEq(t, `{"level":"info"}`, string(data))
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"level":"debug"}`), &c)
		require.NoError(t, err)
		assert.Equal(t, "debug", c.Level.String())
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"level":"trace"}`), &c)
		require.Error(t, err)
	})
}

func TestLogFormat(t *testing.T) {
	t.Run("constants", func(t *testing.T) {
		assert.Equal(t, "text", LogFormatText.String())
		assert.Equal(t, "json", LogFormatJSON.String())
	})

	t.Run("ParseLogFormat valid", func(t *testing.T) {
		tests := []string{"text", "json"}
		for _, v := range tests {
			f, err := ParseLogFormat(v)
			require.NoError(t, err)
			assert.Equal(t, v, f.String())
		}
	})

	t.Run("ParseLogFormat invalid", func(t *testing.T) {
		_, err := ParseLogFormat("xml")
		require.Error(t, err)
	})
}

func TestLogFormat_MarshalUnmarshal(t *testing.T) {
	type config struct {
		Format LogFormat `json:"format"`
	}

	t.Run("marshal", func(t *testing.T) {
		c := config{Format: LogFormatJSON}
		data, err := json.Marshal(c)
		require.NoError(t, err)
		assert.JSONEq(t, `{"format":"json"}`, string(data))
	})

	t.Run("unmarshal valid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"format":"text"}`), &c)
		require.NoError(t, err)
		assert.Equal(t, "text", c.Format.String())
	})

	t.Run("unmarshal invalid", func(t *testing.T) {
		var c config
		err := json.Unmarshal([]byte(`{"format":"xml"}`), &c)
		require.Error(t, err)
	})
}

func TestPtr(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		v := 42
		p := Ptr(v)
		require.NotNil(t, p)
		assert.Equal(t, 42, *p)
	})

	t.Run("string", func(t *testing.T) {
		v := "hello"
		p := Ptr(v)
		require.NotNil(t, p)
		assert.Equal(t, "hello", *p)
	})

	t.Run("struct", func(t *testing.T) {
		type s struct{ Name string }
		v := s{Name: "test"}
		p := Ptr(v)
		require.NotNil(t, p)
		assert.Equal(t, "test", p.Name)
	})
}

func TestValueOrDefault(t *testing.T) {
	t.Run("nil pointer returns default", func(t *testing.T) {
		var p *int
		result := ValueOrDefault(p, 10)
		assert.Equal(t, 10, result)
	})

	t.Run("non-nil pointer returns value", func(t *testing.T) {
		p := Ptr(42)
		result := ValueOrDefault(p, 10)
		assert.Equal(t, 42, result)
	})

	t.Run("empty string default", func(t *testing.T) {
		var p *string
		result := ValueOrDefault(p, "default")
		assert.Equal(t, "default", result)
	})
}

func TestEnsureValid(t *testing.T) {
	t.Run("nil returns error", func(t *testing.T) {
		var p *int
		err := EnsureValid(p, "myField")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "myField")
		assert.Contains(t, err.Error(), "must not be nil")
	})

	t.Run("non-nil returns nil", func(t *testing.T) {
		p := Ptr(42)
		err := EnsureValid(p, "myField")
		require.NoError(t, err)
	})
}

// Helper functions for type checking
func IsEnumError(err error) bool {
	var e *EnumError
	return AsEnumError(err, &e)
}

func AsEnumError(err error, target **EnumError) bool {
	switch e := err.(type) {
	case *EnumError:
		*target = e
		return true
	default:
		return false
	}
}

func IsDurationError(err error) bool {
	var e *DurationError
	return AsDurationError(err, &e)
}

func AsDurationError(err error, target **DurationError) bool {
	switch e := err.(type) {
	case *DurationError:
		*target = e
		return true
	default:
		return false
	}
}
