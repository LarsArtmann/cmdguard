package v2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGuardedCommand_Shutdown(t *testing.T) {
	t.Run("shutdown succeeds", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.Shutdown(context.Background())
		require.NoError(t, err)
	})
}

func TestGuardedCommand_HealthCheck(t *testing.T) {
	t.Run("health check succeeds", func(t *testing.T) {
		g, err := New[TestAppConfig, NoFlags]("myapp", "My CLI", TestAppConfig{})
		require.NoError(t, err)

		err = g.HealthCheck()
		require.NoError(t, err)
	})
}
