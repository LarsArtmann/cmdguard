// Package main demonstrates signal handling with cmdguard v2.
//
// This example shows:
// - WithSignalHandling for automatic context cancellation
// - Graceful shutdown on SIGINT/SIGTERM
// - Long-running commands that respect context cancellation
//
// Usage:
//
//	go run examples/signals/main.go serve
//	# Then press Ctrl+C to trigger graceful shutdown
//	go run examples/signals/main.go ping
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
type AppConfig struct {
	Debug bool `default:"false" flag:"debug" help:"Enable debug mode" short:"d"`
}

// ServeFlags defines flags for the serve command.
type ServeFlags struct {
	Port    int `flag:"port"    short:"p" default:"8080" help:"Server port"`
	Workers int `flag:"workers" short:"w" default:"4"   help:"Number of workers"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("signal-demo", "Signal handling demo", AppConfig{},
		v2.WithSignalHandling[AppConfig](),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	serveCmd, err := v2.NewCommand[AppConfig, *ServeFlags]("serve",
		func(ctx context.Context, _ *AppConfig, flags *ServeFlags) error {
			fmt.Printf("Starting server on :%d with %d workers\n", flags.Port, flags.Workers)
			fmt.Println("Press Ctrl+C to trigger graceful shutdown")

			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()

			count := 0
			for {
				select {
				case <-ctx.Done():
					fmt.Println()
					fmt.Println("Signal received, shutting down gracefully...")

					fmt.Println("Draining connections...")
					time.Sleep(500 * time.Millisecond)

					fmt.Println("Stopping workers...")
					time.Sleep(300 * time.Millisecond)

					fmt.Println("Server stopped")

					return nil
				case <-ticker.C:
					count++
					fmt.Printf("\r  serving... (%ds)", count)
				}
			}
		},
		v2.WithShort[AppConfig, *ServeFlags]("Start a long-running server"),
		v2.WithLong[AppConfig, *ServeFlags](
			`Starts a simulated server that handles SIGINT/SIGTERM for graceful shutdown.

The context passed to the handler is cancelled when a signal is received,
allowing clean cleanup of resources.`,
		),
		v2.WithFlags[AppConfig, *ServeFlags](&ServeFlags{}),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, serveCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	pingCmd, err := v2.NewCommand[AppConfig, v2.NoFlags]("ping",
		func(_ context.Context, _ *AppConfig, _ v2.NoFlags) error {
			fmt.Println("pong")
			return nil
		},
		v2.WithShort[AppConfig, v2.NoFlags]("Quick health check"),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, pingCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	cli.ExecuteAndExit(context.Background())
}
