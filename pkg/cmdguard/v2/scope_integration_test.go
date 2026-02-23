package v2

import (
	"context"
	"testing"
	"time"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScope_Integration(t *testing.T) {
	t.Run("full workflow with DI", func(t *testing.T) {
		// Create root scope
		root := NewScope("app")

		// Register services
		type Config struct {
			Debug bool
		}
		require.NoError(t, ProvideValue(root, Config{Debug: true}))

		require.NoError(t, Provide(root, func(i do.Injector) (string, error) {
			cfg, err := do.Invoke[Config](i)
			if err != nil {
				return "", err
			}
			if cfg.Debug {
				return "debug-mode", nil
			}
			return "production-mode", nil
		}))

		// Verify services
		cfg, err := Invoke[Config](root)
		require.NoError(t, err)
		assert.True(t, cfg.Debug)

		mode, err := Invoke[string](root)
		require.NoError(t, err)
		assert.Equal(t, "debug-mode", mode)

		// Health check
		require.NoError(t, root.HealthCheck())

		// Shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		require.NoError(t, root.Shutdown(ctx))
	})

	t.Run("child scope can override parent services", func(t *testing.T) {
		assertChildInheritsParent(t)
	})
}
