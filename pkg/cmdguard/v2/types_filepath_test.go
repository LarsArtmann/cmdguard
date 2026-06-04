package v2_test

import (
	"testing"

	v2 "github.com/larsartmann/cmdguard/v2/pkg/cmdguard/v2"
	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestFilePath(t *testing.T) {
	t.Parallel()

	t.Run("ParseFilePath relative", func(t *testing.T) {
		t.Parallel()

		fp, err := v2.ParseFilePath("./test.txt", false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		testutil.AssertStringerEq(t, fp, "test.txt")

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

		testutil.AssertStringerEq(t, fp, "/tmp/test.txt")
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
