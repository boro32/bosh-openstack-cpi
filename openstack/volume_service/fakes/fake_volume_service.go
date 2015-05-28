package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"
)

type FakeVolumeService struct {
	CreateCalled           bool
	CreateErr              error
	CreateID               string
	CreateSize             int
	CreateVolumeType       string
	CreateAvailabilityZone string

	DeleteCalled bool
	DeleteErr    error

	FindCalled bool
	FindFound  bool
	FindVolume volume.Volume
	FindErr    error
}

func (v *FakeVolumeService) Create(size int, volumeType string, availabilityZone string) (string, error) {
	v.CreateCalled = true
	v.CreateSize = size
	v.CreateVolumeType = volumeType
	v.CreateAvailabilityZone = availabilityZone
	return v.CreateID, v.CreateErr
}

func (v *FakeVolumeService) Delete(id string) error {
	v.DeleteCalled = true
	return v.DeleteErr
}

func (v *FakeVolumeService) Find(id string) (volume.Volume, bool, error) {
	v.FindCalled = true
	return v.FindVolume, v.FindFound, v.FindErr
}
