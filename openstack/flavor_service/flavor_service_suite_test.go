package flavor_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFlavorService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Flavor Service Suite")
}
