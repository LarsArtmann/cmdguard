package main_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

var _ = Describe("Advanced Flags Example", func() {
	var (
		ctx  context.Context
		root *v2.GuardedCommand[GlobalConfig, v2.NoFlags]
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		root, err = v2.New[GlobalConfig, v2.NoFlags]("advflags", "Advanced Flags Example", GlobalConfig{})
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("Server Command", func() {
		BeforeEach(func() {
			err := root.AddCommand(v2.Command[GlobalConfig, ServerFlags]{
				Use:   "server",
				Short: "Start the server",
				Flags: ServerFlags{},
				RunE: func(ctx context.Context, cfg *GlobalConfig, flags ServerFlags) error {
					return nil
				},
			})
			Expect(err).ToNot(HaveOccurred())
		})

		It("accepts valid server flags", func() {
			// Command should be valid with default flags
			cmd := v2.Command[GlobalConfig, ServerFlags]{
				Use:   "test",
				Short: "Test",
				Flags: ServerFlags{Port: 8080, Host: "localhost"},
				RunE:  func(ctx context.Context, cfg *GlobalConfig, flags ServerFlags) error { return nil },
			}
			err := cmd.Validate()
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Describe("Enum Validation", func() {
		It("validates correct environment", func() {
			flags := EnumFlags{Environment: "production", Region: "us-west-2"}
			err := flags.Validate()
			Expect(err).ToNot(HaveOccurred())
		})

		It("rejects invalid environment", func() {
			flags := EnumFlags{Environment: "invalid", Region: "us-west-2"}
			err := flags.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid environment"))
		})
	})

	Describe("Duration Flags", func() {
		It("parses duration correctly", func() {
			duration, err := v2.ParseDuration("1h30m")
			Expect(err).ToNot(HaveOccurred())
			Expect(duration.Hours()).To(Equal(1.5))
		})

		It("rejects invalid duration", func() {
			_, err := v2.ParseDuration("not-a-duration")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Format Suggestion", func() {
		It("suggests yaml for yam", func() {
			suggestion := suggestFormat("yam")
			Expect(suggestion).To(Equal("yaml"))
		})

		It("suggests json for jsn", func() {
			suggestion := suggestFormat("jsn")
			Expect(suggestion).To(Equal("json"))
		})
	})
})

func TestAdvancedFlagsExample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Advanced Flags Example Suite")
}
