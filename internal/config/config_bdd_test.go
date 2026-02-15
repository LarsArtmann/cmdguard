package config_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/internal/config"
)

var _ = Describe("Config", func() {
	BeforeEach(func() {
		os.Unsetenv("CMDGUARD_LOG_LEVEL")
		os.Unsetenv("CMDGUARD_LOG_FORMAT")
		os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	AfterEach(func() {
		os.Unsetenv("CMDGUARD_LOG_LEVEL")
		os.Unsetenv("CMDGUARD_LOG_FORMAT")
		os.Unsetenv("CMDGUARD_STRICT_MODE")
	})

	Describe("Loading configuration", func() {
		Context("with no environment variables", func() {
			It("should return default values", func() {
				cfg := config.Load()
				Expect(cfg.LogLevel).To(Equal("info"))
				Expect(cfg.LogFormat).To(Equal("text"))
				Expect(cfg.StrictMode).To(BeFalse())
			})
		})

		Context("with environment variables", func() {
			BeforeEach(func() {
				os.Setenv("CMDGUARD_LOG_LEVEL", "debug")
				os.Setenv("CMDGUARD_LOG_FORMAT", "json")
				os.Setenv("CMDGUARD_STRICT_MODE", "true")
			})

			It("should load log level from environment", func() {
				cfg := config.Load()
				Expect(cfg.LogLevel).To(Equal("debug"))
			})

			It("should load log format from environment", func() {
				cfg := config.Load()
				Expect(cfg.LogFormat).To(Equal("json"))
			})

			It("should load strict mode from environment", func() {
				cfg := config.Load()
				Expect(cfg.StrictMode).To(BeTrue())
			})
		})

		Context("with invalid strict mode value", func() {
			BeforeEach(func() {
				os.Setenv("CMDGUARD_STRICT_MODE", "yes")
			})

			It("should not enable strict mode for non-true values", func() {
				cfg := config.Load()
				Expect(cfg.StrictMode).To(BeFalse())
			})
		})
	})

	Describe("Validating configuration", func() {
		Context("with valid log level", func() {
			It("should accept debug level", func() {
				cfg := &config.Config{LogLevel: "debug"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})

			It("should accept info level", func() {
				cfg := &config.Config{LogLevel: "info"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})

			It("should accept warn level", func() {
				cfg := &config.Config{LogLevel: "warn"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})

			It("should accept error level", func() {
				cfg := &config.Config{LogLevel: "error"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})
		})

		Context("with invalid log level", func() {
			It("should reject invalid level", func() {
				cfg := &config.Config{LogLevel: "invalid"}
				err := cfg.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid log level"))
			})
		})

		Context("with valid log format", func() {
			It("should accept text format", func() {
				cfg := &config.Config{LogFormat: "text"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})

			It("should accept json format", func() {
				cfg := &config.Config{LogFormat: "json"}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})
		})

		Context("with invalid log format", func() {
			It("should reject invalid format", func() {
				cfg := &config.Config{LogFormat: "xml"}
				err := cfg.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid log format"))
			})
		})

		Context("with empty values", func() {
			It("should accept empty log level", func() {
				cfg := &config.Config{LogLevel: ""}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})

			It("should accept empty log format", func() {
				cfg := &config.Config{LogFormat: ""}
				Expect(cfg.Validate()).ToNot(HaveOccurred())
			})
		})
	})

	Describe("GetConfigFilePath", func() {
		Context("with empty path", func() {
			It("should return empty string", func() {
				Expect(config.GetConfigFilePath("")).To(Equal(""))
			})
		})

		Context("with relative path", func() {
			It("should return absolute path", func() {
				result := config.GetConfigFilePath("config.yaml")
				Expect(result).ToNot(Equal("config.yaml"))
				Expect(result).To(ContainSubstring("config.yaml"))
			})
		})
	})
})
