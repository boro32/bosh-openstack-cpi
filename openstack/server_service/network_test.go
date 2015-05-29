package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

var _ = Describe("Network", func() {
	var (
		err            error
		dynamicNetwork Network
		manualNetwork  Network
		vipNetwork     Network
		unknownNetwork Network
	)

	BeforeEach(func() {
		dynamicNetwork = Network{
			Type:           "dynamic",
			IP:             "fake-dynamic-network-ip",
			Gateway:        "fake-dynamic-network-gateway",
			Netmask:        "fake-dynamic-network-netmask",
			DNS:            []string{"fake-dynamic-network-dns"},
			Default:        []string{"fake-dynamic-network-default"},
			Network:        "fake-dynamic-network-network",
			SecurityGroups: []string{"fake-dynamic-network-security-group"},
		}

		manualNetwork = Network{
			Type:           "manual",
			IP:             "fake-manual-network-ip",
			Gateway:        "fake-manual-network-gateway",
			Netmask:        "fake-manual-network-netmask",
			DNS:            []string{"fake-manual-network-dns"},
			Default:        []string{"fake-manual-network-default"},
			Network:        "fake-manual-network-network",
			SecurityGroups: []string{"fake-manual-network-security-group"},
		}

		vipNetwork = Network{
			Type:           "vip",
			IP:             "fake-vip-network-ip",
			Gateway:        "fake-vip-network-gateway",
			Netmask:        "fake-vip-network-netmask",
			DNS:            []string{"fake-vip-network-dns"},
			Default:        []string{"fake-vip-network-default"},
			Network:        "fake-vip-network-network",
			SecurityGroups: []string{"fake-vip-network-security-group"},
		}

		unknownNetwork = Network{Type: "unknown"}
	})

	Describe("IsDynamic", func() {
		It("returns true for a dynamic network", func() {
			Expect(dynamicNetwork.IsDynamic()).To(BeTrue())
		})

		It("returns false for a manual network", func() {
			Expect(manualNetwork.IsDynamic()).To(BeFalse())
		})

		It("returns false for a vip network", func() {
			Expect(vipNetwork.IsDynamic()).To(BeFalse())
		})
	})

	Describe("IsManual", func() {
		It("returns false for a dynamic network", func() {
			Expect(dynamicNetwork.IsManual()).To(BeFalse())
		})

		It("returns true for a manual network", func() {
			Expect(manualNetwork.IsManual()).To(BeTrue())
		})

		It("returns false for a vip network", func() {
			Expect(vipNetwork.IsManual()).To(BeFalse())
		})
	})

	Describe("IsVip", func() {
		It("returns false for a dynamic network", func() {
			Expect(dynamicNetwork.IsVip()).To(BeFalse())
		})

		It("returns false for a manual network", func() {
			Expect(manualNetwork.IsDynamic()).To(BeFalse())
		})

		It("returns true for a vip network", func() {
			Expect(vipNetwork.IsVip()).To(BeTrue())
		})
	})

	Describe("Validate", func() {
		Context("Dynamic Network", func() {
			It("does not return error if network properties are valid", func() {
				err = dynamicNetwork.Validate()
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("Manual Network", func() {
			It("does not return error if network properties are valid", func() {
				err = manualNetwork.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error if does not have a network name", func() {
				manualNetwork.Network = ""

				err = manualNetwork.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Manual Networks must provide a Network name"))
			})
		})

		Context("VIP Network", func() {
			It("does not return error if network properties are valid", func() {
				err = vipNetwork.Validate()
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns an error if does not have n IP Address", func() {
				vipNetwork.IP = ""

				err = vipNetwork.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("VIP Networks must provide an IP Address"))
			})
		})

		Context("Unknown Network", func() {
			It("returns an error", func() {
				err = unknownNetwork.Validate()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Network type 'unknown' not supported"))
			})
		})
	})

	Describe("IPAddress", func() {
		It("returns an empty IP address for a dynamic network", func() {
			Expect(dynamicNetwork.IPAddress()).To(BeEmpty())
		})

		It("returns the network IP address for a manual network", func() {
			Expect(manualNetwork.IPAddress()).To(Equal("fake-manual-network-ip"))
		})

		It("returns network IP address for a vip network", func() {
			Expect(vipNetwork.IPAddress()).To(Equal("fake-vip-network-ip"))
		})
	})

	Describe("DNSList", func() {
		It("returns the network DNS list for a dynamic network", func() {
			Expect(dynamicNetwork.DNSList()).To(Equal([]string{"fake-dynamic-network-dns"}))
		})

		It("returns the network DNS list for a manual network", func() {
			Expect(manualNetwork.DNSList()).To(Equal([]string{"fake-manual-network-dns"}))
		})

		It("returns an empty DNS list for a vip network", func() {
			Expect(vipNetwork.DNSList()).To(BeEmpty())
		})
	})

	Describe("NetworkName", func() {
		It("returns the network name for a dynamic network", func() {
			Expect(dynamicNetwork.NetworkName()).To(Equal("fake-dynamic-network-network"))
		})

		It("returns the network name for a manual network", func() {
			Expect(manualNetwork.NetworkName()).To(Equal("fake-manual-network-network"))
		})

		It("returns an empty network name for a vip network", func() {
			Expect(vipNetwork.NetworkName()).To(BeEmpty())
		})
	})

	Describe("SecurityGroupsList", func() {
		It("returns the network security groups list for a dynamic network", func() {
			Expect(dynamicNetwork.SecurityGroupsList()).To(Equal([]string{"fake-dynamic-network-security-group"}))
		})

		It("returns the network security groups list for a manual network", func() {
			Expect(manualNetwork.SecurityGroupsList()).To(Equal([]string{"fake-manual-network-security-group"}))
		})

		It("returns an empty security groups list for a vip network", func() {
			Expect(vipNetwork.SecurityGroupsList()).To(BeEmpty())
		})
	})
})
