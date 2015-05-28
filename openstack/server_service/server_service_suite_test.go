package server_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestServerService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Service Suite")
}
