package v4

import (
	"testing"
)

func FuzzParseURL(f *testing.F) {
	f.Add("https://example.com")
	f.Add("http://localhost:8080/path?q=1")
	f.Add("")
	f.Add("not-a-url")
	f.Add("://missing-scheme")
	f.Add("ftp://user:pass@host:21/path")
	f.Add("https://example.com:443/with/path")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParseURL(s)
	})
}

func FuzzParseEmail(f *testing.F) {
	f.Add("user@example.com")
	f.Add("")
	f.Add("not-an-email")
	f.Add("a@b")
	f.Add("user+tag@domain.co.uk")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParseEmail(s)
	})
}

func FuzzParsePort(f *testing.F) {
	f.Add("80")
	f.Add("443")
	f.Add("0")
	f.Add("65536")
	f.Add("")
	f.Add("http")
	f.Add("-1")
	f.Add("abc")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParsePort(s)
	})
}

func FuzzParseHostPort(f *testing.F) {
	f.Add("localhost:8080")
	f.Add("example.com:443")
	f.Add(":8080")
	f.Add("")
	f.Add("no-port")
	f.Add("host:99999")
	f.Add("127.0.0.1:3000")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParseHostPort(s)
	})
}

func FuzzParseDuration(f *testing.F) {
	f.Add("30s")
	f.Add("5m")
	f.Add("1h30m")
	f.Add("")
	f.Add("not-a-duration")
	f.Add("100ms")
	f.Add("-5s")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ParseDuration(s)
	})
}

func FuzzParseFlagTags(f *testing.F) {
	f.Add("valid")
	f.Add("empty")
	f.Add("nil")
	f.Add("string")

	f.Fuzz(func(t *testing.T, typ string) {
		switch typ {
		case "valid":
			type cfg struct {
				Name string `flag:"name" help:"a name" default:"default"`
			}
			_, _ = ParseFlagTags(&cfg{})
		case "empty":
			type cfg struct{}
			_, _ = ParseFlagTags(&cfg{})
		case "nil":
			_, _ = ParseFlagTags(nil)
		default:
			_, _ = ParseFlagTags(typ)
		}
	})
}

func FuzzSetField(f *testing.F) {
	type config struct {
		Name    string
		Port    int
		Enabled bool
	}

	f.Add("Name", "hello")
	f.Add("Port", "42")
	f.Add("Enabled", "true")
	f.Add("Nonexistent", "value")
	f.Add("", "")

	f.Fuzz(func(t *testing.T, fieldName, value string) {
		cfg := &config{}
		_ = SetField(cfg, fieldName, value)
	})
}
