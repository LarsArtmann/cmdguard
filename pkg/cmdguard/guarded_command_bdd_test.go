package cmdguard_test

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/pkg/cmdguard"
	"github.com/spf13/cobra"
)

var _ = Describe("GuardedCommand", func() {
	var root *cmdguard.GuardedCommand

	BeforeEach(func() {
		os.Unsetenv("CMDGUARD_LOG_LEVEL")
		os.Unsetenv("CMDGUARD_LOG_FORMAT")
		os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("Creating a GuardedCommand", func() {
		Context("with default settings", func() {
			BeforeEach(func() {
				root = cmdguard.New("testapp", "Test application")
			})

			It("should create a command with the given name", func() {
				Expect(root.Command().Name()).To(Equal("testapp"))
			})

			It("should create a command with the given description", func() {
				Expect(root.Command().Short).To(Equal("Test application"))
			})

			It("should not be in strict mode by default", func() {
				Expect(root.IsStrictMode()).To(BeFalse())
			})

			It("should have built-in version command", func() {
				cmd := root.Command()
				versionCmd, _, err := cmd.Find([]string{"version"})
				Expect(err).ToNot(HaveOccurred())
				Expect(versionCmd).ToNot(BeNil())
				Expect(versionCmd.Name()).To(Equal("version"))
			})

			It("should have built-in validate command", func() {
				cmd := root.Command()
				validateCmd, _, err := cmd.Find([]string{"validate"})
				Expect(err).ToNot(HaveOccurred())
				Expect(validateCmd).ToNot(BeNil())
				Expect(validateCmd.Name()).To(Equal("validate"))
			})
		})

		Context("when loading configuration from environment", func() {
			BeforeEach(func() {
				os.Setenv("CMDGUARD_STRICT_MODE", "true")
				root = cmdguard.New("testapp", "Test application")
			})

			AfterEach(func() {
				os.Unsetenv("CMDGUARD_STRICT_MODE")
			})

			It("should enable strict mode from environment", func() {
				Expect(root.IsStrictMode()).To(BeTrue())
			})
		})
	})

	Describe("Adding commands", func() {
		BeforeEach(func() {
			root = cmdguard.New("testapp", "Test application")
		})

		Context("when command is valid", func() {
			It("should accept commands with Run handler", func() {
				cmd := &cobra.Command{
					Use:  "valid",
					Short: "A valid command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should accept commands with RunE handler", func() {
				cmd := &cobra.Command{
					Use:   "valid",
					Short: "A valid command",
					RunE:  func(cmd *cobra.Command, args []string) error { return nil },
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should accept parent commands with subcommands", func() {
				parent := &cobra.Command{
					Use:   "parent",
					Short: "Parent command",
				}
				child := &cobra.Command{
					Use:  "child",
					Short: "Child command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				parent.AddCommand(child)
				Expect(func() { root.AddCommand(parent) }).NotTo(Panic())
			})
		})

		Context("when command is invalid", func() {
			It("should panic on missing handler", func() {
				cmd := &cobra.Command{
					Use:   "invalid",
					Short: "Invalid command",
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})

			It("should panic on missing name", func() {
				cmd := &cobra.Command{
					Short: "No name command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})

			It("should panic with message containing cmdguard prefix", func() {
				cmd := &cobra.Command{
					Use:   "invalid",
					Short: "Invalid command",
				}
				panicFn := func() { root.AddCommand(cmd) }
				Expect(panicFn).To(PanicWith(ContainSubstring("cmdguard:")))
			})
		})

		Context("in strict mode", func() {
			BeforeEach(func() {
				os.Setenv("CMDGUARD_STRICT_MODE", "true")
				root = cmdguard.New("testapp", "Test application")
			})

			AfterEach(func() {
				os.Unsetenv("CMDGUARD_STRICT_MODE")
			})

			It("should accept RunE handlers", func() {
				cmd := &cobra.Command{
					Use:   "valid",
					Short: "Valid in strict mode",
					RunE:  func(cmd *cobra.Command, args []string) error { return nil },
				}
				Expect(func() { root.AddCommand(cmd) }).NotTo(Panic())
			})

			It("should reject Run handlers", func() {
				cmd := &cobra.Command{
					Use:   "invalid",
					Short: "Invalid in strict mode",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).To(Panic())
			})

			It("should mention strict mode in panic message", func() {
				cmd := &cobra.Command{
					Use:   "invalid",
					Short: "Invalid in strict mode",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				panicFn := func() { root.AddCommand(cmd) }
				Expect(panicFn).To(PanicWith(ContainSubstring("strict mode")))
			})
		})
	})

	Describe("Adding subcommands", func() {
		var parent *cobra.Command

		BeforeEach(func() {
			root = cmdguard.New("testapp", "Test application")
			parent = &cobra.Command{
				Use:   "parent",
				Short: "Parent command",
				Run:   func(cmd *cobra.Command, args []string) {},
			}
			root.AddCommand(parent)
		})

		Context("when subcommand is valid", func() {
			It("should accept subcommands with Run handler", func() {
				child := &cobra.Command{
					Use:  "child",
					Short: "Child command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddSubcommand(parent, child) }).NotTo(Panic())
			})
		})

		Context("when subcommand is invalid", func() {
			It("should panic on missing handler", func() {
				child := &cobra.Command{
					Use:   "invalid",
					Short: "Invalid child",
				}
				Expect(func() { root.AddSubcommand(parent, child) }).To(Panic())
			})
		})
	})

	Describe("Command execution", func() {
		BeforeEach(func() {
			root = cmdguard.New("testapp", "Test application")
		})

		Context("when executing version command", func() {
			It("should run without error", func() {
				root.Command().SetArgs([]string{"version"})
				err := root.Execute(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when executing validate command", func() {
			It("should run without error", func() {
				root.Command().SetArgs([]string{"validate"})
				err := root.Execute(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Describe("Version function", func() {
		It("should return a version string", func() {
			version := cmdguard.Version()
			Expect(version).ToNot(BeEmpty())
		})

		It("should return dev by default", func() {
			version := cmdguard.Version()
			Expect(version).To(Equal("dev"))
		})
	})

	Describe("Accessing configuration", func() {
		BeforeEach(func() {
			root = cmdguard.New("testapp", "Test application")
		})

		It("should provide access to config", func() {
			cfg := root.Config()
			Expect(cfg).ToNot(BeNil())
		})

		It("should provide access to underlying command", func() {
			cmd := root.Command()
			Expect(cmd).ToNot(BeNil())
			Expect(cmd.Name()).To(Equal("testapp"))
		})
	})

	Describe("ExecuteAndExit", func() {
		BeforeEach(func() {
			root = cmdguard.New("testapp", "Test application")
		})

		Context("when command executes successfully", func() {
			It("should not call os.Exit", func() {
				root.Command().SetArgs([]string{"version"})
				Expect(func() { root.ExecuteAndExit(context.Background()) }).NotTo(Panic())
			})
		})
	})

	Describe("Adding commands after execution", func() {
		Context("when trying to add command after Execute", func() {
			It("should panic when adding command after execution", func() {
				root = cmdguard.New("testapp", "Test application")
				root.Command().SetArgs([]string{"version"})
				_ = root.Execute(context.Background())

				cmd := &cobra.Command{
					Use:  "late",
					Short: "Late command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddCommand(cmd) }).To(PanicWith(ContainSubstring("cannot add commands after execution")))
			})

			It("should panic when adding subcommand after execution", func() {
				root = cmdguard.New("testapp", "Test application")
				parent := &cobra.Command{
					Use:  "parent",
					Short: "Parent command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				root.AddCommand(parent)
				root.Command().SetArgs([]string{"version"})
				_ = root.Execute(context.Background())

				child := &cobra.Command{
					Use:  "child",
					Short: "Child command",
					Run:   func(cmd *cobra.Command, args []string) {},
				}
				Expect(func() { root.AddSubcommand(parent, child) }).To(PanicWith(ContainSubstring("cannot add commands after execution")))
			})
		})
	})

	Describe("Command execution", func() {
		Context("when running commands", func() {
			It("should execute version command successfully", func() {
				root = cmdguard.New("testapp", "Test application")
				root.Command().SetArgs([]string{"version"})

				err := root.Execute(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("with CMDGUARD_LOG_FORMAT env var", func() {
			BeforeEach(func() {
				os.Setenv("CMDGUARD_LOG_FORMAT", "json")
			})

			AfterEach(func() {
				os.Unsetenv("CMDGUARD_LOG_FORMAT")
			})

			It("should use json format", func() {
				root = cmdguard.New("testapp", "Test application")
				Expect(root.Config().LogFormat).To(Equal("json"))
			})
		})
	})
})

