package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"
	volumefakes "github.com/frodenas/bosh-openstack-cpi/openstack/volume_service/fakes"
	volumetypefakes "github.com/frodenas/bosh-openstack-cpi/openstack/volume_type_service/fakes"

	"github.com/frodenas/bosh-openstack-cpi/api"
)

var _ = Describe("CreateDisk", func() {
	var (
		err        error
		diskCID    DiskCID
		vmCID      VMCID
		cloudProps DiskCloudProperties

		volumeService                *volumefakes.FakeVolumeService
		volumeTypeService            *volumetypefakes.FakeVolumeTypeService
		serverService                *serverfakes.FakeServerService
		ignoreServerAvailabilityZone bool

		createDisk CreateDisk
	)

	BeforeEach(func() {
		volumeService = &volumefakes.FakeVolumeService{}
		volumeTypeService = &volumetypefakes.FakeVolumeTypeService{}
		serverService = &serverfakes.FakeServerService{}
		createDisk = NewCreateDisk(volumeService, volumeTypeService, serverService, ignoreServerAvailabilityZone)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			vmCID = ""
			cloudProps = DiskCloudProperties{}
			volumeService.CreateID = "fake-volume-id"
		})

		It("creates the disk", func() {
			diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.FindCalled).To(BeFalse())
			Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
			Expect(volumeService.CreateCalled).To(BeTrue())
			Expect(volumeService.CreateSize).To(Equal(32))
			Expect(volumeService.CreateVolumeType).To(BeEmpty())
			Expect(volumeService.CreateAvailabilityZone).To(BeEmpty())
			Expect(diskCID).To(Equal(DiskCID("fake-volume-id")))
		})

		It("returns an error if volumeService create call returns an error", func() {
			volumeService.CreateErr = errors.New("fake-volume-service-error")

			_, err = createDisk.Run(32768, cloudProps, vmCID)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-volume-service-error"))
			Expect(serverService.FindCalled).To(BeFalse())
			Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
			Expect(volumeService.CreateCalled).To(BeTrue())
		})

		Context("when vmCID is set", func() {
			BeforeEach(func() {
				vmCID = VMCID("fake-server-cid")
			})

			It("creates the disk at the vm zone", func() {
				serverService.FindFound = true

				diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).NotTo(HaveOccurred())
				Expect(serverService.FindCalled).To(BeTrue())
				Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
				Expect(volumeService.CreateCalled).To(BeTrue())
				Expect(volumeService.CreateSize).To(Equal(32))
				Expect(volumeService.CreateVolumeType).To(BeEmpty())
				//Expect(volumeService.CreateAvailabilityZone).To(Equal("fake-server-zone"))
				Expect(diskCID).To(Equal(DiskCID("fake-volume-id")))
			})

			It("returns an error if serverService find call returns an error", func() {
				serverService.FindErr = errors.New("fake-server-service-error")

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
				Expect(serverService.FindCalled).To(BeTrue())
				Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
				Expect(volumeService.CreateCalled).To(BeFalse())
			})

			It("returns an error if server is not found", func() {
				serverService.FindFound = false

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(api.NewVMNotFoundError(string(vmCID)).Error()))
				Expect(serverService.FindCalled).To(BeTrue())
				Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
				Expect(volumeService.CreateCalled).To(BeFalse())
			})

			Context("when ignoreServerAvailabilityZone is set to true", func() {
				BeforeEach(func() {
					ignoreServerAvailabilityZone = true
					createDisk = NewCreateDisk(volumeService, volumeTypeService, serverService, ignoreServerAvailabilityZone)
				})

				It("does not check the vm's availability zones", func() {
					diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
					Expect(err).NotTo(HaveOccurred())
					Expect(serverService.FindCalled).To(BeFalse())
					Expect(volumeTypeService.FindByNameCalled).To(BeFalse())
					Expect(volumeService.CreateCalled).To(BeTrue())
					Expect(volumeService.CreateSize).To(Equal(32))
					Expect(volumeService.CreateVolumeType).To(BeEmpty())
					Expect(diskCID).To(Equal(DiskCID("fake-volume-id")))
				})
			})
		})

		Context("when volume type is set", func() {
			BeforeEach(func() {
				cloudProps = DiskCloudProperties{VolumeType: "fake-volume-type"}
			})

			It("creates the disk using the appropiate volume type", func() {
				volumeTypeService.FindByNameFound = true

				diskCID, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).NotTo(HaveOccurred())
				Expect(serverService.FindCalled).To(BeFalse())
				Expect(volumeTypeService.FindByNameCalled).To(BeTrue())
				Expect(volumeService.CreateCalled).To(BeTrue())
				Expect(volumeService.CreateSize).To(Equal(32))
				Expect(volumeService.CreateVolumeType).To(Equal("fake-volume-type"))
				Expect(volumeService.CreateAvailabilityZone).To(BeEmpty())
				Expect(diskCID).To(Equal(DiskCID("fake-volume-id")))
			})

			It("returns an error if volumeTypeService find call returns an error", func() {
				volumeTypeService.FindByNameErr = errors.New("fake-volume-type-service-error")

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-volume-type-service-error"))
				Expect(serverService.FindCalled).To(BeFalse())
				Expect(volumeTypeService.FindByNameCalled).To(BeTrue())
				Expect(volumeService.CreateCalled).To(BeFalse())
			})

			It("returns an error if volume type is not found", func() {
				volumeTypeService.FindFound = false

				_, err = createDisk.Run(32768, cloudProps, vmCID)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Volume Type 'fake-volume-type' does not exists"))
				Expect(serverService.FindCalled).To(BeFalse())
				Expect(volumeTypeService.FindByNameCalled).To(BeTrue())
				Expect(volumeService.CreateCalled).To(BeFalse())
			})
		})
	})

})
