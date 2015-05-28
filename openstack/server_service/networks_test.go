package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

var _ = Describe("Network", func() {
	var (
		dynamicNetwork Network
		vipNetwork     Network
	)

	BeforeEach(func() {
		dynamicNetwork = Network{
			Type:    "dynamic",
			IP:      "fake-dynamic-network-ip",
			Gateway: "fake-dynamic-network-gateway",
			Netmask: "fake-dynamic-network-netmask",
			DNS:     []string{"fake-dynamic-network-dns"},
			Default: []string{"fake-dynamic-network-default"},
			Network: "fake-dynamic-network-network",
		}

		vipNetwork = Network{
			Type:    "vip",
			IP:      "fake-vip-network-ip",
			Gateway: "fake-vip-network-gateway",
			Netmask: "fake-vip-network-netmask",
			DNS:     []string{"fake-vip-network-dns"},
			Default: []string{"fake-vip-network-default"},
			Network: "fake-vip-network-network",
		}
	})

	Describe("IsDynamic", func() {
		It("returns true for a dynamic network", func() {
			Expect(dynamicNetwork.IsDynamic()).To(BeTrue())
		})

		It("returns false for a vip network", func() {
			Expect(vipNetwork.IsDynamic()).To(BeFalse())
		})
	})

	Describe("IsVip", func() {
		It("returns true for a vip network", func() {
			Expect(vipNetwork.IsVip()).To(BeTrue())
		})

		It("returns false for a dynamic network", func() {
			Expect(dynamicNetwork.IsVip()).To(BeFalse())
		})
	})
})

var _ = Describe("Networks", func() {
	var (
		err            error
		dynamicNetwork Network
		vipNetwork     Network
		networks       Networks
	)

	BeforeEach(func() {
		dynamicNetwork = Network{
			Type:    "dynamic",
			IP:      "fake-dynamic-network-ip",
			Gateway: "fake-dynamic-network-gateway",
			Netmask: "fake-dynamic-network-netmask",
			DNS:     []string{"fake-dynamic-network-dns"},
			Default: []string{"fake-dynamic-network-default"},
			Network: "fake-dynamic-network-network",
		}

		vipNetwork = Network{
			Type:    "vip",
			IP:      "fake-vip-network-ip",
			Gateway: "fake-vip-network-gateway",
			Netmask: "fake-vip-network-netmask",
			DNS:     []string{"fake-vip-network-dns"},
			Default: []string{"fake-vip-network-default"},
			Network: "fake-vip-network-network",
		}

		networks = Networks{
			"fake-dynamic-network": dynamicNetwork,
			"fake-vip-network":     vipNetwork,
		}
	})

	Describe("Validate", func() {
		It("should not return an error", func() {
			err = networks.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("At least one 'dynamic' network should be defined"))
			})
		})

		Context("when there is more than one dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-dynamic-network-1": dynamicNetwork,
					"fake-dynamic-network-2": dynamicNetwork,
				}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one dynamic network is allowed"))
			})
		})

		Context("when there is NOT a vip network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("should not return an error", func() {
				err = networks.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when there is more than one vip network", func() {
			BeforeEach(func() {
				networks = Networks{
					"fake-vip-network-1": vipNetwork,
					"fake-vip-network-2": vipNetwork,
				}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Only one VIP network is allowed"))
			})
		})

		Context("when VIP network does not have an IP", func() {
			BeforeEach(func() {
				vipNetwork.IP = ""
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an error", func() {
				err = networks.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("VIP Network must have an IP address"))
			})
		})
	})

	Describe("DynamicNetwork", func() {
		It("should return the dynamic network", func() {
			Expect(networks.DynamicNetwork()).To(Equal(dynamicNetwork))
		})

		Context("when there is NOT a dynamic network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-vip-network": vipNetwork}
			})

			It("should return an emtpy network", func() {
				Expect(networks.DynamicNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("VipNetwork", func() {
		It("should return the vip network", func() {
			Expect(networks.VipNetwork()).To(Equal(vipNetwork))
		})

		Context("when there is NOT a vip network", func() {
			BeforeEach(func() {
				networks = Networks{"fake-dynamic-network": dynamicNetwork}
			})

			It("should return an emtpy network", func() {
				Expect(networks.VipNetwork()).To(Equal(Network{}))
			})
		})
	})

	Describe("DNS", func() {
		It("should return only the dynamic network DNS servers", func() {
			Expect(networks.DNS()).To(Equal([]string{"fake-dynamic-network-dns"}))
		})
	})

	Describe("Network", func() {
		It("should return only the dynamic network network", func() {
			Expect(networks.Network()).To(Equal("fake-dynamic-network-network"))
		})

		Context("when Network Name is empty", func() {
			BeforeEach(func() {
				dynamicNetwork.Network = ""
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-vip-network":     vipNetwork,
				}
			})

			It("should return the default network name", func() {
				Expect(networks.Network()).To(BeEmpty())
			})
		})
	})
})
