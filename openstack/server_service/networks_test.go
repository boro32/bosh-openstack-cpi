package server_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

var _ = Describe("Networks", func() {
	var (
		err            error
		dynamicNetwork Network
		manualNetwork  Network
		vipNetwork     Network
		networks       Networks

		expectedDNSList            []string
		expectedSecurityGroupsList []string
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

		networks = Networks{
			"fake-dynamic-network": dynamicNetwork,
			"fake-manual-network":  manualNetwork,
			"fake-vip-network":     vipNetwork,
		}

		expectedDNSList = []string{
			"fake-dynamic-network-dns",
			"fake-manual-network-dns",
		}
		expectedSecurityGroupsList = []string{
			"fake-dynamic-network-security-group",
			"fake-manual-network-security-group",
		}
	})

	Describe("Validate", func() {
		It("does not return an error if networks are valid", func() {
			err = networks.Validate()
			Expect(err).NotTo(HaveOccurred())
		})

		It("returns an error if networks are not valid", func() {
			networks = Networks{"fake-network-name": Network{Type: "unknown"}}

			err = networks.Validate()
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if there are not any dynamic and manual networks", func() {
			networks = Networks{
				"fake-vip-network": vipNetwork,
			}

			err = networks.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("At least one 'dynamic' or 'manual' network should be defined"))
		})

		It("returns an error if there is more than 1 VIP network", func() {
			networks = Networks{
				"fake-dynamic-network": dynamicNetwork,
				"fake-vip-network-1":   vipNetwork,
				"fake-vip-network-2":   vipNetwork,
			}

			err = networks.Validate()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Only one VIP network is allowed"))
		})
	})

	Describe("DNSList", func() {
		It("returns the DNS list from the dynamic and manual networks", func() {
			Expect(networks.DNSList()).To(ConsistOf(expectedDNSList))
		})

		Context("when there are duplicates", func() {
			BeforeEach(func() {
				dynamicNetwork.DNS = []string{
					"fake-dynamic-network-dns",
					"fake-duplicate-network-dns",
				}
				manualNetwork.DNS = []string{
					"fake-manual-network-dns",
					"fake-duplicate-network-dns",
				}
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-manual-network":  manualNetwork,
				}
				expectedDNSList = []string{
					"fake-dynamic-network-dns",
					"fake-manual-network-dns",
					"fake-duplicate-network-dns",
				}
			})

			It("returns the DNS list without duplicates", func() {
				Expect(networks.DNSList()).To(ConsistOf(expectedDNSList))
			})
		})
	})

	Describe("SecurityGroupsList", func() {
		It("returns the Security Groups list from the dynamic and manual networks", func() {
			Expect(networks.SecurityGroupsList()).To(ConsistOf(expectedSecurityGroupsList))
		})

		Context("when there are duplicates", func() {
			BeforeEach(func() {
				dynamicNetwork.SecurityGroups = []string{
					"fake-dynamic-network-security-group",
					"fake-duplicate-network-security-group",
				}
				manualNetwork.SecurityGroups = []string{
					"fake-manual-network-security-group",
					"fake-duplicate-network-security-group",
				}
				networks = Networks{
					"fake-dynamic-network": dynamicNetwork,
					"fake-manual-network":  manualNetwork,
				}
				expectedSecurityGroupsList = []string{
					"fake-dynamic-network-security-group",
					"fake-manual-network-security-group",
					"fake-duplicate-network-security-group",
				}
			})

			It("returns the Security Groups list without duplicates", func() {
				Expect(networks.SecurityGroupsList()).To(ConsistOf(expectedSecurityGroupsList))
			})
		})
	})
})
