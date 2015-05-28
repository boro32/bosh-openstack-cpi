package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"

	registryfakes "github.com/frodenas/bosh-registry/client/fakes"
)

var _ = Describe("DeleteVM", func() {
	var (
		err error

		serverService  *serverfakes.FakeServerService
		registryClient *registryfakes.FakeClient

		deleteVM DeleteVM
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		registryClient = &registryfakes.FakeClient{}
		deleteVM = NewDeleteVM(serverService, registryClient)
	})

	Describe("Run", func() {
		It("deletes the vm", func() {
			_, err = deleteVM.Run("fake-server-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(serverService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeTrue())
		})

		It("returns an error if serverService delete network configuration call returns an error", func() {
			serverService.DeleteNetworkConfigurationErr = errors.New("fake-server-service-error")

			_, err = deleteVM.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(serverService.DeleteCalled).To(BeFalse())
			Expect(registryClient.DeleteCalled).To(BeFalse())
		})

		It("returns an error if serverService delete call returns an error", func() {
			serverService.DeleteErr = errors.New("fake-server-service-error")

			_, err = deleteVM.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(serverService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeFalse())
		})

		It("returns an error if registryClient delete call returns an error", func() {
			registryClient.DeleteErr = errors.New("fake-registry-client-error")

			_, err = deleteVM.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(serverService.DeleteNetworkConfigurationCalled).To(BeTrue())
			Expect(serverService.DeleteCalled).To(BeTrue())
			Expect(registryClient.DeleteCalled).To(BeTrue())
		})
	})
})
