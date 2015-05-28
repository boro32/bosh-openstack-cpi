package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"
)

var _ = Describe("HasVM", func() {
	var (
		err   error
		found bool

		serverService *serverfakes.FakeServerService

		hasVM HasVM
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		hasVM = NewHasVM(serverService)
	})

	Describe("Run", func() {
		It("returns true if vm ID exist", func() {
			serverService.FindFound = true

			found, err = hasVM.Run("fake-server-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeTrue())
			Expect(serverService.FindCalled).To(BeTrue())
		})

		It("returns false if vm ID does not exist", func() {
			serverService.FindFound = false

			found, err = hasVM.Run("fake-server-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(found).To(BeFalse())
			Expect(serverService.FindCalled).To(BeTrue())
		})

		It("returns an error if serverService find call returns an error", func() {
			serverService.FindErr = errors.New("fake-server-service-error")

			_, err = hasVM.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.FindCalled).To(BeTrue())
		})
	})
})
