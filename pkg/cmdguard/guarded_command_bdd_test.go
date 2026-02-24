package cmdguard_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
)

var _ = Describe("GuardedCommand - User Expectations", func() {
	var root *cmdguard.GuardedCommand

	BeforeEach(func() {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("As a CLI developer creating a new application", func() {
		Context("when I create a new GuardedCommand with basic settings", func() {
			BeforeEach(func() {
				root = cmdguard.New("myapp", "My awesome CLI tool")
			})

			It("should have my application name for help output", func() {
				Expect(root.Command().Name()).To(Equal("myapp"))
			})

			It("should show my description when users run --help", func() {
				Expect(root.Command().Short).To(Equal("My awesome CLI tool"))
			})

			assertBuiltinCommandExists := func(cmdName string) {
				cmd := root.Command()
				foundCmd, _, err := cmd.Find([]string{cmdName})
				Expect(err).ToNot(HaveOccurred())
				Expect(foundCmd).ToNot(BeNil())
				Expect(foundCmd.Name()).To(Equal(cmdName))
			}

			It("should provide version out of the box without me writing code", func() {
				assertBuiltinCommandExists("version")
			})

			It("should provide a validate command to check my CLI setup", func() {
				assertBuiltinCommandExists("validate")
			})
		})

		Context("when I want to configure behavior via environment variables", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
				root = cmdguard.New("myapp", "My CLI")
			})

			AfterEach(func() {
				_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
			})

			It("should pick up strict mode setting automatically", func() {
				Expect(root.IsStrictMode()).To(BeTrue())
			})
		})
	})

	Describe("As a library user adding commands to my CLI", func() {
		BeforeEach(func() {
			root = cmdguard.New("myapp", "My CLI")
		})

		Context("when I add a properly implemented command", func() {
			It("should accept commands with error-returning handlers (RunE)", func() {
				cmd := &cobra.Command{
					Use:   "deploy",
					Short: "Deploy the application",
					RunE:  func(cmd *cobra.Command, args []string) error { return nil },
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should accept commands with simple handlers (Run)", func() {
				cmd := &cobra.Command{
					Use:   "status",
					Short: "Show status",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should accept command groups with subcommands (like 'git remote add')", func() {
				parent := &cobra.Command{
					Use:   "remote",
					Short: "Manage remotes",
				}
				child := &cobra.Command{
					Use:   "add",
					Short: "Add a remote",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				parent.AddCommand(child)
				Expect(func() { root.AddCommand(parent) }).NotTo(Panic())
			})
		})

		Context("when I accidentally forget to implement a command", func() {
			It("should fail fast with a clear error if I forget the handler", func() {
				cmd := &cobra.Command{
					Use:   "incomplete",
					Short: "I forgot to add Run or RunE",
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})

			It("should tell me this is a cmdguard validation failure", func() {
				cmd := &cobra.Command{
					Use:   "incomplete",
					Short: "Missing handler",
				}
				panicFn := func() { root.AddCommand(cmd) }
				Expect(panicFn).To(PanicWith(ContainSubstring("cmdguard:")))
			})

			It("should fail fast if I forget the command name", func() {
				cmd := &cobra.Command{
					Short: "I forgot Use field",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})
		})

		Context("when I add subcommands to an existing command", func() {
			var parent *cobra.Command

			BeforeEach(func() {
				parent = &cobra.Command{
					Use:   "db",
					Short: "Database operations",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				root.AddCommand(parent)
			})

			It("should accept valid subcommands", func() {
				child := &cobra.Command{
					Use:   "migrate",
					Short: "Run migrations",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddSubcommand(parent, child) }).NotTo(Panic())
			})

			It("should catch incomplete subcommands", func() {
				child := &cobra.Command{
					Use:   "broken",
					Short: "Missing handler",
				}
				Expect(func() { root.AddSubcommand(parent, child) }).To(Panic())
			})
		})
	})

	Describe("As a security-conscious operator in strict environments", func() {
		BeforeEach(func() {
			_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
			root = cmdguard.New("myapp", "My CLI")
		})

		AfterEach(func() {
			_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
		})

		Context("when strict mode is enabled", func() {
			It("should require proper error handling (RunE) for all commands", func() {
				cmd := &cobra.Command{
					Use:   "deploy",
					Short: "Production deployment",
					RunE:  func(cmd *cobra.Command, args []string) error { return nil },
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should reject commands that silently swallow errors (Run)", func() {
				cmd := &cobra.Command{
					Use:   "unsafe",
					Short: "This could hide errors",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})

			It("should explain why Run handlers are rejected in strict mode", func() {
				cmd := &cobra.Command{
					Use:   "unsafe",
					Short: "No error handling",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				panicFn := func() { root.AddCommand(cmd) }
				Expect(panicFn).To(PanicWith(ContainSubstring("strict mode")))
			})
		})
	})

	Describe("As a user running the CLI application", func() {
		BeforeEach(func() {
			root = cmdguard.New("myapp", "My CLI")
		})

		Context("when I run built-in commands", func() {
			It("should execute 'version' successfully", func() {
				root.Command().SetArgs([]string{"version"})
				err := root.Execute(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})

			It("should execute 'validate' successfully", func() {
				root.Command().SetArgs([]string{"validate"})
				err := root.Execute(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when I configure logging via environment", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_LOG_FORMAT", "json")
			})

			AfterEach(func() {
				_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
			})

			It("should use my configured format", func() {
				root = cmdguard.New("myapp", "My CLI")
				Expect(root.Config().LogFormat).To(Equal("json"))
			})
		})
	})

	Describe("As a developer needing to inspect configuration", func() {
		BeforeEach(func() {
			root = cmdguard.New("myapp", "My CLI")
		})

		Context("when I need programmatic access to settings", func() {
			It("should expose the configuration object", func() {
				cfg := root.Config()
				Expect(cfg).ToNot(BeNil())
			})

			It("should expose the underlying Cobra command for advanced usage", func() {
				cmd := root.Command()
				Expect(cmd).ToNot(BeNil())
				Expect(cmd.Name()).To(Equal("myapp"))
			})
		})
	})

	Describe("As a developer using ExecuteAndExit for main() simplicity", func() {
		BeforeEach(func() {
			root = cmdguard.New("myapp", "My CLI")
		})

		Context("when the command succeeds", func() {
			It("should complete without calling os.Exit (no panic)", func() {
				root.Command().SetArgs([]string{"version"})
				Expect(func() { root.ExecuteAndExit(context.Background()) }).NotTo(Panic())
			})
		})
	})

	Describe("As a developer making mistakes after initialization", func() {
		Context("when I try to add commands after the CLI has started", func() {
			It("should prevent runtime modification of commands", func() {
				root = cmdguard.New("myapp", "My CLI")
				root.Command().SetArgs([]string{"version"})
				_ = root.Execute(context.Background())

				lateCmd := &cobra.Command{
					Use:   "too-late",
					Short: "Added after execution",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(lateCmd) }).To(PanicWith(ContainSubstring("cannot add commands after execution")))
			})

			It("should also prevent adding subcommands after execution", func() {
				root = cmdguard.New("myapp", "My CLI")
				parent := &cobra.Command{
					Use:   "db",
					Short: "Database commands",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				root.AddCommand(parent)
				root.Command().SetArgs([]string{"version"})
				_ = root.Execute(context.Background())

				child := &cobra.Command{
					Use:   "migrate",
					Short: "Too late to add",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddSubcommand(parent, child) }).To(PanicWith(ContainSubstring("cannot add commands after execution")))
			})
		})
	})

	Describe("As a developer embedding cmdguard in my application", func() {
		Context("when I call the Version function", func() {
			It("should return a version string for my --version output", func() {
				version := cmdguard.Version()
				Expect(version).ToNot(BeEmpty())
			})

			It("should return 'dev' when not built with ldflags", func() {
				version := cmdguard.Version()
				Expect(version).To(Equal("dev"))
			})
		})
	})
})
