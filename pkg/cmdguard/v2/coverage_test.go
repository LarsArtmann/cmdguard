package v2

import (
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"charm.land/fang/v2"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestNewCLI(t *testing.T) {
	t.Parallel()

	t.Run("success returns CLI", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[struct{}]("test", "test app", struct{}{})
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})

	t.Run("empty name returns error", func(t *testing.T) {
		t.Parallel()

		_, err := NewCLI[struct{}]("", "test", struct{}{})
		testutil.AssertExpectedError(t, err)
	})
}

func TestAddCommand(t *testing.T) {
	t.Parallel()

	t.Run("success adds command", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[struct{}]("test", "test app", struct{}{})
		testutil.AssertNoError(t, err)
		cmd, err := NewCommand(
			"hello",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
			WithShort("Say hello"),
		)
		testutil.AssertNoError(t, err)

		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("duplicate command returns error", func(t *testing.T) {
		t.Parallel()

		cli, err := NewCLI[struct{}]("test", "test app", struct{}{})
		testutil.AssertNoError(t, err)
		cmd, err := NewCommand(
			"hello",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
		)
		testutil.AssertNoError(t, err)

		testutil.AssertNoError(t, AddCommand(cli, cmd))

		cmd2, _ := NewCommand(
			"hello",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
		)
		err = AddCommand(cli, cmd2)
		testutil.AssertExpectedError(t, err)
	})
}

func TestWithSignalHandling(t *testing.T) {
	t.Parallel()

	t.Run("sets signalHandling flag", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cli, err := NewCLI[cfg]("test", "test", cfg{}, WithSignalHandling[cfg]())
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})
}

func TestWithCLICommit(t *testing.T) {
	t.Parallel()

	t.Run("sets commit via fang option", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cli, err := NewCLI[cfg]("test", "test", cfg{}, WithCLICommit[cfg]("abc123"))
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})
}

func TestWithFangOptions(t *testing.T) {
	t.Parallel()

	t.Run("sets custom fang options", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cli, err := NewCLI[cfg]("test", "test", cfg{}, WithFangOptions[cfg]())
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})
}

func TestWithFangErrorHandler(t *testing.T) {
	t.Parallel()

	t.Run("sets custom error handler", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		handler := func(w io.Writer, styles fang.Styles, err error) {}
		cli, err := NewCLI[cfg]("test", "test", cfg{}, WithFangErrorHandler[cfg](handler))
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})
}

func TestWithFangColorScheme(t *testing.T) {
	t.Parallel()

	t.Run("sets custom color scheme", func(t *testing.T) {
		t.Parallel()

		type cfg struct{}

		cs := fang.DefaultColorScheme
		cli, err := NewCLI[cfg]("test", "test", cfg{}, WithFangColorScheme[cfg](cs))
		testutil.AssertNoError(t, err)
		testutil.AssertNotNil(t, cli)
	})
}

func TestCommandAccessors(t *testing.T) {
	t.Parallel()

	t.Run("Version returns version", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"test",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
			func(s *commandSpec) { s.version = "1.2.3" },
		)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, "1.2.3", cmd.Version())
	})

	t.Run("SilenceErrors returns silenceErrors", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"test",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
			func(s *commandSpec) { s.silenceErrors = true },
		)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, cmd.SilenceErrors(), "should be true")
	})

	t.Run("SilenceUsage returns silenceUsage", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"test",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
			func(s *commandSpec) { s.silenceUsage = true },
		)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, cmd.SilenceUsage(), "should be true")
	})

	t.Run("Group returns group", func(t *testing.T) {
		t.Parallel()

		cmd, err := NewCommand(
			"test",
			NoFlags{},
			func(ctx context.Context, cfg *struct{}, flags NoFlags) error { return nil },
			WithGroupID("mygroup"),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertEqual(t, "mygroup", cmd.Group())
	})
}

func TestBranchWithDuration(t *testing.T) {
	t.Parallel()

	t.Run("creates child with timeout", func(t *testing.T) {
		t.Parallel()

		parent := NewBranchingFlowContext(context.Background())
		child, cancel := parent.BranchWithDuration("child", 5*time.Second)
		defer cancel()

		testutil.AssertNotNil(t, child)
		testutil.AssertEqual(t, "child", child.PathString())
		testutil.AssertBoolTrue(t, child.Context != nil, "context should not be nil")
		testutil.AssertBoolTrue(t, !child.IsRoot(), "child should not be root")
	})
}

func TestBranchWithDeadlineTime(t *testing.T) {
	t.Parallel()

	t.Run("creates child with deadline", func(t *testing.T) {
		t.Parallel()

		parent := NewBranchingFlowContext(context.Background())
		deadline := time.Now().Add(10 * time.Second)
		child, cancel := parent.BranchWithDeadlineTime("child", deadline)
		defer cancel()

		testutil.AssertNotNil(t, child)
		testutil.AssertEqual(t, "child", child.PathString())
	})
}

func TestRegisterFlagValidator(t *testing.T) {
	t.Parallel()

	t.Run("registers custom validator on instance", func(t *testing.T) {
		t.Parallel()

		registry, err := NewFlagRegistry(&struct {
			Name string `flag:"name" validate:"custom"`
		}{})
		testutil.AssertNoError(t, err)

		registry.RegisterFlagValidator("custom", func(value string) error {
			return nil
		})
	})
}

func TestRegisterTypeHandlerOnFlagRegistry(t *testing.T) {
	t.Parallel()

	t.Run("registers custom type on instance", func(t *testing.T) {
		t.Parallel()

		registry, err := NewFlagRegistry(&struct {
			Name string `flag:"name"`
		}{})
		testutil.AssertNoError(t, err)

		registry.RegisterTypeHandler(reflect.TypeFor[string](), TypeHandlerFunc{
			ParseFunc: func(value string, _ FlagTag) (any, error) { return "custom:" + value, nil },
		})

		testutil.AssertNotNil(t, registry)
	})
}

func TestFlagRegistry_RegisterGoDurationHandler(t *testing.T) {
	t.Parallel()

	t.Run("registers time.Duration on instance", func(t *testing.T) {
		t.Parallel()

		registry, err := NewFlagRegistry(&struct {
			Timeout string `flag:"timeout"`
		}{})
		testutil.AssertNoError(t, err)

		registry.RegisterGoDurationHandler()
		testutil.AssertNotNil(t, registry)
	})
}

func TestFilePath_Exists(t *testing.T) {
	t.Parallel()

	t.Run("existing file returns true", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath("/etc/passwd", true)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, fp.Exists(), "/etc/passwd should exist")
	})

	t.Run("nonexistent file returns false", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath("/nonexistent/path/to/file.txt", false)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, !fp.Exists(), "nonexistent file should not exist")
	})
}

func TestFilePath_IsDir(t *testing.T) {
	t.Parallel()

	t.Run("directory returns true", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath(".", true)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, fp.IsDir(), "current dir should be a dir")
	})

	t.Run("file returns false", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath("/etc/passwd", true)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, !fp.IsDir(), "file should not be a dir")
	})
}

func TestFilePath_IsFile(t *testing.T) {
	t.Parallel()

	t.Run("regular file returns true", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath("/etc/passwd", true)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, fp.IsFile(), "/etc/passwd should be a file")
	})

	t.Run("directory returns false", func(t *testing.T) {
		t.Parallel()

		fp, err := ParseFilePath(".", false)
		testutil.AssertNoError(t, err)
		testutil.AssertBoolTrue(t, !fp.IsFile(), "dir should not be a file")
	})
}
