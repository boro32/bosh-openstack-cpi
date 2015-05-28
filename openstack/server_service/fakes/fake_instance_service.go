package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type FakeServerService struct {
	AddNetworkConfigurationCalled bool
	AddNetworkConfigurationErr    error

	AttachVolumeCalled     bool
	AttachVolumeErr        error
	AttachVolumeDeviceName string
	AttachVolumeDevicePath string

	AttachedVolumesCalled bool
	AttachedVolumesErr    error
	AttachedVolumesList   server.AttachedVolumes

	CleanUpCalled bool

	CreateCalled           bool
	CreateErr              error
	CreateID               string
	CreateServerProps      *server.Properties
	CreateNetworks         server.Networks
	CreateRegistryEndpoint string

	DeleteCalled bool
	DeleteErr    error

	DeleteNetworkConfigurationCalled bool
	DeleteNetworkConfigurationErr    error

	DetachVolumeCalled bool
	DetachVolumeErr    error

	FindCalled bool
	FindFound  bool
	FindServer *servers.Server
	FindErr    error

	RebootCalled bool
	RebootErr    error

	SetMetadataCalled         bool
	SetMetadataErr            error
	SetMetadataServerMetadata server.Metadata

	SetTagsCalled bool
	SetTagsErr    error

	UpdateNetworkConfigurationCalled bool
	UpdateNetworkConfigurationErr    error
}

func (i *FakeServerService) AddNetworkConfiguration(id string, networks server.Networks) error {
	i.AddNetworkConfigurationCalled = true
	return i.AddNetworkConfigurationErr
}

func (i *FakeServerService) AttachVolume(id string, volumeID string) (string, string, error) {
	i.AttachVolumeCalled = true
	return i.AttachVolumeDeviceName, i.AttachVolumeDevicePath, i.AttachVolumeErr
}

func (i *FakeServerService) AttachedVolumes(id string) (server.AttachedVolumes, error) {
	i.AttachedVolumesCalled = true
	return i.AttachedVolumesList, i.AttachedVolumesErr
}

func (i *FakeServerService) CleanUp(id string) {
	i.CleanUpCalled = true
	return
}

func (i *FakeServerService) Create(serverProps *server.Properties, networks server.Networks, registryEndpoint string) (string, error) {
	i.CreateCalled = true
	i.CreateServerProps = serverProps
	i.CreateNetworks = networks
	i.CreateRegistryEndpoint = registryEndpoint
	return i.CreateID, i.CreateErr
}

func (i *FakeServerService) Delete(id string) error {
	i.DeleteCalled = true
	return i.DeleteErr
}

func (i *FakeServerService) DeleteNetworkConfiguration(id string) error {
	i.DeleteNetworkConfigurationCalled = true
	return i.DeleteNetworkConfigurationErr
}

func (i *FakeServerService) DetachVolume(id string, diskID string) error {
	i.DetachVolumeCalled = true
	return i.DetachVolumeErr
}

func (i *FakeServerService) Find(id string) (*servers.Server, bool, error) {
	i.FindCalled = true
	return i.FindServer, i.FindFound, i.FindErr
}

func (i *FakeServerService) Reboot(id string) error {
	i.RebootCalled = true
	return i.RebootErr
}

func (i *FakeServerService) SetMetadata(id string, serverMetadata server.Metadata) error {
	i.SetMetadataCalled = true
	i.SetMetadataServerMetadata = serverMetadata
	return i.SetMetadataErr
}

func (i *FakeServerService) UpdateNetworkConfiguration(id string, networks server.Networks) error {
	i.UpdateNetworkConfigurationCalled = true
	return i.UpdateNetworkConfigurationErr
}
