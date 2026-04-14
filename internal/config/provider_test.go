package config

import (
	"strings"
	"testing"
)

// loadTestCase holds test case data for TestLoad.
type loadTestCase struct {
	name       string
	envLevel   string
	envFormat  string
	envStrict  string
	wantLevel  string
	wantFormat string
	wantStrict bool
}

// newLoadTestCase creates a loadTestCase with the given values.
func newLoadTestCase(name, envLevel, envFormat, envStrict, wantLevel, wantFormat string, wantStrict bool) loadTestCase {
	return loadTestCase{
		name:       name,
		envLevel:   envLevel,
		envFormat:  envFormat,
		envStrict:  envStrict,
		wantLevel:  wantLevel,
		wantFormat: wantFormat,
		wantStrict: wantStrict,
	}
}

func TestLoad(t *testing.T) {
	tests := []loadTestCase{
		newLoadTestCase("creates config with defaults", "", "", "", "info", "text", false),
		newLoadTestCase("loads log level from env", "debug", "", "", "debug", "text", false),
		newLoadTestCase("loads log format from env", "", "json", "", "info", "json", false),
		newLoadTestCase("loads strict mode from env", "", "", "true", "info", "text", true),
		newLoadTestCase("loads all from env", "error", "json", "true", "error", "json", true),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envLevel != "" {
				t.Setenv("CMDGUARD_LOG_LEVEL", tt.envLevel)
			}

			if tt.envFormat != "" {
				t.Setenv("CMDGUARD_LOG_FORMAT", tt.envFormat)
			}

			if tt.envStrict != "" {
				t.Setenv("CMDGUARD_STRICT_MODE", tt.envStrict)
			}

			cfg := Load()

			if cfg == nil {
				t.Fatal("Load() returned nil, expected non-nil config")
			}

			if cfg.LogLevel != tt.wantLevel {
				t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, tt.wantLevel)
			}

			if cfg.LogFormat != tt.wantFormat {
				t.Errorf("cfg.LogFormat = %q, want %q", cfg.LogFormat, tt.wantFormat)
			}

			if cfg.StrictMode != tt.wantStrict {
				t.Errorf("cfg.StrictMode = %v, want %v", cfg.StrictMode, tt.wantStrict)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid log level debug",
			config:  Config{LogLevel: "debug"},
			wantErr: false,
		},
		{
			name:    "valid log level info",
			config:  Config{LogLevel: "info"},
			wantErr: false,
		},
		{
			name:    "valid log level warn",
			config:  Config{LogLevel: "warn"},
			wantErr: false,
		},
		{
			name:    "valid log level error",
			config:  Config{LogLevel: "error"},
			wantErr: false,
		},
		{
			name:    "invalid log level",
			config:  Config{LogLevel: "invalid"},
			wantErr: true,
			errMsg:  "invalid log level",
		},
		{
			name:    "empty log level is valid",
			config:  Config{LogLevel: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.errMsg)
				}

				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errMsg)
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		configFile string
		want       string
	}{
		{
			name:       "empty path returns empty",
			configFile: "",
			want:       "",
		},
		{
			name:       "relative path is converted to absolute",
			configFile: "config.yaml",
			want:       "/", // Will start with /
		},
		{
			name:       "absolute path is returned as-is",
			configFile: "/etc/cmdguard/config.yaml",
			want:       "/etc/cmdguard/config.yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GetConfigFilePath(tt.configFile)

			switch tt.want {
			case "":
				if got != "" {
					t.Errorf("GetConfigFilePath(%q) = %q, want empty", tt.configFile, got)
				}
			case "/":
				if got == "" {
					t.Errorf("GetConfigFilePath(%q) = empty, want non-empty", tt.configFile)
				}
			default:
				if got != tt.want {
					t.Errorf("GetConfigFilePath(%q) = %q, want %q", tt.configFile, got, tt.want)
				}
			}
		})
	}
}
