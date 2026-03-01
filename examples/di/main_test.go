package main_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	v2 "github.com/larsartmann/cmdguard/pkg/cmdguard/v2"
)

var _ = Describe("DI Example", func() {
	var (
		ctx  context.Context
		root *v2.GuardedCommand[Config, v2.NoFlags]
	)

	BeforeEach(func() {
		ctx = context.Background()
		var err error
		root, err = v2.New[Config, v2.NoFlags]("di-app", "DI Example App", Config{})
		Expect(err).ToNot(HaveOccurred())

		// Register services
		Expect(v2.Provide(root.ScopeStruct(), NewDatabaseService)).To(Succeed())
		Expect(v2.Provide(root.ScopeStruct(), NewAPIService)).To(Succeed())
	})

	Describe("Service Registration", func() {
		It("registers database service", func() {
			db, err := v2.Invoke[*DatabaseService](root.ScopeStruct())
			Expect(err).ToNot(HaveOccurred())
			Expect(db).ToNot(BeNil())
			Expect(db.IsConnected()).To(BeTrue())
		})

		It("registers API service", func() {
			api, err := v2.Invoke[*APIService](root.ScopeStruct())
			Expect(err).ToNot(HaveOccurred())
			Expect(api).ToNot(BeNil())
		})
	})

	Describe("Health Checks", func() {
		It("passes health check when connected", func() {
			err := root.HealthCheckWithContext(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		It("fails health check after shutdown", func() {
			// First verify it's healthy
			Expect(root.HealthCheckWithContext(ctx)).To(Succeed())

			// Shutdown
			shutdownCtx, cancel := context.WithTimeout(ctx, 5)
			defer cancel()
			Expect(root.Shutdown(shutdownCtx)).To(Succeed())

			// Health check should fail
			err := root.HealthCheckWithContext(ctx)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("DI Helpers", func() {
		It("uses MustInvoke for service retrieval", func() {
			db := v2.MustInvoke[*DatabaseService](root.ScopeStruct())
			Expect(db).ToNot(BeNil())
			Expect(db.IsConnected()).To(BeTrue())
		})
	})
})

func TestDIExample(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DI Example Suite")
}
