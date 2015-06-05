package extension_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/extension_service"
)

var _ = Describe("OpenStackComputeExtensions", func() {
	var (
		computeExtensions OpenStackComputeExtensions
	)

	BeforeEach(func() {
		computeExtensions = make(map[string]struct{})
	})

	Describe("HasBlockDeviceMappingSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasBlockDeviceMappingSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-block-device-mapping-v2-boot"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasBlockDeviceMappingSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasConfigDriveSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasConfigDriveSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-config-drive"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasConfigDriveSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasFloatingIpsSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasFloatingIpsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-floating-ips"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasFloatingIpsSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasKeyPairsSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasKeyPairsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-keypairs"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasKeyPairsSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasSchedulerHintsSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasSchedulerHintsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["OS-SCH-HNT"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasSchedulerHintsSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasSecurityGroupsSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasSecurityGroupsSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-security-groups"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasSecurityGroupsSupport()).To(BeTrue())
			})
		})
	})

	Describe("HasTenantNetworksSupport", func() {
		It("returns true", func() {
			Expect(computeExtensions.HasTenantNetworksSupport()).To(BeFalse())
		})

		Context("when extension is enabled", func() {
			BeforeEach(func() {
				computeExtensions["os-tenant-networks"] = struct{}{}
			})

			It("returns true", func() {
				Expect(computeExtensions.HasTenantNetworksSupport()).To(BeTrue())
			})
		})
	})
})
