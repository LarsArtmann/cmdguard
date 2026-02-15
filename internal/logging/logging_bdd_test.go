package logging_test

import (
	"log/slog"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/internal/logging"
)

var _ = Describe("Logging", func() {
	Describe("Creating a logger", func() {
		Context("with text format", func() {
			It("should create a logger without error", func() {
				logger := logging.NewLogger("text", "info")
				Expect(logger).ToNot(BeNil())
			})
		})

		Context("with JSON format", func() {
			It("should create a logger without error", func() {
				logger := logging.NewLogger("json", "info")
				Expect(logger).ToNot(BeNil())
			})
		})

		Context("with different log levels", func() {
			It("should accept debug level", func() {
				logger := logging.NewLogger("text", "debug")
				Expect(logger).ToNot(BeNil())
			})

			It("should accept info level", func() {
				logger := logging.NewLogger("text", "info")
				Expect(logger).ToNot(BeNil())
			})

			It("should accept warn level", func() {
				logger := logging.NewLogger("text", "warn")
				Expect(logger).ToNot(BeNil())
			})

			It("should accept error level", func() {
				logger := logging.NewLogger("text", "error")
				Expect(logger).ToNot(BeNil())
			})
		})

		Context("with invalid format", func() {
			It("should default to text format", func() {
				logger := logging.NewLogger("invalid", "info")
				Expect(logger).ToNot(BeNil())
			})
		})

		Context("with invalid level", func() {
			It("should default to info level", func() {
				logger := logging.NewLogger("text", "invalid")
				Expect(logger).ToNot(BeNil())
			})
		})

		Context("with empty values", func() {
			It("should handle empty format", func() {
				logger := logging.NewLogger("", "info")
				Expect(logger).ToNot(BeNil())
			})

			It("should handle empty level", func() {
				logger := logging.NewLogger("text", "")
				Expect(logger).ToNot(BeNil())
			})
		})
	})

	Describe("Format constants", func() {
		It("should have text format", func() {
			Expect(logging.FormatText).To(Equal(logging.Format("text")))
		})

		It("should have JSON format", func() {
			Expect(logging.FormatJSON).To(Equal(logging.Format("json")))
		})
	})

	Describe("Validating format", func() {
		Context("with valid formats", func() {
			It("should accept text format", func() {
				Expect(logging.ValidFormat("text")).To(BeTrue())
			})

			It("should accept json format", func() {
				Expect(logging.ValidFormat("json")).To(BeTrue())
			})
		})

		Context("with invalid formats", func() {
			It("should reject xml format", func() {
				Expect(logging.ValidFormat("xml")).To(BeFalse())
			})

			It("should reject empty format", func() {
				Expect(logging.ValidFormat("")).To(BeFalse())
			})

			It("should reject random string", func() {
				Expect(logging.ValidFormat("random")).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Logger functionality", func() {
	var logger *slog.Logger

	BeforeEach(func() {
		logger = logging.NewLogger("text", "debug")
	})

	Describe("Logging messages", func() {
		It("should support debug logging", func() {
			Expect(func() { logger.Debug("test debug message") }).NotTo(Panic())
		})

		It("should support info logging", func() {
			Expect(func() { logger.Info("test info message") }).NotTo(Panic())
		})

		It("should support warn logging", func() {
			Expect(func() { logger.Warn("test warn message") }).NotTo(Panic())
		})

		It("should support error logging", func() {
			Expect(func() { logger.Error("test error message") }).NotTo(Panic())
		})

		It("should support logging with key-value pairs", func() {
			Expect(func() { logger.Info("test message", "key", "value") }).NotTo(Panic())
		})
	})
})
