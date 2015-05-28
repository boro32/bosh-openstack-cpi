package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_type_service"
)

type FakeVolumeTypeService struct {
	FindCalled     bool
	FindFound      bool
	FindVolumeType volumetype.VolumeType
	FindErr        error

	FindByNameCalled     bool
	FindByNameFound      bool
	FindByNameVolumeType volumetype.VolumeType
	FindByNameErr        error
}

func (vt *FakeVolumeTypeService) Find(id string) (volumetype.VolumeType, bool, error) {
	vt.FindCalled = true
	return vt.FindVolumeType, vt.FindFound, vt.FindErr
}

func (vt *FakeVolumeTypeService) FindByName(name string) (volumetype.VolumeType, bool, error) {
	vt.FindByNameCalled = true
	return vt.FindByNameVolumeType, vt.FindByNameFound, vt.FindByNameErr
}
