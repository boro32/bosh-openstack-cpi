package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	snapshotfakes "github.com/frodenas/bosh-openstack-cpi/openstack/snapshot_service/fakes"
)

var _ = Describe("SnapshotDisk", func() {
	var (
		err        error
		metadata   SnapshotMetadata
		snapshotID SnapshotCID

		snapshotService *snapshotfakes.FakeSnapshotService

		snapshotDisk SnapshotDisk
	)

	BeforeEach(func() {
		snapshotService = &snapshotfakes.FakeSnapshotService{}
		snapshotDisk = NewSnapshotDisk(snapshotService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			snapshotService.CreateID = "fake-snapshot-id"
			metadata = SnapshotMetadata{Deployment: "fake-deployment", Job: "fake-job", Index: "fake-index"}
		})

		Context("creates a snaphot", func() {
			It("with the proper description", func() {
				snapshotID, err = snapshotDisk.Run("fake-volume-id", metadata)
				Expect(err).NotTo(HaveOccurred())
				Expect(snapshotService.CreateCalled).To(BeTrue())
				Expect(snapshotService.CreateVolumeID).To(Equal("fake-volume-id"))
				Expect(snapshotService.CreateDescription).To(Equal("fake-deployment/fake-job/fake-index"))
				Expect(snapshotID).To(Equal(SnapshotCID("fake-snapshot-id")))
			})

			Context("when metadata is empty", func() {
				BeforeEach(func() {
					metadata = SnapshotMetadata{}
				})

				It("with an empty description", func() {
					snapshotID, err = snapshotDisk.Run("fake-volume-id", metadata)
					Expect(err).NotTo(HaveOccurred())
					Expect(snapshotService.CreateCalled).To(BeTrue())
					Expect(snapshotService.CreateVolumeID).To(Equal("fake-volume-id"))
					Expect(snapshotService.CreateDescription).To(BeEmpty())
					Expect(snapshotID).To(Equal(SnapshotCID("fake-snapshot-id")))
				})
			})
		})

		It("returns an error if snapshotService create call returns an error", func() {
			snapshotService.CreateErr = errors.New("fake-snapshot-service-error")

			_, err = snapshotDisk.Run("fake-volume-id", metadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-snapshot-service-error"))
			Expect(snapshotService.CreateCalled).To(BeTrue())
		})
	})
})
