package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	imagefakes "github.com/frodenas/bosh-openstack-cpi/openstack/image_service/fakes"
)

var _ = Describe("DeleteStemcell", func() {
	var (
		err error

		imageService *imagefakes.FakeImageService

		deleteStemcell DeleteStemcell
	)

	BeforeEach(func() {
		imageService = &imagefakes.FakeImageService{}
		deleteStemcell = NewDeleteStemcell(imageService)
	})

	Describe("Run", func() {
		It("deletes the stemcell", func() {
			_, err = deleteStemcell.Run("fake-image-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(imageService.DeleteCalled).To(BeTrue())
		})

		It("returns an error if imageService delete call returns an error", func() {
			imageService.DeleteErr = errors.New("fake-image-service-error")

			_, err = deleteStemcell.Run("fake-image-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-image-service-error"))
			Expect(imageService.DeleteCalled).To(BeTrue())
		})
	})
})
