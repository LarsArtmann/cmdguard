package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

func TestURL(t *testing.T) {
	t.Parallel()

	t.Run("ParseURL valid", func(t *testing.T) {
		t.Parallel()
		u, err := v2.ParseURL("https://example.com:8080/path")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.String() != "https://example.com:8080/path" {
			t.Errorf("String() = %q, want %q", u.String(), "https://example.com:8080/path")
		}
		if u.Scheme() != "https" {
			t.Errorf("Scheme() = %q, want %q", u.Scheme(), "https")
		}
		if u.Host() != "example.com:8080" {
			t.Errorf("Host() = %q, want %q", u.Host(), "example.com:8080")
		}
		if u.Hostname() != "example.com" {
			t.Errorf("Hostname() = %q, want %q", u.Hostname(), "example.com")
		}
		if u.Port() != "8080" {
			t.Errorf("Port() = %q, want %q", u.Port(), "8080")
		}
		if u.Path() != "/path" {
			t.Errorf("Path() = %q, want %q", u.Path(), "/path")
		}
	})

	t.Run("ParseURL empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseURL("")
		if err == nil {
			t.Fatal("expected error for empty URL")
		}
	})

	t.Run("ParseURL missing scheme", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseURL("example.com/path")
		if err == nil {
			t.Fatal("expected error for URL without scheme")
		}
	})

	t.Run("ParseURL missing host", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseURL("http:///path")
		if err == nil {
			t.Fatal("expected error for URL without host")
		}
	})

	t.Run("MustParseURL valid", func(t *testing.T) {
		t.Parallel()
		u := v2.MustParseURL("https://example.com")
		if u.String() != "https://example.com" {
			t.Errorf("String() = %q, want %q", u.String(), "https://example.com")
		}
	})

	t.Run("MustParseURL panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid URL")
			}
		}()
		v2.MustParseURL("")
	})

	t.Run("URL IsEmpty", func(t *testing.T) {
		t.Parallel()
		u, _ := v2.ParseURL("https://example.com")
		if u.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("URL URL returns copy", func(t *testing.T) {
		t.Parallel()
		u, _ := v2.ParseURL("https://example.com")
		url1 := u.URL()
		url2 := u.URL()
		if url1 == url2 {
			t.Error("URL() should return different pointers")
		}
	})

	t.Run("URL MarshalText", func(t *testing.T) {
		t.Parallel()
		u, _ := v2.ParseURL("https://example.com/path")
		data, err := u.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "https://example.com/path" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "https://example.com/path")
		}
	})

	t.Run("URL UnmarshalText", func(t *testing.T) {
		t.Parallel()
		var u v2.URL
		err := u.UnmarshalText([]byte("https://example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if u.String() != "https://example.com" {
			t.Errorf("String() = %q, want %q", u.String(), "https://example.com")
		}
	})
}

func TestEmail(t *testing.T) {
	t.Parallel()

	t.Run("ParseEmail valid", func(t *testing.T) {
		t.Parallel()
		e, err := v2.ParseEmail("user@example.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if e.String() != "user@example.com" {
			t.Errorf("String() = %q, want %q", e.String(), "user@example.com")
		}
		if e.Address() != "user@example.com" {
			t.Errorf("Address() = %q, want %q", e.Address(), "user@example.com")
		}
		if e.Local() != "user" {
			t.Errorf("Local() = %q, want %q", e.Local(), "user")
		}
		if e.Domain() != "example.com" {
			t.Errorf("Domain() = %q, want %q", e.Domain(), "example.com")
		}
	})

	t.Run("ParseEmail empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseEmail("")
		if err == nil {
			t.Fatal("expected error for empty email")
		}
	})

	t.Run("ParseEmail invalid", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseEmail("not-an-email")
		if err == nil {
			t.Fatal("expected error for invalid email")
		}
	})

	t.Run("MustParseEmail valid", func(t *testing.T) {
		t.Parallel()
		e := v2.MustParseEmail("test@test.com")
		if e.String() != "test@test.com" {
			t.Errorf("String() = %q, want %q", e.String(), "test@test.com")
		}
	})

	t.Run("MustParseEmail panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid email")
			}
		}()
		v2.MustParseEmail("invalid")
	})

	t.Run("Email IsEmpty", func(t *testing.T) {
		t.Parallel()
		e, _ := v2.ParseEmail("user@example.com")
		if e.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("Email MarshalText", func(t *testing.T) {
		t.Parallel()
		e, _ := v2.ParseEmail("user@example.com")
		data, err := e.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "user@example.com" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "user@example.com")
		}
	})

	t.Run("Email UnmarshalText", func(t *testing.T) {
		t.Parallel()
		var e v2.Email
		err := e.UnmarshalText([]byte("test@example.com"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if e.String() != "test@example.com" {
			t.Errorf("String() = %q, want %q", e.String(), "test@example.com")
		}
	})
}

func TestPort(t *testing.T) {
	t.Parallel()

	t.Run("ParsePort numeric", func(t *testing.T) {
		t.Parallel()
		p, err := v2.ParsePort("8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}
		if p.String() != "8080" {
			t.Errorf("String() = %q, want %q", p.String(), "8080")
		}
	})

	t.Run("ParsePort named http", func(t *testing.T) {
		t.Parallel()
		p, err := v2.ParsePort("http")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 80 {
			t.Errorf("Int() = %d, want %d", p.Int(), 80)
		}
	})

	t.Run("ParsePort named https", func(t *testing.T) {
		t.Parallel()
		p, err := v2.ParsePort("https")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 443 {
			t.Errorf("Int() = %d, want %d", p.Int(), 443)
		}
	})

	t.Run("ParsePort named ssh", func(t *testing.T) {
		t.Parallel()
		p, err := v2.ParsePort("ssh")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 22 {
			t.Errorf("Int() = %d, want %d", p.Int(), 22)
		}
	})

	t.Run("ParsePort empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParsePort("")
		if err == nil {
			t.Fatal("expected error for empty port")
		}
	})

	t.Run("ParsePort out of range", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParsePort("70000")
		if err == nil {
			t.Fatal("expected error for out of range port")
		}
	})

	t.Run("ParsePort zero", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParsePort("0")
		if err == nil {
			t.Fatal("expected error for port 0")
		}
	})

	t.Run("ParsePort negative", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParsePort("-1")
		if err == nil {
			t.Fatal("expected error for negative port")
		}
	})

	t.Run("PortFromInt valid", func(t *testing.T) {
		t.Parallel()
		p, err := v2.PortFromInt(8080)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 8080 {
			t.Errorf("Int() = %d, want %d", p.Int(), 8080)
		}
	})

	t.Run("PortFromInt invalid", func(t *testing.T) {
		t.Parallel()
		_, err := v2.PortFromInt(70000)
		if err == nil {
			t.Fatal("expected error for out of range port")
		}
	})

	t.Run("Port IsValid", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
		if !p.IsValid() {
			t.Error("IsValid() = false, want true")
		}
	})

	t.Run("Port IsWellKnown", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("80")
		if !p.IsWellKnown() {
			t.Error("IsWellKnown() = false, want true for port 80")
		}
		p2, _ := v2.ParsePort("8080")
		if p2.IsWellKnown() {
			t.Error("IsWellKnown() = true, want false for port 8080")
		}
	})

	t.Run("Port IsRegistered", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
		if !p.IsRegistered() {
			t.Error("IsRegistered() = false, want true for port 8080")
		}
	})

	t.Run("Port IsDynamic", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("50000")
		if !p.IsDynamic() {
			t.Error("IsDynamic() = false, want true for port 50000")
		}
	})

	t.Run("MustParsePort valid", func(t *testing.T) {
		t.Parallel()
		p := v2.MustParsePort("443")
		if p.Int() != 443 {
			t.Errorf("Int() = %d, want %d", p.Int(), 443)
		}
	})

	t.Run("MustParsePort panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid port")
			}
		}()
		v2.MustParsePort("invalid")
	})

	t.Run("Port MarshalText", func(t *testing.T) {
		t.Parallel()
		p, _ := v2.ParsePort("8080")
		data, err := p.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "8080" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "8080")
		}
	})

	t.Run("Port UnmarshalText", func(t *testing.T) {
		t.Parallel()
		var p v2.Port
		err := p.UnmarshalText([]byte("9090"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if p.Int() != 9090 {
			t.Errorf("Int() = %d, want %d", p.Int(), 9090)
		}
	})
}

func TestFilePath(t *testing.T) {
	t.Parallel()

	t.Run("ParseFilePath relative", func(t *testing.T) {
		t.Parallel()
		fp, err := v2.ParseFilePath("./test.txt", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fp.String() != "test.txt" {
			t.Errorf("String() = %q, want %q", fp.String(), "test.txt")
		}
		if fp.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("ParseFilePath empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseFilePath("", false)
		if err == nil {
			t.Fatal("expected error for empty path")
		}
	})

	t.Run("FilePath Absolute not empty", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("test.txt", false)
		if fp.Absolute() == "" {
			t.Error("Absolute() should not be empty")
		}
	})

	t.Run("FilePath Dir", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("/home/user/file.txt", false)
		if fp.Dir() == "" {
			t.Error("Dir() should not be empty")
		}
	})

	t.Run("FilePath Base", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("/home/user/file.txt", false)
		if fp.Base() != "file.txt" {
			t.Errorf("Base() = %q, want %q", fp.Base(), "file.txt")
		}
	})

	t.Run("FilePath Ext", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("/home/user/file.txt", false)
		if fp.Ext() != ".txt" {
			t.Errorf("Ext() = %q, want %q", fp.Ext(), ".txt")
		}
	})

	t.Run("MustParseFilePath valid", func(t *testing.T) {
		t.Parallel()
		fp := v2.MustParseFilePath("/tmp/test.txt", false)
		if fp.Base() != "test.txt" {
			t.Errorf("Base() = %q, want %q", fp.Base(), "test.txt")
		}
	})

	t.Run("MustParseFilePath panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for empty path")
			}
		}()
		v2.MustParseFilePath("", false)
	})

	t.Run("FilePath MarshalText", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("/tmp/test.txt", false)
		data, err := fp.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "/tmp/test.txt" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "/tmp/test.txt")
		}
	})

	t.Run("FilePath UnmarshalText", func(t *testing.T) {
		t.Parallel()
		var fp v2.FilePath
		err := fp.UnmarshalText([]byte("/tmp/test.txt"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fp.String() != "/tmp/test.txt" {
			t.Errorf("String() = %q, want %q", fp.String(), "/tmp/test.txt")
		}
	})

	t.Run("FilePath Join", func(t *testing.T) {
		t.Parallel()
		fp, _ := v2.ParseFilePath("/tmp", false)
		joined := fp.Join("subdir", "file.txt")
		if joined.Base() != "file.txt" {
			t.Errorf("Base() = %q, want %q", joined.Base(), "file.txt")
		}
	})
}

func TestHostPort(t *testing.T) {
	t.Parallel()

	t.Run("ParseHostPort valid", func(t *testing.T) {
		t.Parallel()
		hp, err := v2.ParseHostPort("localhost:8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hp.String() != "localhost:8080" {
			t.Errorf("String() = %q, want %q", hp.String(), "localhost:8080")
		}
		if hp.Host() != "localhost" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "localhost")
		}
		if hp.Port().Int() != 8080 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 8080)
		}
	})

	t.Run("ParseHostPort any host", func(t *testing.T) {
		t.Parallel()
		hp, err := v2.ParseHostPort(":8080")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hp.Host() != "" {
			t.Errorf("Host() = %q, want empty", hp.Host())
		}
		if !hp.IsAnyHost() {
			t.Error("IsAnyHost() = false, want true")
		}
	})

	t.Run("ParseHostPort empty", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseHostPort("")
		if err == nil {
			t.Fatal("expected error for empty host:port")
		}
	})

	t.Run("ParseHostPort invalid", func(t *testing.T) {
		t.Parallel()
		_, err := v2.ParseHostPort("not-valid")
		if err == nil {
			t.Fatal("expected error for invalid host:port")
		}
	})

	t.Run("NewHostPort valid", func(t *testing.T) {
		t.Parallel()
		hp, err := v2.NewHostPort("example.com", "443")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hp.Host() != "example.com" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "example.com")
		}
		if hp.Port().Int() != 443 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 443)
		}
	})

	t.Run("MustParseHostPort valid", func(t *testing.T) {
		t.Parallel()
		hp := v2.MustParseHostPort("127.0.0.1:3000")
		if hp.Host() != "127.0.0.1" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "127.0.0.1")
		}
		if hp.Port().Int() != 3000 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 3000)
		}
	})

	t.Run("MustParseHostPort panic", func(t *testing.T) {
		t.Parallel()
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic for invalid host:port")
			}
		}()
		v2.MustParseHostPort("invalid")
	})

	t.Run("HostPort IsEmpty", func(t *testing.T) {
		t.Parallel()
		hp, _ := v2.ParseHostPort("localhost:8080")
		if hp.IsEmpty() {
			t.Error("IsEmpty() = true, want false")
		}
	})

	t.Run("HostPort MarshalText", func(t *testing.T) {
		t.Parallel()
		hp, _ := v2.ParseHostPort("example.com:443")
		data, err := hp.MarshalText()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if string(data) != "example.com:443" {
			t.Errorf("MarshalText() = %q, want %q", string(data), "example.com:443")
		}
	})

	t.Run("HostPort UnmarshalText", func(t *testing.T) {
		t.Parallel()
		var hp v2.HostPort
		err := hp.UnmarshalText([]byte("localhost:9090"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hp.Host() != "localhost" {
			t.Errorf("Host() = %q, want %q", hp.Host(), "localhost")
		}
		if hp.Port().Int() != 9090 {
			t.Errorf("Port().Int() = %d, want %d", hp.Port().Int(), 9090)
		}
	})
}
