package client_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOpenStackClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "OpenStack Client Suite")
}
