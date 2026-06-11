package v2

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/samber/do/v2"

	"github.com/larsartmann/cmdguard/v2/pkg/testutil"
)

func TestDoctorCommand(t *testing.T) {
	t.Parallel()

	newDoctorCLI := func(t *testing.T) *CLI[testConfig] {
		t.Helper()

		cli, err := NewCLI[testConfig]("myapp", "My app", testConfig{}, WithFang[testConfig](false))
		testutil.AssertNoError(t, err)

		return cli
	}

	// newFailingDoctor builds a doctor command whose only check returns the given error.
	// The command is also registered against the CLI for execution in tests.
	newFailingDoctor := func(t *testing.T, cli *CLI[testConfig], name, errMsg string) {
		t.Helper()

		cmd, err := DoctorCommand[testConfig](
			cli,
			WithDoctorCheck[testConfig](name, func(_ context.Context) error {
				return errors.New(errMsg)
			}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	}

	t.Run("reports no custom checks as healthy", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		cmd, err := DoctorCommand[testConfig](cli)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		testutil.AssertNoError(t, cli.Execute(context.Background()))

		output := out.String()
		if !containsAll(output, "passed", "0 failed") {
			t.Errorf("expected all passed, got: %s", output)
		}
	})

	t.Run("reports healthy DI service", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		testutil.AssertNoError(t, ProvideValue(cli.Scope(), &doctorHealthyService{}))

		cmd, err := DoctorCommand[testConfig](cli)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		testutil.AssertNoError(t, cli.Execute(context.Background()))

		output := out.String()
		if !strings.Contains(output, "✓") {
			t.Errorf("expected checkmark, got: %s", output)
		}

		if !strings.Contains(output, "passed") {
			t.Errorf("expected 'passed', got: %s", output)
		}
	})

	t.Run("reports unhealthy DI service with exit code 1", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		testutil.AssertNoError(t, ProvideValue(cli.Scope(), &doctorUnhealthyService{}))

		cmd, err := DoctorCommand[testConfig](cli)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		err = cli.Execute(context.Background())
		if err == nil {
			t.Fatal("expected error for unhealthy service")
		}

		output := out.String()
		if !strings.Contains(output, "✗") {
			t.Errorf("expected cross mark, got: %s", output)
		}

		if !strings.Contains(output, "failed") {
			t.Errorf("expected failed count, got: %s", output)
		}
	})

	t.Run("runs custom checks after DI health checks", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		customCheckCalled := false

		cmd, err := DoctorCommand[testConfig](
			cli,
			WithDoctorCheck[testConfig]("custom", func(ctx context.Context) error {
				customCheckCalled = true

				return nil
			}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		testutil.AssertNoError(t, cli.Execute(context.Background()))

		if !customCheckCalled {
			t.Error("expected custom check to be called")
		}

		if !strings.Contains(out.String(), "✓ custom") {
			t.Errorf("expected custom check output, got: %s", out.String())
		}
	})

	t.Run("custom check failure causes exit code 1", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)

		newFailingDoctor(t, cli, "failing", "connection refused")

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		err := cli.Execute(context.Background())
		if err == nil {
			t.Fatal("expected error for failing check")
		}

		output := out.String()
		if !strings.Contains(output, "✗ failing") {
			t.Errorf("expected failing check output, got: %s", output)
		}

		if !strings.Contains(output, "connection refused") {
			t.Errorf("expected error message in output, got: %s", output)
		}
	})

	t.Run("sorts results alphabetically", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)

		cmd, err := DoctorCommand[testConfig](
			cli,
			WithDoctorCheck[testConfig]("zebra", func(ctx context.Context) error {
				return nil
			}),
			WithDoctorCheck[testConfig]("alpha", func(ctx context.Context) error {
				return nil
			}),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		testutil.AssertNoError(t, cli.Execute(context.Background()))

		output := out.String()
		alphaIdx := strings.Index(output, "alpha")
		zebraIdx := strings.Index(output, "zebra")

		if alphaIdx > zebraIdx {
			t.Errorf("expected alpha before zebra, got: %s", output)
		}
	})

	t.Run("WithDoctorShort", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		cmd, err := DoctorCommand[testConfig](
			cli,
			WithDoctorShort[testConfig]("Run diagnostics"),
		)
		testutil.AssertNoError(t, err)
		if cmd.Short() != "Run diagnostics" {
			t.Errorf("expected short 'Run diagnostics', got %q", cmd.Short())
		}
	})

	t.Run("WithDoctorGroupID", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		cmd, err := DoctorCommand[testConfig](
			cli,
			WithDoctorGroupID[testConfig]("system"),
		)
		testutil.AssertNoError(t, err)
		testutil.AssertNoError(t, AddCommand(cli, cmd))
	})

	t.Run("mixed healthy and unhealthy checks", func(t *testing.T) {
		t.Parallel()

		cli := newDoctorCLI(t)
		testutil.AssertNoError(t, ProvideValue(cli.Scope(), &doctorHealthyService{}))

		newFailingDoctor(t, cli, "db", "timeout")

		var out strings.Builder
		cli.rootCmd.SetOut(&out)
		cli.rootCmd.SetArgs([]string{"doctor"})

		err := cli.Execute(context.Background())
		if err == nil {
			t.Fatal("expected error")
		}

		output := out.String()
		if !containsAll(output, "✓", "✗") {
			t.Errorf("expected both checkmark and cross, got: %s", output)
		}
	})
}

var (
	_ do.HealthcheckerWithContext = (*doctorHealthyService)(nil)
	_ do.HealthcheckerWithContext = (*doctorUnhealthyService)(nil)
)

type doctorHealthyService struct{}

func (h *doctorHealthyService) HealthCheck(_ context.Context) error { return nil }

type doctorUnhealthyService struct{}

func (u *doctorUnhealthyService) HealthCheck(_ context.Context) error {
	return errors.New("service is broken")
}
