package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/internal/config"
)

var _ = Describe("Configuration - User Expectations", func() {
	BeforeEach(func() {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	AfterEach(func() {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("As a developer deploying to production", func() {
		Context("when I don't set any environment variables", func() {
			It("should work out of the box with sensible defaults", func() {
				cfg := config.Load()
				Expect(cfg.Validate()).To(Succeed(), "Default config should be valid")

				Expect(cfg.LogLevel).To(Equal("info"))
				Expect(cfg.LogFormat).To(Equal("text"))
				Expect(cfg.StrictMode).To(BeFalse())
			})
		})

		Context("when I want structured logs for my log aggregator", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_LOG_FORMAT", "json")
			})

			It("should switch to JSON format", func() {
				cfg := config.Load()
				Expect(cfg.LogFormat).To(Equal("json"))
				Expect(cfg.Validate()).To(Succeed())
			})
		})

		Context("when I want verbose logs for debugging", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_LOG_LEVEL", "debug")
			})

			It("should enable debug logging", func() {
				cfg := config.Load()
				Expect(cfg.LogLevel).To(Equal("debug"))
				Expect(cfg.Validate()).To(Succeed())
			})
		})
	})

	Describe("As a platform operator running in strict environments", func() {
		Context("when I need guaranteed error handling behavior", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
			})

			It("should enable strict mode", func() {
				cfg := config.Load()
				Expect(cfg.StrictMode).To(BeTrue())
			})
		})

		Context("when strict mode is set to non-true values", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_STRICT_MODE", "yes")
			})

			It("should NOT enable strict mode (only 'true' is accepted)", func() {
				cfg := config.Load()
				Expect(cfg.StrictMode).To(BeFalse())
			})
		})
	})

	Describe("As a user configuring via CI/CD environment variables", func() {
		Context("when I combine multiple settings", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_LOG_LEVEL", "warn")
				_ = os.Setenv("CMDGUARD_LOG_FORMAT", "json")
				_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
			})

			It("should apply all settings together", func() {
				cfg := config.Load()
				Expect(cfg.LogLevel).To(Equal("warn"))
				Expect(cfg.LogFormat).To(Equal("json"))
				Expect(cfg.StrictMode).To(BeTrue())
				Expect(cfg.Validate()).To(Succeed())
			})
		})
	})

	Describe("As a user making configuration mistakes", func() {
		Context("when I typo the log level", func() {
			It("should tell me exactly what went wrong and what's valid", func() {
				cfg := &config.Config{LogLevel: "debbug"}
				err := cfg.Validate()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid log level"))
				Expect(err.Error()).To(ContainSubstring("debbug"))
				Expect(err.Error()).To(ContainSubstring("debug"))
				Expect(err.Error()).To(ContainSubstring("info"))
				Expect(err.Error()).To(ContainSubstring("warn"))
				Expect(err.Error()).To(ContainSubstring("error"))
			})
		})

		Context("when I typo the log format", func() {
			It("should tell me exactly what went wrong and what's valid", func() {
				cfg := &config.Config{LogFormat: "jso"}
				err := cfg.Validate()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid log format"))
				Expect(err.Error()).To(ContainSubstring("jso"))
				Expect(err.Error()).To(ContainSubstring("text"))
				Expect(err.Error()).To(ContainSubstring("json"))
			})
		})
	})

	Describe("As a user with case-sensitive input", func() {
		Context("when I use uppercase log level", func() {
			It("should reject it (case-sensitive)", func() {
				cfg := &config.Config{LogLevel: "DEBUG"}
				Expect(cfg.Validate()).To(HaveOccurred())
			})
		})

		Context("when I use uppercase log format", func() {
			It("should reject it (case-sensitive)", func() {
				cfg := &config.Config{LogFormat: "JSON"}
				Expect(cfg.Validate()).To(HaveOccurred())
			})
		})
	})

	Describe("As a user who leaves settings empty", func() {
		Context("when log level is empty", func() {
			It("should be valid (uses default)", func() {
				cfg := &config.Config{LogLevel: ""}
				Expect(cfg.Validate()).To(Succeed())
			})
		})

		Context("when log format is empty", func() {
			It("should be valid (uses default)", func() {
				cfg := &config.Config{LogFormat: ""}
				Expect(cfg.Validate()).To(Succeed())
			})
		})
	})
})

var _ = Describe("Config file path - User Expectations", func() {
	Describe("As a user specifying a config file", func() {
		Context("when I don't specify a config file", func() {
			It("should return empty string (no config file)", func() {
				Expect(config.GetConfigFilePath("")).To(Equal(""))
			})
		})

		Context("when I specify a relative path", func() {
			It("should convert to absolute path for consistency", func() {
				result := config.GetConfigFilePath("config.yaml")
				Expect(result).ToNot(Equal("config.yaml"))
				Expect(result).To(ContainSubstring("config.yaml"))
			})
		})
	})
})

var _ = Describe("Environment variable precedence - User Expectations", func() {
	BeforeEach(func() {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	AfterEach(func() {
		_ = os.Unsetenv("CMDGUARD_LOG_LEVEL")
		_ = os.Unsetenv("CMDGUARD_LOG_FORMAT")
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("As a DevOps engineer understanding configuration priority", func() {
		Context("when environment variable overrides default", func() {
			BeforeEach(func() {
				_ = os.Setenv("CMDGUARD_LOG_LEVEL", "error")
			})

			It("should use environment variable over default", func() {
				cfg := config.Load()
				Expect(cfg.LogLevel).To(Equal("error"), "Env var should override default 'info'")
			})
		})

		Context("when I set all common log levels", func() {
			It("should accept all standard levels", func() {
				for _, level := range []string{"debug", "info", "warn", "error"} {
					_ = os.Setenv("CMDGUARD_LOG_LEVEL", level)
					cfg := config.Load()
					Expect(cfg.LogLevel).To(Equal(level))
					Expect(cfg.Validate()).To(Succeed(), "Level %s should be valid", level)
				}
			})
		})
	})
})

var _ = Describe("Strict mode behavior - User Expectations", func() {
	BeforeEach(func() {
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	AfterEach(func() {
		_ = os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("As a user understanding strict mode activation", func() {
		Context("when I set strict mode to various truthy-looking values", func() {
			truthyAttempts := []string{"yes", "1", "TRUE", "True", "on", "enabled"}

			for _, value := range truthyAttempts {
				value := value // capture range variable
				It("should NOT enable strict mode for '"+value+"' (only 'true' works)", func() {
					_ = os.Setenv("CMDGUARD_STRICT_MODE", value)
					cfg := config.Load()
					Expect(cfg.StrictMode).To(BeFalse(), "Only 'true' should enable strict mode, not '%s'", value)
				})
			}
		})

		Context("when I set strict mode to 'true' (lowercase)", func() {
			It("should enable strict mode", func() {
				_ = os.Setenv("CMDGUARD_STRICT_MODE", "true")
				cfg := config.Load()
				Expect(cfg.StrictMode).To(BeTrue())
			})
		})
	})
})
