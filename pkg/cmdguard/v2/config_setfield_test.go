package v2

import (
	"strings"
	"testing"
	"time"
)

func TestSetField(t *testing.T) {
	t.Parallel()
	t.Run("set string field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Name string
		}{}

		err := SetField(cfg, "Name", "test")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Name != "test" {
			t.Errorf("expected Name 'test', got %q", cfg.Name)
		}
	})

	t.Run("set int field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Count int
		}{}

		err := SetField(cfg, "Count", 42)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Count != 42 {
			t.Errorf("expected Count 42, got %d", cfg.Count)
		}
	})

	t.Run("set bool field", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Enabled bool
		}{}

		err := SetField(cfg, "Enabled", true)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if !cfg.Enabled {
			t.Error("expected Enabled to be true")
		}
	})

	t.Run("non-pointer config", func(t *testing.T) {
		t.Parallel()

		cfg := struct{ Name string }{}

		err := SetField(cfg, "Name", "test")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "pointer to struct") {
			t.Errorf("expected error to contain 'pointer to struct', got: %v", err)
		}
	})

	t.Run("field not found", func(t *testing.T) {
		t.Parallel()

		cfg := &struct{ Name string }{}

		err := SetField(cfg, "Missing", "test")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("expected error to contain 'not found', got: %v", err)
		}
	})

	t.Run("string to LogLevel", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Level LogLevel
		}{}

		err := SetField(cfg, "Level", "debug")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Level.String() != "debug" {
			t.Errorf("expected Level 'debug', got %q", cfg.Level.String())
		}
	})

	t.Run("string to LogFormat", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Format LogFormat
		}{}

		err := SetField(cfg, "Format", "json")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Format.String() != "json" {
			t.Errorf("expected Format 'json', got %q", cfg.Format.String())
		}
	})

	t.Run("string to Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", "30s")
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Timeout.Duration() != 30*time.Second {
			t.Errorf("expected Timeout 30s, got %v", cfg.Timeout.Duration())
		}
	})

	t.Run("time.Duration to Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", 45*time.Second)
		if err != nil {
			t.Fatalf("expected no error, got: %v", err)
		}

		if cfg.Timeout.Duration() != 45*time.Second {
			t.Errorf("expected Timeout 45s, got %v", cfg.Timeout.Duration())
		}
	})

	t.Run("invalid LogLevel", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Level LogLevel
		}{}

		err := SetField(cfg, "Level", "invalid")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid LogFormat", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Format LogFormat
		}{}

		err := SetField(cfg, "Format", "invalid-format")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid Duration", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Timeout Duration
		}{}

		err := SetField(cfg, "Timeout", "not-a-duration")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("incompatible types", func(t *testing.T) {
		t.Parallel()

		cfg := &struct {
			Name string
		}{}
		// Use a struct type which is truly incompatible with string
		err := SetField(cfg, "Name", Duration{duration: 5 * time.Minute})
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if !strings.Contains(err.Error(), "cannot convert") {
			t.Errorf("expected error to contain 'cannot convert', got: %v", err)
		}
	})
}
