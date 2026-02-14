package cmdguard_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCmdguard(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmdguard Suite")
}
