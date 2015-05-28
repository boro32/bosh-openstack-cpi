package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"

	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

var _ = Describe("SetVMMetadata", func() {
	var (
		err        error
		vmMetadata VMMetadata

		serverService *serverfakes.FakeServerService

		setVMMetadata SetVMMetadata
	)

	BeforeEach(func() {
		vmMetadata = map[string]interface{}{
			"deployment": "fake-deployment",
			"job":        "fake-job",
			"index":      "fake-index",
		}
		serverService = &serverfakes.FakeServerService{}
		setVMMetadata = NewSetVMMetadata(serverService)
	})

	Describe("Run", func() {
		It("set the vm metadata", func() {
			_, err = setVMMetadata.Run("fake-server-id", vmMetadata)
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.SetMetadataCalled).To(BeTrue())
			Expect(serverService.SetMetadataServerMetadata).To(Equal(server.Metadata(vmMetadata)))
		})

		It("returns an error if serverService set metadata call returns an error", func() {
			serverService.SetMetadataErr = errors.New("fake-server-service-error")

			_, err = setVMMetadata.Run("fake-server-id", vmMetadata)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.SetMetadataCalled).To(BeTrue())
		})
	})
})
