package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	imagefakes "github.com/frodenas/bosh-openstack-cpi/openstack/image_service/fakes"
)

var _ = Describe("CreateStemcell", func() {
	var (
		err         error
		stemcellCID StemcellCID
		cloudProps  StemcellCloudProperties

		imageService   *imagefakes.FakeImageService
		createStemcell CreateStemcell
	)

	BeforeEach(func() {
		imageService = &imagefakes.FakeImageService{}
		createStemcell = NewCreateStemcell(imageService)
	})

	Describe("Run", func() {
		BeforeEach(func() {
			cloudProps = StemcellCloudProperties{
				Name:           "fake-stemcell-name",
				Version:        "fake-stemcell-version",
				Infrastructure: "openstack",
				ImageUUID:      "",
			}
		})

		Context("when infrastructure is not openstack", func() {
			BeforeEach(func() {
				cloudProps.Infrastructure = "fake-insfrastructure"
			})

			It("returns an error", func() {
				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Invalid 'fake-insfrastructure' infrastructure"))
				Expect(imageService.CreateCalled).To(BeFalse())
			})
		})

		Context("when cloud properties has an Image UUID", func() {
			BeforeEach(func() {
				cloudProps.ImageUUID = "fake-image-uuid"
			})

			It("returns the Image UUID", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(imageService.CreateCalled).To(BeFalse())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-image-uuid")))
			})
		})

		Context("from a stemcell tarball", func() {
			BeforeEach(func() {
				imageService.CreateID = "fake-image-id"
			})

			It("creates the stemcell", func() {
				stemcellCID, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).NotTo(HaveOccurred())
				Expect(imageService.CreateCalled).To(BeTrue())
				Expect(stemcellCID).To(Equal(StemcellCID("fake-image-id")))
				Expect(imageService.CreateImagePath).To(Equal("fake-stemcell-tarball"))
				Expect(imageService.CreateDescription).To(Equal("fake-stemcell-name/fake-stemcell-version"))
			})

			It("returns an error if imageService create call returns an error", func() {
				imageService.CreateErr = errors.New("fake-image-service-error")

				_, err = createStemcell.Run("fake-stemcell-tarball", cloudProps)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("fake-image-service-error"))
				Expect(imageService.CreateCalled).To(BeTrue())
			})
		})
	})
})
