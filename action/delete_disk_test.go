package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	volumefakes "github.com/frodenas/bosh-openstack-cpi/openstack/volume_service/fakes"
)

var _ = Describe("DeleteDisk", func() {
	var (
		err error

		volumeService *volumefakes.FakeVolumeService

		deleteDisk DeleteDisk
	)

	BeforeEach(func() {
		volumeService = &volumefakes.FakeVolumeService{}
		deleteDisk = NewDeleteDisk(volumeService)
	})

	Describe("Run", func() {
		It("deletes the disk", func() {
			_, err = deleteDisk.Run("fake-volume-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(volumeService.DeleteCalled).To(BeTrue())
		})

		It("returns an error if volumeService delete call returns an error", func() {
			volumeService.DeleteErr = errors.New("fake-volume-service-error")

			_, err = deleteDisk.Run("fake-volume-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-volume-service-error"))
			Expect(volumeService.DeleteCalled).To(BeTrue())
		})
	})
})
