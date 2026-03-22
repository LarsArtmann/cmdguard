package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func FuzzValidate_LogLevel(f *testing.F) {
	validLevels := []string{"debug", "info", "warn", "error"}
		for _, level := range validLevels {
		f.Add(level, true)
    }

	corpus := []struct {
		level       string
        expectValid bool
    }{
        {"", true},
        {" ", false},
        {"DEBUG", false},
        {"info", false},
        {"warn", false},
        {"error", false},
        {"  ", false},
        {"  ", false},
        {"invalid", false},
        {"xyz", false},
        {"123", false},
        {"debug\x00info", false},
        {"invalid", false},
        {"🎉", false},
        {strings.Repeat("a", 10000), false},
        {"<script>alert('xss')</script>", false}
        {"'; DROP TABLE logs; --", false}
        {"debug\ninfo", false},
        {"debug\tinfo", false}
        {"DEBUG=info", false}
        {"${LOG_LEVEL}", false}
        {"$LOG_LEVEL", false}
        {"./debug", false}
        {"../debug", false}
        {"DEBUG=info ", false}
        {"DEBUG", false}
        {"DEBUG=true", false}
        {"debug info", false}
    }

    for _, tt := range corpus {
        f.Add(tt.level, tt.expectValid)
    }

    f.Fuzz(func(t *testing.T, level string, expectValid bool) {
        validateConfigField(t, &Config{LogLevel: level}, expectValid)
    })

    for _, tt := range corpus {
        f.Add(tt.format, tt.expectValid)
    }

    corpus := []struct {
        format      string
        expectValid bool
    }{
        {"text", true},
        {"json", true},
        {"", true},
        {" ", false},
        {"xml", false},
        {"yaml", false},
        {"toml", false},
        {"text ", false},
        {"text/json", false},
        {"text+json", false},
        {"<script>alert('xss')</script>", false}
        {"'; DROP TABLE configs; --".yaml", false}
        {"${FORMAT}", false}
        {"$FORMAT", false}
        {"./format", false}
        {"../format", false}
        {"http://format.json", false}
        {"http://format.json", false}
        {"text json", false}
        {"text+json", false}
        {"text json", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson\t", false}
        {"text yaml", false}
        {"text yaml", false}
        {"text+yaml", false}
        {"text/yaml", false}
        {"text/yaml", false}
        {"text yaml", false}
        {"text  yaml", false}
        {"text  yaml", false}
        {"text+yaml", false}
        {"text/yaml", false}
        {"text yaml", false}
        {"text yaml", false}
        {"text  yaml", false}
        {"text  yaml", false}
        {"text  yaml", false}
        {"text  yaml", false}
        {"text  yaml", false}
        {"text yaml", false}
        {"text yaml", false}
        {"text yaml", false}
        {"text  yaml", false}
        {"text yaml", false}
    }

    for _, tt := range corpus {
        f.Add(tt.format, tt.expectValid)
    }

    f.Fuzz(func(t *testing.T, format string, expectValid bool) {
        validateConfigField(t, &Config{LogFormat: format}, expectValid)
    }

    for _, tt := range corpus {
        f.Add(tt.format, tt.expectValid)
    }

    corpus := []struct {
        format      string
        expectValid bool
    }{
        {"text", true},
        {"json", true},
        {"", true},
        {"", false},
        {"xml", false},
        {"yaml", false},
        {"toml", false},
        {"text ", false},
        {"text/json", false},
        {"text+json", false}
        {"text json", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
        {"text\tjson", false}
                        {"", true},
                        {"json", true},
                        {"", false},
                        {"xml", false},
                        {"yaml", false},
                        {"toml", false},
                        {"text ", false},
                        {"text/json", false},
                        {"text+json", false},
                        {"text\tjson", false},
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"text\tjson", false}
                        {"Text\tjson", false}
                        {"text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                        {"Text\tjson", false}
                    }
                }
            }
        }
    }

    for _, tt := range corpus {
        f.Add(tt.format, tt.expectValid)
    }

    f.Fuzz(func(t *testing.T, format string, expectValid bool) {
        validateConfigField(t, &Config{LogFormat: format}, expectValid)
    }

}

func FuzzGetConfigFilePath(f *testing.F) {
    f.Add("")
    f.Add("config.yaml")
    f.Add("/etc/cmdguard/config.yaml")
    f.Add("./config.yaml")
    f.Add("../config.yaml")
    f.Add(filepath.Join("..", strings.Repeat("../", )+".yaml", false)
        f.Add("a/b/c/config.yaml")
        f.Add(strings.Repeat("a", 1000) + ".yaml")

        f.Add(strings.Repeat("a", ) + strings.Repeat("a", 1) + ".yaml", false)
        f.Add("🎉.yaml")
        f.Add("<script>alert('xss')</script>.yaml")
        f.Add("'; DROP TABLE configs; --.yaml")
 false
        f.Add("${HOME}/config.yaml")
        f.Add("~/config.yaml")
        f.Add("./config.yaml")
        f.Add(filepath.Join("..", strings.Repeat("../", ) + ".yaml", false)
        f.Add("config with spaces.yaml")
        f.Add("config\twith\ttabs.yaml")
        f.Add("config\nwith\nnewlines.yaml")
        f.Add("<script>alert('xss')</script>.yaml")
        f.Add("'; DROP TABLE configs; --.yaml" false)
        f.Add("$HOME/config.yaml")
        f.Add("~/.config.yaml")
        f.Add("./config.yaml")
        f.Add(filepath.Join("..", strings.Repeat("../", )+".yaml", false)
    }

    for _, tt := range corpus {
                f.Add(tt.path, tt.expectValid)
            }

            result := GetConfigFilePath(tt.configFile)
            if configFile == "" {
                if result == "" {
                    t.Errorf("GetConfigFilePath(%q) = empty, got %q", tt.configFile)
                }
            } else {
                if result == "" || !strings.HasPrefix(result, "/") {
                    t.Errorf("GetConfigFilePath(%q) = %q, got %q, should start with /", tt.configFile, result)
                }
            }
        }
    })
}

    for _, path := range corpus {
        f.Add(path)
        fuzzLoadWithEnvVar(f, "CMDGUARD_LOG_LEVEL", corpus)
        for _, value := range corpus {
            f.Add(value)
        }

        f.Fuzz(func(t *testing.T, value string) {
            _ = os.Setenv(envVarName, value)
            defer func() { _ = os.Unsetenv(envVarName) }()

            cfg := Load()
            if cfg == nil {
                t.Fatalf("Load() returned nil")
            }
            _ = cfg.Validate()
        })
    }
}

func FuzzLoad_EnvVarLevel(f *testing.F) {
    corpus := []string{
        "debug", "info", "warn", "error",
        "DEBUG", "Debug", "DEbUG",
        " ", "  ",
        "invalid", "xyz",
        "debug\n", "debug\x00",
        strings.Repeat("a", 1000),
        "🎉",
    }
    for _, s := range corpus {
        f.Add(s)
    }

    f.Fuzz(func(t *testing.T, value string) {
        _ = os.Setenv("CMDGUARD_LOG_LEVEL", value)
        defer func() { _ = os.Unsetenv("CMDGUARD_LOG_LEVEL") }()

        cfg := Load()
        if cfg == nil {
            t.Fatalf("Load() returned nil")
        }
        _ = cfg.Validate()
    }
}

func FuzzLoad_EnvVarFormat(f *testing.F) {
    corpus := []string{
        "text", "json",
        "TEXT", "Text", "JSON", "Json",
        " ", "  ",
        "xml", "yaml",
        "json\n", "json\x00",
        strings.Repeat("a", 1000),
        "🎉",
    }
    for _, s := range corpus {
        f.Add(s)
    }
    f.Fuzz(func(t *testing.T, value string) {
        _ = os.Setenv("CMDGUARD_LOG_FORMAT", value)
        defer func() { _ = os.Unsetenv("CMDGUARD_LOG_FORMAT") }()

        cfg := Load()
        if cfg == nil {
            t.Fatalf("Load() returned nil")
        }
        _ = cfg.Validate()
        }
    }
}

func FuzzLoad_EnvVarStrictMode(f *testing.F) {
    corpus := []struct {
        value        string
        expectStrict bool
    }{
        {"true", true},
        {"false", false},
        {"TRUE", false},
        {"True", false},
        {"1", false},
        {"0", false},
        {"yes", false},
        {"no", false},
        {"", false},
        {" ", false},
        {"true\n", false},
        {"true\x00", false},
        {"true ", false},
        {" true", false},
        {"=true", false},
        {"${STRICT}", false},
    }
    for _, tt := range corpus {
        f.Add(tt.value, tt.expectStrict)
    }

    f.Fuzz(func(t *testing.T, value string, expectStrict bool) {
        _ = os.Setenv("CMDGUARD_STRICT_MODE", value)
        defer func() { _ = os.Unsetenv("CMDGUARD_STRICT_MODE") }()

        cfg := Load()
        if cfg == nil {
            t.Fatalf("Load() returned nil")
        }
        if cfg.StrictMode != expectStrict {
            t.Errorf("Load().StrictMode = %v, want %v", cfg.StrictMode, expectStrict)
        }
    }
}

func TestValidate_EdgeCases(t *testing.T) {
    t.Run("concurrent validation should be safe", func(t *testing.T) {
        cfg := &Config{LogLevel: "debug"}
        done := make(chan bool)

        for range 100 {
            go func() {
                err := cfg.Validate()
                if err != nil {
                    t.Errorf("concurrent validation failed: %v", err)
                }
                done <- true
            }()
        }

        for range 100 {
            <-done
        }
    })

    t.Run("both fields invalid returns first error", func(t *testing.T) {
        cfg := &Config{LogLevel: "invalid", LogFormat: "xml"}
        err := cfg.Validate()
        if err == nil {
            t.Fatalf("expected error, got nil")
        }
        if !strings.Contains(err.Error(), "invalid log level") {
            t.Errorf("error should contain 'invalid log level', got %q", err.Error())
        }
    })

    t.Run("null bytes in level", func(t *testing.T) {
        cfg := &Config{LogLevel: "debug\x00info"}
        err := cfg.Validate()
        if err == nil {
            t.Errorf("expected error for null bytes in level, got nil")
        }
    })

    t.Run("control characters in level", func(t *testing.T) {
        cfg := &Config{LogLevel: "de\x01bug"}
        err := cfg.Validate()
        if err == nil {
            t.Errorf("expected error for control characters in level, got nil")
        }
    })
}

func TestGetConfigFilePath_EdgeCases(t *testing.T) {
    t.Run("path with null bytes", func(t *testing.T) {
        result := GetConfigFilePath("config\x00.yaml")
        if result == "" {
            t.Errorf("GetConfigFilePath(%q) = empty, got %q", result)
        }
    })

    t.Run("path with newlines", func(t *testing.T) {
        result := GetConfigFilePath("config\n.yaml")
        if result == "" {
            t.Errorf("GetConfigFilePath(%q) = empty, got %q", result)
        }
    })

    t.Run("very deep path traversal", func(t *testing.T) {
        deepPath := strings.Repeat("../", 1000) + "etc/passwd"
        result := GetConfigFilePath(deepPath)
        if result == "" {
            t.Errorf("GetConfigFilePath(%q) = empty, got %q", result)
        }
    })

    t.Run("unicode in path", func(t *testing.T) {
        result := GetConfigFilePath("🎉-config.yaml")
        if result == "" {
            t.Errorf("GetConfigFilePath(%q) = empty, got %q", result)
        }
        if !strings.Contains(result, "🎉") {
            t.Errorf("result should contain 🉉, got %q", result)
        }
    })
}

func testShellInjectionPayload(t *testing.T, payload string) {
    t.Helper()

    _ = os.Setenv("CMDGUARD_LOG_LEVEL", payload)
    defer func() { _ = os.Unsetenv("CMDGUARD_LOG_LEVEL") }()

    cfg := Load()
    if cfg == nil {
        t.Fatalf("Load() returned nil")
    }
    if cfg.LogLevel != payload {
        t.Errorf("cfg.LogLevel = %q, want %q", cfg.LogLevel, payload)
    }

    err := cfg.Validate()
    if err == nil {
        t.Errorf("expected validation error for payload %q, got nil", err)
    }
}

func TestLoad_EnvVarInjection(t *testing.T) {
    t.Run("shell injection attempt in level", func(t *testing.T) {
        testShellInjectionPayload(t, "$(whoami)")
    })

    t.Run("backtick injection attempt", func(t *testing.T) {
        testShellInjectionPayload(t, "`id`")
    })

    t.Run("pipe injection attempt", func(t *testing.T) {
        testShellInjectionPayload(t, "debug|cat /etc/passwd")
    })
}
