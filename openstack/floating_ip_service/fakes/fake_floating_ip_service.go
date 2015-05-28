package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/floating_ip_service"
)

type FakeFloatingIPService struct {
	FindCalled     bool
	FindFound      bool
	FindFloatingIP floatingip.FloatingIP
	FindErr        error

	FindByIPCalled     bool
	FindByIPFound      bool
	FindByIPFloatingIP floatingip.FloatingIP
	FindByIPErr        error
}

func (fip *FakeFloatingIPService) Find(id string) (floatingip.FloatingIP, bool, error) {
	fip.FindCalled = true
	return fip.FindFloatingIP, fip.FindFound, fip.FindErr
}

func (fip *FakeFloatingIPService) FindByIP(ipAddress string) (floatingip.FloatingIP, bool, error) {
	fip.FindByIPCalled = true
	return fip.FindByIPFloatingIP, fip.FindByIPFound, fip.FindByIPErr
}
