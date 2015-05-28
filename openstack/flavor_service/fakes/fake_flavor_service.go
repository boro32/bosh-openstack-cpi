package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
)

type FakeFlavorService struct {
	FindCalled bool
	FindFound  bool
	FindFlavor flavor.Flavor
	FindErr    error
}

func (f *FakeFlavorService) Find(id string) (flavor.Flavor, bool, error) {
	f.FindCalled = true
	return f.FindFlavor, f.FindFound, f.FindErr
}
