package extension_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/extension_service"
)

var _ = Describe("OpenStackNetworkExtensions", func() {
	var (
		networkExtensions OpenStackNetworkExtensions
	)

	BeforeEach(func() {
		networkExtensions = make(map[string]struct{})
	})

	Describe("HasFloatingIpsSupport", func() {
		It("returns true", func() {
			Expect(networkExtensions.HasFloatingIpsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				networkExtensions["floatingips"] = struct{}{}
			})

			It("returns true", func() {
				Expect(networkExtensions.HasFloatingIpsSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasSecurityGroupsSupport", func() {
		It("returns true", func() {
			Expect(networkExtensions.HasSecurityGroupsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				networkExtensions["security-group"] = struct{}{}
			})

			It("returns true", func() {
				Expect(networkExtensions.HasSecurityGroupsSupport()).To(BeTrue())
			})
		})
	})
})
