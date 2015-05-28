package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"

	registryfakes "github.com/frodenas/bosh-registry/client/fakes"

	"github.com/frodenas/bosh-registry/client"
)

var _ = Describe("ConfigureNetworks", func() {
	var (
		err                   error
		networks              Networks
		expectedAgentSettings registry.AgentSettings

		serverService  *serverfakes.FakeServerService
		registryClient *registryfakes.FakeClient

		configureNetworks ConfigureNetworks
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		registryClient = &registryfakes.FakeClient{}
		configureNetworks = NewConfigureNetworks(serverService, registryClient)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			networks = Networks{
				"fake-network-name": Network{
					Type:    "dynamic",
					IP:      "fake-network-ip",
					Gateway: "fake-network-gateway",
					Netmask: "fake-network-netmask",
					DNS:     []string{"fake-network-dns"},
					Default: []string{"fake-network-default"},
					CloudProperties: NetworkCloudProperties{
						Network:        "fake-network-cloud-network",
						SecurityGroups: []string{"fake-network-security-group"},
					},
				},
			}
			registryClient.FetchSettings = registry.AgentSettings{}
			expectedAgentSettings = registry.AgentSettings{
				Networks: registry.NetworksSettings{
					"fake-network-name": registry.NetworkSettings{
						Type:    "dynamic",
						IP:      "fake-network-ip",
						Gateway: "fake-network-gateway",
						Netmask: "fake-network-netmask",
						DNS:     []string{"fake-network-dns"},
						Default: []string{"fake-network-default"},
					},
				},
			}
		})

		It("configures the network", func() {
			_, err = configureNetworks.Run("fake-server-id", networks)
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
		})

		It("returns an error if networks are not valid", func() {
			_, err = configureNetworks.Run("fake-server-id", Networks{})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Configuring networks for vm"))
			Expect(serverService.UpdateNetworkConfigurationCalled).To(BeFalse())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if serverService update network configuration call returns an error", func() {
			serverService.UpdateNetworkConfigurationErr = errors.New("fake-server-service-error")

			_, err = configureNetworks.Run("fake-server-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = configureNetworks.Run("fake-server-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(serverService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = configureNetworks.Run("fake-server-id", networks)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(serverService.UpdateNetworkConfigurationCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
