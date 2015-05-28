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

var _ = Describe("AttachDisk", func() {
	var (
		err                   error
		expectedAgentSettings registry.AgentSettings

		serverService  *serverfakes.FakeServerService
		registryClient *registryfakes.FakeClient

		attachDisk AttachDisk
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		registryClient = &registryfakes.FakeClient{}
		attachDisk = NewAttachDisk(serverService, registryClient)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			serverService.AttachVolumeDeviceName = "fake-volume-device-name"
			serverService.AttachVolumeDevicePath = "fake-volume-device-path"
			registryClient.FetchSettings = registry.AgentSettings{}
			expectedAgentSettings = registry.AgentSettings{
				Disks: registry.DisksSettings{
					Persistent: map[string]registry.PersistentSettings{
						"fake-volume-id": {
							ID:       "fake-volume-id",
							VolumeID: "fake-volume-device-name",
							Path:     "fake-volume-device-path",
						},
					},
				},
			}
		})

		It("attaches the disk", func() {
			_, err = attachDisk.Run("fake-server-id", "fake-volume-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
			Expect(registryClient.UpdateSettings).To(Equal(expectedAgentSettings))
		})

		It("returns an error if serverService attach disk call returns an error", func() {
			serverService.AttachVolumeErr = errors.New("fake-server-service-error")

			_, err = attachDisk.Run("fake-server-id", "fake-volume-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeFalse())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient fetch call returns an error", func() {
			registryClient.FetchErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-server-id", "fake-volume-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(serverService.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeFalse())
		})

		It("returns an error if registryClient update call returns an error", func() {
			registryClient.UpdateErr = errors.New("fake-registry-client-error")

			_, err = attachDisk.Run("fake-server-id", "fake-volume-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-registry-client-error"))
			Expect(serverService.AttachVolumeCalled).To(BeTrue())
			Expect(registryClient.FetchCalled).To(BeTrue())
			Expect(registryClient.UpdateCalled).To(BeTrue())
		})
	})
})
