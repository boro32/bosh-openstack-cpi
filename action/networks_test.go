package action_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("Networks", func() {
	var (
		networks Networks
	)

	BeforeEach(func() {
		networks = Networks{
			"fake-network-1-name": Network{
				Type:    "fake-network-1-type",
				IP:      "fake-network-1-ip",
				Gateway: "fake-network-1-gateway",
				Netmask: "fake-network-1-netmask",
				DNS:     []string{"fake-network-1-dns"},
				Default: []string{"fake-network-1-default"},
				CloudProperties: NetworkCloudProperties{
					Network:        "fake-network-1-cloud-network",
					SecurityGroups: []string{"fake-network-1-security-group"},
				},
			},
			"fake-network-2-name": Network{
				Type: "fake-network-2-type",
				IP:   "fake-network-2-ip",
			},
		}
	})

	Describe("AsServerServiceNetworks", func() {
		It("returns networks for the server service", func() {
			expectedInstanceNetworks := server.Networks{
				"fake-network-1-name": server.Network{
					Type:           "fake-network-1-type",
					IP:             "fake-network-1-ip",
					Gateway:        "fake-network-1-gateway",
					Netmask:        "fake-network-1-netmask",
					DNS:            []string{"fake-network-1-dns"},
					Default:        []string{"fake-network-1-default"},
					Network:        "fake-network-1-cloud-network",
					SecurityGroups: []string{"fake-network-1-security-group"},
				},
				"fake-network-2-name": server.Network{
					Type: "fake-network-2-type",
					IP:   "fake-network-2-ip",
				},
			}

			Expect(networks.AsServerServiceNetworks()).To(Equal(expectedInstanceNetworks))
		})
	})

	Describe("AsRegistryNetworks", func() {
		It("returns networks for the registry", func() {
			expectedRegistryNetworks := registry.NetworksSettings{
				"fake-network-1-name": registry.NetworkSettings{
					Type:    "fake-network-1-type",
					IP:      "fake-network-1-ip",
					Gateway: "fake-network-1-gateway",
					Netmask: "fake-network-1-netmask",
					DNS:     []string{"fake-network-1-dns"},
					Default: []string{"fake-network-1-default"},
				},
				"fake-network-2-name": registry.NetworkSettings{
					Type: "fake-network-2-type",
					IP:   "fake-network-2-ip",
				},
			}

			Expect(networks.AsRegistryNetworks()).To(Equal(expectedRegistryNetworks))
		})
	})
})
