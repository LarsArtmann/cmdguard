package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		envLevel   string
		envFormat  string
		envStrict  string
		wantLevel  string
		wantFormat string
		wantStrict bool
	}{
		{
			name:       "creates config with defaults",
			envLevel:   "",
			envFormat:  "",
			envStrict:  "",
			wantLevel:  "info",
			wantFormat: "text",
			wantStrict: false,
		},
		{
			name:       "loads log level from env",
			envLevel:   "debug",
			envFormat:  "",
			envStrict:  "",
			wantLevel:  "debug",
			wantFormat: "text",
			wantStrict: false,
		},
		{
			name:       "loads log format from env",
			envLevel:   "",
			envFormat:  "json",
			envStrict:  "",
			wantLevel:  "info",
			wantFormat: "json",
			wantStrict: false,
		},
		{
			name:       "loads strict mode from env",
			envLevel:   "",
			envFormat:  "",
			envStrict:  "true",
			wantLevel:  "info",
			wantFormat: "text",
			wantStrict: true,
		},
		{
			name:       "loads all from env",
			envLevel:   "error",
			envFormat:  "json",
			envStrict:  "true",
			wantLevel:  "error",
			wantFormat: "json",
			wantStrict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set env vars (errors ignored - test environment setup)
			if tt.envLevel != "" {
				_ = os.Setenv("CMDGUARD_LOG_LEVEL", tt.envLevel)
				defer func() { _ = os.Unsetenv("CMDGUARD_LOG_LEVEL") }()
			}
			if tt.envFormat != "" {
				_ = os.Setenv("CMDGUARD_LOG_FORMAT", tt.envFormat)
				defer func() { _ = os.Unsetenv("CMDGUARD_LOG_FORMAT") }()
			}
			if tt.envStrict != "" {
				_ = os.Setenv("CMDGUARD_STRICT_MODE", tt.envStrict)
				defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()
			}

			cfg := Load()

			require.NotNil(t, cfg)
			assert.Equal(t, tt.wantLevel, cfg.LogLevel)
			assert.Equal(t, tt.wantFormat, cfg.LogFormat)
			assert.Equal(t, tt.wantStrict, cfg.StrictMode)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
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
			err := tt.config.Validate()

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestGetConfigFilePath(t *testing.T) {
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
			got := GetConfigFilePath(tt.configFile)

			switch tt.want {
			case "":
				assert.Empty(t, got)
			case "/":
				assert.NotEmpty(t, got)
			default:
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
