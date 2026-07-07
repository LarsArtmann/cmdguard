package v3

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	t.Parallel()
	t.Run("valid config", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name  string `flag:"name"`
			Count int    `flag:"count"`
		}

		err := ValidateConfig(TestConfig{Name: "test", Count: 10})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("nil config", func(t *testing.T) {
		t.Parallel()

		err := ValidateConfig(nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "must not be nil")
	})

	t.Run("non-struct config", func(t *testing.T) {
		t.Parallel()

		err := ValidateConfig("not a struct")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "expected struct")
	})

	t.Run("valid enum value", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}

		err := ValidateConfig(TestConfig{Level: "info"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("invalid enum value", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Level string `flag:"level" values:"debug,info,warn"`
		}

		err := ValidateConfig(TestConfig{Level: "invalid"})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		assertErrorContains(t, err, "config validation")
	})

	t.Run("pointer to config", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Name string `flag:"name"`
		}

		err := ValidateConfig(&TestConfig{Name: "test"})
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("LogLevel field with values", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Level LogLevel `flag:"level" values:"debug,info,warn,error"`
		}

		cfg := TestConfig{Level: LogLevelInfo}

		err := ValidateConfig(cfg)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})

	t.Run("LogFormat field with values", func(t *testing.T) {
		t.Parallel()

		type TestConfig struct {
			Format LogFormat `flag:"format" values:"text,json"`
		}

		cfg := TestConfig{Format: LogFormatText}

		err := ValidateConfig(cfg)
		if err != nil {
			t.Errorf("expected no error, got: %v", err)
		}
	})
}
