package securitygroup_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSecurityGroupService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SecurityGroup Service Suite")
}
