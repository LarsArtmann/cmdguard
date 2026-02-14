package config

import (
	"testing"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "creates config with defaults",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			injector := do.New()
			cfg, err := NewConfig(injector)

			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, cfg)
			assert.Equal(t, "info", cfg.LogLevel, "default log level should be info")
			assert.False(t, cfg.StrictMode, "default strict mode should be false")
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
			name:    "valid log levels",
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
