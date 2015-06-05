package extension_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestExtensionService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Extension Service Suite")
}
