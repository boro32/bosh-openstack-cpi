package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
)

type FakeFlavorService struct {
	FindCalled bool
	FindFound  bool
	FindFlavor flavor.Flavor
	FindErr    error

	FindByNameCalled bool
	FindByNameFound  bool
	FindByNameFlavor flavor.Flavor
	FindByNameErr    error
}

func (f *FakeFlavorService) Find(id string) (flavor.Flavor, bool, error) {
	f.FindCalled = true
	return f.FindFlavor, f.FindFound, f.FindErr
}

func (f *FakeFlavorService) FindByName(name string) (flavor.Flavor, bool, error) {
	f.FindByNameCalled = true
	return f.FindByNameFlavor, f.FindByNameFound, f.FindByNameErr
}
