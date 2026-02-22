package logging_test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"regexp"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/internal/logging"
)

// validationTestCase represents a single validation test case
type validationTestCase struct {
	input    string
	expected bool
}

// testValidation is a parameterized helper for testing validation functions
func testValidation(name string, validator func(string) bool, cases []validationTestCase) {
	It("should NOT accept uppercase "+name+" (case-sensitive)", func() {
		for _, tc := range cases {
			Expect(validator(tc.input)).To(Equal(tc.expected))
		}
	})
}

// captureStderr captures stderr output during test execution
func captureStderr(fn func()) string {
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	_ = w.Close()
	os.Stderr = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

var _ = Describe("Logging - User Expectations", func() {
	Describe("As a DevOps engineer setting up log aggregation", func() {
		Context("when I configure JSON format for structured logging", func() {
			It("should produce valid JSON that my log aggregator can parse", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info("user logged in", "user_id", "12345", "ip", "192.168.1.1")
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed(), "Output should be valid JSON")

				Expect(parsed["msg"]).To(Equal("user logged in"))
				Expect(parsed["level"]).To(Equal("INFO"))
				Expect(parsed["user_id"]).To(Equal("12345"))
				Expect(parsed["ip"]).To(Equal("192.168.1.1"))
			})

			It("should include timestamp for log correlation", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info("test message")
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed())
				Expect(parsed).To(HaveKey("time"), "JSON logs should include timestamp")
			})
		})
	})

	Describe("As a developer debugging an issue", func() {
		Context("when I set log level to debug", func() {
			It("should show all log levels including debug details", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "debug")
					logger.Debug("debugging connection", "host", "localhost")
					logger.Info("connection established")
					logger.Warn("slow response", "latency_ms", 500)
					logger.Error("timeout occurred")
				})

				Expect(output).To(ContainSubstring(`"level":"DEBUG"`))
				Expect(output).To(ContainSubstring(`"level":"INFO"`))
				Expect(output).To(ContainSubstring(`"level":"WARN"`))
				Expect(output).To(ContainSubstring(`"level":"ERROR"`))
			})
		})

		Context("when I set log level to error to reduce noise", func() {
			It("should only show error messages, hiding debug/info/warn", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "error")
					logger.Debug("this should be hidden")
					logger.Info("this should be hidden too")
					logger.Warn("this should also be hidden")
					logger.Error("only this should appear")
				})

				Expect(output).ToNot(ContainSubstring("this should be hidden"))
				Expect(output).To(ContainSubstring("only this should appear"))
				Expect(output).To(ContainSubstring(`"level":"ERROR"`))
			})
		})
	})

	Describe("As a terminal user reading logs", func() {
		Context("when I use text format for human-readable output", func() {
			It("should produce easy-to-read output with key=value pairs", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("text", "info")
					logger.Info("request completed", "status", 200, "duration_ms", 42)
				})

				Expect(output).To(ContainSubstring("request completed"))
				Expect(output).To(ContainSubstring("status=200"))
				Expect(output).To(ContainSubstring("duration_ms=42"))
			})

			It("should not be JSON (should be plain text)", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("text", "info")
					logger.Info("test message")
				})

				Expect(strings.HasPrefix(strings.TrimSpace(output), "{")).To(BeFalse(), "Text format should not start with JSON brace")
			})
		})
	})

	Describe("As a platform operator handling misconfiguration", func() {
		Context("when an invalid log format is provided", func() {
			It("should gracefully fall back to text format instead of crashing", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("xml", "info")
					logger.Info("app started")
				})

				Expect(output).To(ContainSubstring("app started"))
				Expect(strings.HasPrefix(strings.TrimSpace(output), "{")).To(BeFalse(), "Should use text format as fallback")
			})
		})

		Context("when an invalid log level is provided", func() {
			It("should gracefully fall back to info level instead of crashing", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("text", "trace")
					logger.Info("app started")
					logger.Debug("this should be hidden since level defaults to info")
				})

				Expect(output).To(ContainSubstring("app started"))
				Expect(output).ToNot(ContainSubstring("this should be hidden"))
			})
		})

		Context("when empty format or level is provided", func() {
			It("should use sensible defaults", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("", "")
					logger.Info("using defaults")
				})

				Expect(output).To(ContainSubstring("using defaults"))
			})
		})
	})

	Describe("As a security auditor reviewing logs", func() {
		Context("when logging sensitive operation data", func() {
			It("should log structured data without losing context", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info("user action",
						"user_id", "user-123",
						"action", "delete_resource",
						"resource_id", "res-456",
						"ip_address", "10.0.0.1",
					)
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed())

				Expect(parsed["user_id"]).To(Equal("user-123"))
				Expect(parsed["action"]).To(Equal("delete_resource"))
				Expect(parsed["resource_id"]).To(Equal("res-456"))
				Expect(parsed["ip_address"]).To(Equal("10.0.0.1"))
			})
		})
	})

	Describe("As an SRE monitoring production", func() {
		Context("when errors occur", func() {
			It("should capture error details for incident response", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Error("database connection failed",
						"error", "connection refused",
						"host", "db.production.internal",
						"port", 5432,
						"retry_count", 3,
					)
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed())

				Expect(parsed["level"]).To(Equal("ERROR"))
				Expect(parsed["msg"]).To(Equal("database connection failed"))
				Expect(parsed["host"]).To(Equal("db.production.internal"))
				Expect(parsed["retry_count"]).To(Equal(float64(3)))
			})
		})
	})
})

var _ = Describe("Format validation - User Expectations", func() {
	Describe("As a CLI developer validating user input", func() {
		Context("when checking if a format string is acceptable", func() {
			It("should accept 'text' for human-readable logs", func() {
				Expect(logging.ValidFormat("text")).To(BeTrue())
			})

			It("should accept 'json' for machine-parseable logs", func() {
				Expect(logging.ValidFormat("json")).To(BeTrue())
			})

			It("should reject formats that would cause parsing issues downstream", func() {
				Expect(logging.ValidFormat("xml")).To(BeFalse())
				Expect(logging.ValidFormat("yaml")).To(BeFalse())
				Expect(logging.ValidFormat("csv")).To(BeFalse())
			})

			It("should reject empty format to prevent ambiguity", func() {
				Expect(logging.ValidFormat("")).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Level validation - User Expectations", func() {
	Describe("As a CLI developer validating log level configuration", func() {
		Context("when checking if a level string is acceptable", func() {
			It("should accept standard log levels", func() {
				Expect(logging.ValidLevel("debug")).To(BeTrue())
				Expect(logging.ValidLevel("info")).To(BeTrue())
				Expect(logging.ValidLevel("warn")).To(BeTrue())
				Expect(logging.ValidLevel("error")).To(BeTrue())
			})

			It("should reject non-standard levels to prevent confusion", func() {
				Expect(logging.ValidLevel("trace")).To(BeFalse())
				Expect(logging.ValidLevel("fatal")).To(BeFalse())
				Expect(logging.ValidLevel("critical")).To(BeFalse())
				Expect(logging.ValidLevel("verbose")).To(BeFalse())
			})

			It("should reject empty level to force explicit configuration", func() {
				Expect(logging.ValidLevel("")).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Level precedence - User Expectations", func() {
	Describe("As a user understanding log level filtering", func() {
		Context("when understanding which levels include which", func() {
			levels := []struct {
				name   string
				allows []string
				blocks []string
			}{
				{"debug", []string{"debug", "info", "warn", "error"}, []string{}},
				{"info", []string{"info", "warn", "error"}, []string{"debug"}},
				{"warn", []string{"warn", "error"}, []string{"debug", "info"}},
				{"error", []string{"error"}, []string{"debug", "info", "warn"}},
			}

			for _, level := range levels {
				level := level // capture range variable
				Context("when log level is "+level.name, func() {
					It("should allow expected levels and block others", func() {
						output := captureStderr(func() {
							logger := logging.NewLogger("json", level.name)
							logger.Debug("debug message")
							logger.Info("info message")
							logger.Warn("warn message")
							logger.Error("error message")
						})

						for _, allowed := range level.allows {
							Expect(output).To(ContainSubstring(allowed+" message"),
								"Level %s should allow %s", level.name, allowed)
						}
						for _, blocked := range level.blocks {
							Expect(output).ToNot(ContainSubstring(blocked+" message"),
								"Level %s should block %s", level.name, blocked)
						}
					})
				})
			}
		})
	})
})

var _ = Describe("Edge cases - User Expectations", func() {
	Describe("As a user handling unusual logging scenarios", func() {
		Context("when logging messages with special characters", func() {
			It("should handle quotes and newlines safely in JSON", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info(`message with "quotes" and \backslashes`)
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed(), "Should produce valid JSON despite special chars")
				Expect(parsed["msg"]).To(ContainSubstring("quotes"))
			})
		})

		Context("when logging very long messages", func() {
			It("should not truncate or crash", func() {
				longMessage := strings.Repeat("a", 10000)
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info(longMessage)
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed())
				Expect(parsed["msg"]).To(Equal(longMessage))
			})
		})

		Context("when logging Unicode characters", func() {
			It("should preserve international characters", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("json", "info")
					logger.Info("用户登录", "用户", "测试") // Chinese: "User login", "user", "test"
				})

				var parsed map[string]any
				Expect(json.Unmarshal([]byte(output), &parsed)).To(Succeed())
				Expect(parsed["msg"]).To(Equal("用户登录"))
				Expect(parsed["用户"]).To(Equal("测试"))
			})
		})
	})
})

var _ = Describe("Output destination - User Expectations", func() {
	Describe("As a user understanding where logs go", func() {
		Context("when logs are written", func() {
			It("should write to stderr, not stdout (following CLI conventions)", func() {
				stderrOutput := captureStderr(func() {
					logger := logging.NewLogger("text", "info")
					logger.Info("test message")
				})

				Expect(stderrOutput).To(ContainSubstring("test message"))
			})
		})
	})
})

var _ = Describe("Case sensitivity - User Expectations", func() {
	Describe("As a user providing configuration values", func() {
		Context("when I provide uppercase or mixed-case values", func() {
			testValidation("level", logging.ValidLevel, []validationTestCase{
				{"DEBUG", false},
				{"INFO", false},
				{"Error", false},
			})

			testValidation("format", logging.ValidFormat, []validationTestCase{
				{"JSON", false},
				{"TEXT", false},
				{"Json", false},
			})

			It("should fall back gracefully when given uppercase", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("JSON", "INFO")
					logger.Info("test")
				})

				Expect(output).To(ContainSubstring("test"), "Should still work with fallback to defaults")
			})
		})
	})
})

var _ = Describe("Text format parsing - User Expectations", func() {
	Describe("As a user reading text logs", func() {
		Context("when parsing text format with grep/awk", func() {
			It("should have predictable format for text parsing", func() {
				output := captureStderr(func() {
					logger := logging.NewLogger("text", "info")
					logger.Info("test message", "key1", "value1")
				})

				levelPattern := regexp.MustCompile(`level=INFO`)
				msgPattern := regexp.MustCompile(`msg="?test message"?`)
				keyPattern := regexp.MustCompile(`key1=value1`)

				Expect(levelPattern.MatchString(output)).To(BeTrue(), "Should contain level=INFO")
				Expect(msgPattern.MatchString(output)).To(BeTrue(), "Should contain msg=test message")
				Expect(keyPattern.MatchString(output)).To(BeTrue(), "Should contain key=value pairs")
			})
		})
	})
})
