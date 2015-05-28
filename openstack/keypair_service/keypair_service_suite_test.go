package keypair_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestKeyPairService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "KeyPair Service Suite")
}
