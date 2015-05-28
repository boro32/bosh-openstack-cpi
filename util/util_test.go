package util_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/util"
)

var _ = Describe("Util", func() {
	Describe("ConvertMib2Gib", func() {
		It("converts Mib to Gib", func() {
			Expect(ConvertMib2Gib(32768)).To(Equal(32))
		})
	})
})
