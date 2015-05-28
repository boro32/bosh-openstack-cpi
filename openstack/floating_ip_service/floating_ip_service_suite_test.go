package floatingip_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFloatingIPService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Floating IP Service Suite")
}
