package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"

	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

var _ = Describe("GetDisks", func() {
	var (
		err                 error
		attachedVolumesList []string
		disks               server.AttachedVolumes

		serverService *serverfakes.FakeServerService

		getDisks GetDisks
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		getDisks = NewGetDisks(serverService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			attachedVolumesList = []string{"fake-volume-1", "fake-volume-2"}
			serverService.AttachedVolumesList = server.AttachedVolumes(attachedVolumesList)
		})

		It("returns the list of attached disks", func() {
			disks, err = getDisks.Run("fake-server-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.AttachedVolumesCalled).To(BeTrue())
			Expect(disks).To(Equal(server.AttachedVolumes(attachedVolumesList)))
		})

		Context("when there are not any attached disks", func() {
			BeforeEach(func() {
				serverService.AttachedVolumesList = server.AttachedVolumes{}
			})

			It("returns an empty array", func() {
				disks, err = getDisks.Run("fake-server-id")
				Expect(err).NotTo(HaveOccurred())
				Expect(serverService.AttachedVolumesCalled).To(BeTrue())
				Expect(disks).To(BeEmpty())
			})
		})

		It("returns an error if vmService attached disks call returns an error", func() {
			serverService.AttachedVolumesErr = errors.New("fake-server-service-error")

			_, err = getDisks.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.AttachedVolumesCalled).To(BeTrue())
		})
	})
})
