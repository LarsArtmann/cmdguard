// Package main demonstrates error handling patterns with cmdguard v2.
//
// This example shows:
// - Sentinel errors with errors.Is()
// - Wrapped errors with errors.As() for typed context
// - FlagError with typo suggestions
// - ValidationError via PreRunE
// - CommandError wrapping
//
// Usage:
//
//	go run examples/error-handling/main.go fetch --url=https://example.com
//	go run examples/error-handling/main.go fetch --url=invalid
//	go run examples/error-handling/main.go fetch --port=99999
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// AppConfig is the application-level configuration.
type AppConfig struct {
	Verbose bool `default:"false" flag:"verbose" help:"Show detailed errors" short:"v"`
}

// FetchFlags demonstrates error handling with validation.
type FetchFlags struct {
	URL     string `flag:"url"     short:"u" help:"URL to fetch"          default:""`
	Port    int    `flag:"port"    short:"p" help:"Port number (1-65535)" default:"443"`
	Timeout int    `flag:"timeout"           help:"Timeout in seconds"    default:"30"`
}

func main() {
	cli, err := v2.NewCLI[AppConfig]("error-demo", "Error handling demo", AppConfig{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fetchCmd, err := v2.NewCommand[AppConfig, *FetchFlags]("fetch",
		func(_ context.Context, cfg *AppConfig, flags *FetchFlags) error {
			fmt.Printf("Fetching from %s:%d (timeout: %ds)\n", flags.URL, flags.Port, flags.Timeout)

			if cfg.Verbose {
				fmt.Println("[verbose] fetch completed successfully")
			}

			return nil
		},
		v2.WithShort[AppConfig, *FetchFlags]("Fetch a URL with error handling"),
		v2.WithFlags[AppConfig, *FetchFlags](&FetchFlags{}),
		v2.WithPreRunE[AppConfig, *FetchFlags](
			func(_ context.Context, _ *AppConfig, flags *FetchFlags) error {
				var errs []error

				switch flags.URL {
				case "":
					errs = append(errs, errors.New("--url is required"))
				case "invalid":
					errs = append(errs, v2.NewFlagError("url",
						fmt.Errorf("%q is not a valid URL", flags.URL)))
				}

				if flags.Port < 1 || flags.Port > 65535 {
					errs = append(errs, v2.NewFlagError("port",
						fmt.Errorf("port %d out of range (1-65535)", flags.Port)))
				}

				if len(errs) > 0 {
					for _, e := range errs {
						fmt.Fprintf(os.Stderr, "  error: %v\n", e)
					}

					return fmt.Errorf("validation failed with %d error(s)", len(errs))
				}

				return nil
			},
		),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := v2.AddCommand(cli, fetchCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cli.Execute(context.Background())
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "--- Error Analysis ---")

	if errors.Is(err, v2.ErrFlagParseFailed) {
		fmt.Fprintln(os.Stderr, "Type: Flag parse error")
	}

	var flagErr *v2.FlagError
	if errors.As(err, &flagErr) {
		fmt.Fprintf(os.Stderr, "Flag: %s\n", flagErr.FlagName)

		if flagErr.Suggestion != "" {
			fmt.Fprintf(os.Stderr, "Did you mean: --%s?\n", flagErr.Suggestion)
		}
	}

	if cfg := cli.Config(); cfg.Verbose {
		fmt.Fprintf(os.Stderr, "Full error: %v\n", err)
	}

	os.Exit(1)
}
