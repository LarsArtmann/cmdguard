// Package main demonstrates input validation patterns with cmdguard v2.
//
// This example shows:
// - PreRunE for validating flag values before command execution
// - Validation with clear error messages and exit codes
// - Severity-based validation (errors block execution, warnings proceed)
// - Reusable validation helpers for common patterns
//
// Validation is a first-class concern in cmdguard v2. Use PreRunE to validate
// user input, flag combinations, and prerequisites before the main handler runs.
//
// Usage:
//
//	go run examples/validation/main.go greet --name=Alice
//	go run examples/validation/main.go greet --name= --count=5
//	go run examples/validation/main.go process --input=file.txt --workers=10
//	go run examples/validation/main.go process --input=file.txt --workers=100
package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	examplesinternal "github.com/larsartmann/cmdguard/examples/internal"
	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

// Config is the application-level configuration.
type Config struct {
	Verbose bool   `default:"false" flag:"verbose"  help:"Enable verbose output"          short:"v"`
	LogFile string `default:""      flag:"log-file" help:"Log file path (empty = stdout)"`
}

// ValidateName validates the name flag.
func ValidateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.New("name cannot be empty")
	}

	if len(name) < 2 {
		return errors.New("name must be at least 2 characters")
	}

	if len(name) > 50 {
		return errors.New("name must be at most 50 characters")
	}

	return nil
}

// ValidateCount validates the count flag.
func ValidateCount(count int) error {
	if count < 1 {
		return errors.New("count must be at least 1")
	}

	if count > 10 {
		return errors.New("count must be at most 10 (use --verbose for more)")
	}

	return nil
}

// ValidateEmail performs a basic email format check.
func ValidateEmail(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	if !strings.Contains(email, "@") {
		return errors.New("email must contain @")
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("email format is invalid")
	}

	if parts[0] == "" || parts[1] == "" {
		return errors.New("email format is invalid")
	}

	return nil
}

// ValidationError wraps validation failures with context.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// addValidationError appends a validation error to the list if err is not nil.
func addValidationError(errs []error, err error, field string) []error {
	if err != nil {
		return append(errs, &ValidationError{Field: field, Message: err.Error()})
	}

	return errs
}

// ValidateFlags performs comprehensive flag validation.
func ValidateFlags(name string, count int, email string) []error {
	errs := addValidationError(nil, ValidateName(name), "name")
	errs = addValidationError(errs, ValidateCount(count), "count")
	errs = addValidationError(errs, ValidateEmail(email), "email")

	return errs
}

// fatal prints the error to stderr and exits with code 1.
func fatal(format string, args ...any) {
	examplesinternal.Fatalf(format, args...)
}

func main() {
	fmt.Println("=== Validation Example (v2 API) ===")
	fmt.Println()

	cli, err := v2.NewCLI[Config]("validation-example", "Validation demonstration", Config{})
	if err != nil {
		fatal("Failed to create CLI: %v\n", err)
	}

	// Add commands
	if err := addCommands(cli); err != nil {
		fatal("Failed to add commands: %v\n", err)
	}

	// Execute
	examplesinternal.Execute(context.Background(), cli)
}

// GreetFlags defines flags for the greet command.
type GreetFlags struct {
	Name      string `default:"World" flag:"name"  help:"Name to greet"               short:"n"`
	Count     int    `default:"1"     flag:"count" help:"Number of times to greet"    short:"c"`
	Email     string `default:""      flag:"email" help:"Email for notification"      short:"e"`
	Uppercase bool   `default:"false" flag:"upper" help:"Print greeting in uppercase"`
}

// greetRunE is the greet command handler.
func greetRunE() func(context.Context, *Config, *GreetFlags) error {
	return func(_ context.Context, _ *Config, flags *GreetFlags) error {
		msg := fmt.Sprintf("Hello, %s!", flags.Name)

		if flags.Uppercase {
			msg = strings.ToUpper(msg)
		}

		for i := 0; i < flags.Count; i++ {
			fmt.Println(msg)
		}

		if flags.Email != "" {
			fmt.Printf("[Would notify: %s]\n", flags.Email)
		}

		return nil
	}
}

// greetPreRunE validates the greet command flags.
func greetPreRunE() func(context.Context, *Config, *GreetFlags) error {
	return func(_ context.Context, _ *Config, flags *GreetFlags) error {
		if errs := ValidateFlags(flags.Name, flags.Count, flags.Email); len(errs) > 0 {
			fmt.Fprintln(os.Stderr, "Validation failed:")

			for _, err := range errs {
				fmt.Fprintf(os.Stderr, "  - %v\n", err)
			}

			return errors.New("validation failed")
		}

		return nil
	}
}

// ProcessFlags defines flags for the process command.
type ProcessFlags struct {
	InputFile string `flag:"input"   help:"Input file path"              required:"true" short:"i"`
	OutputDir string `flag:"output"  help:"Output directory"                             short:"o" default:"./output"`
	Workers   int    `flag:"workers" help:"Number of parallel workers"                   short:"w" default:"1"`
	DryRun    bool   `flag:"dry-run" help:"Show what would be processed"`
}

// processRunE is the process command handler.
func processRunE() func(context.Context, *Config, *ProcessFlags) error {
	return func(_ context.Context, _ *Config, flags *ProcessFlags) error {
		fmt.Printf("Processing %s with %d workers\n", flags.InputFile, flags.Workers)
		fmt.Printf("Output directory: %s\n", flags.OutputDir)

		if flags.DryRun {
			fmt.Println("[Dry run - no actual processing]")
		} else {
			fmt.Println("[Processing complete]")
		}

		return nil
	}
}

// processPreRunE validates the process command flags.
func processPreRunE() func(context.Context, *Config, *ProcessFlags) error {
	return func(_ context.Context, cfg *Config, flags *ProcessFlags) error {
		if flags.InputFile == "" {
			return errors.New("--input is required")
		}

		if flags.Workers < 1 {
			return errors.New("--workers must be at least 1")
		}

		if flags.Workers > 50 {
			return errors.New("--workers must be at most 50 (resource limit)")
		}

		if cfg.Verbose {
			fmt.Printf("[DEBUG] Validating input file: %s\n", flags.InputFile)
		}

		return nil
	}
}

// ConfigFlags defines flags for the config command.
type ConfigFlags struct {
	Key   string `flag:"key"   help:"Configuration key"                   required:"true" short:"k"`
	Value string `flag:"value" help:"Configuration value"                                 short:"v"`
	Get   bool   `flag:"get"   help:"Get config value"                                    short:"g"`
	Set   bool   `flag:"set"   help:"Set config value (requires --value)"`
}

// configRunE is the config command handler.
func configRunE() func(context.Context, *Config, *ConfigFlags) error {
	return func(_ context.Context, _ *Config, flags *ConfigFlags) error {
		if flags.Get {
			fmt.Printf("config.%s = <stored-value>\n", flags.Key)

			return nil
		}

		if flags.Set {
			if flags.Value == "" {
				return errors.New("--value is required when using --set")
			}

			fmt.Printf("config.%s = %s\n", flags.Key, flags.Value)

			return nil
		}

		return errors.New("use --get or --set")
	}
}

// configPreRunE validates the config command flags.
func configPreRunE() func(context.Context, *Config, *ConfigFlags) error {
	return func(_ context.Context, _ *Config, flags *ConfigFlags) error {
		if flags.Key == "" {
			return errors.New("--key is required")
		}

		// Validate key format (alphanumeric and dots only)
		for _, c := range flags.Key {
			if !strings.ContainsRune("abcdefghijklmnopqrstuvwxyz0123456789_.", c) {
				return fmt.Errorf("--key contains invalid character: %c", c)
			}
		}

		return nil
	}
}

// addCommands adds all example commands to the CLI.
func addCommands(cli *v2.CLI[Config]) error {
	// Greet command with validation
	err := v2.AddCommand(cli, v2.Command[Config, *GreetFlags]{
		Use:   "greet",
		Short: "Greet someone",
		Long:  "Greets the named person a specified number of times.",
		Example: `
  # Greet Alice twice
  greet --name Alice --count 2

  # Greet in uppercase
  greet --name Bob --upper`,
		Flags:   &GreetFlags{},
		PreRunE: greetPreRunE(),
		RunE:    greetRunE(),
	})
	if err != nil {
		return fmt.Errorf("adding greet command: %w", err)
	}

	// Process command with validation
	err = v2.AddCommand(cli, v2.Command[Config, *ProcessFlags]{
		Use:   "process",
		Short: "Process a file",
		Long:  "Processes an input file with configurable parallelism.",
		Example: `
  # Process a file with 4 workers
  process --input data.txt --workers 4

  # Dry run to preview
  process --input data.txt --dry-run`,
		Flags:   &ProcessFlags{},
		PreRunE: processPreRunE(),
		RunE:    processRunE(),
	})
	if err != nil {
		return fmt.Errorf("adding process command: %w", err)
	}

	// Config command with validation
	err = v2.AddCommand(cli, v2.Command[Config, *ConfigFlags]{
		Use:   "config",
		Short: "Get or set configuration",
		Long:  "Manages application configuration.",
		Example: `
  # Get a config value
  config --get --key database.host

  # Set a config value
  config --set --key database.host --value localhost`,
		Flags:   &ConfigFlags{},
		PreRunE: configPreRunE(),
		RunE:    configRunE(),
	})
	if err != nil {
		return fmt.Errorf("adding config command: %w", err)
	}

	return nil
}

// =============================================================================
// VALIDATION PATTERNS REFERENCE
// =============================================================================
//
// cmdguard v2 supports validation through the PreRunE hook. This example
// demonstrates common validation patterns:
//
// PATTERN 1: Simple field validation
// ---------------------------------
// func validateName(name string) error {
//     if name == "" {
//         return errors.New("name is required")
//     }
//     return nil
// }
//
// PATTERN 2: Composite validation
// ------------------------------
// func validateAll(name string, count int) []error {
//     var errs []error
//     if err := validateName(name); err != nil {
//         errs = append(errs, err)
//     }
//     if count < 1 || count > 100 {
//         errs = append(errs, errors.New("count out of range"))
//     }
//     return errs
// }
//
// PATTERN 3: Integration with validation libraries
// -------------------------------------------------
// PreRunE can call any validation library:
//
//   PreRunE: func(ctx, cfg, flags) error {
//       // Using github.com/artmann/businessrules
//       result := validator.Validate(flags)
//       if result.HasErrors() {
//           return fmt.Errorf("validation: %w", result.Error())
//       }
//       return nil
//   }
//
// PATTERN 4: Severity-based validation
// -------------------------------------
// Use warnings for non-blocking issues:
//
//   PreRunE: func(ctx, cfg, flags) error {
//       if flags.Deprecated {
//           fmt.Fprintf(os.Stderr, "WARNING: using deprecated feature\n")
//       }
//       if err := validateRequired(flags); err != nil {
//           return err  // This blocks
//       }
//       return nil
//   }
//
// =============================================================================
