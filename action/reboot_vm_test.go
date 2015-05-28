package action_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/frodenas/bosh-openstack-cpi/action"

	serverfakes "github.com/frodenas/bosh-openstack-cpi/openstack/server_service/fakes"
)

var _ = Describe("RebootVM", func() {
	var (
		err error

		serverService *serverfakes.FakeServerService

		rebootVM RebootVM
	)

	BeforeEach(func() {
		serverService = &serverfakes.FakeServerService{}
		rebootVM = NewRebootVM(serverService)
	})

	Describe("Run", func() {
		It("reboots the vm", func() {
			_, err = rebootVM.Run("fake-server-id")
			Expect(err).NotTo(HaveOccurred())
			Expect(serverService.RebootCalled).To(BeTrue())
		})

		It("returns an error if serverService reboot call returns an error", func() {
			serverService.RebootErr = errors.New("fake-server-service-error")

			_, err = rebootVM.Run("fake-server-id")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("fake-server-service-error"))
			Expect(serverService.RebootCalled).To(BeTrue())
		})
	})
})
