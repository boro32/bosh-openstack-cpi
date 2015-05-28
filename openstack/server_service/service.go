package server

import (
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

type Service interface {
	AddNetworkConfiguration(id string, networks Networks) error
	AttachVolume(id string, volumeID string) (string, string, error)
	AttachedVolumes(id string) (AttachedVolumes, error)
	CleanUp(id string)
	Create(vmProps *Properties, networks Networks, registryEndpoint string) (string, error)
	Delete(id string) error
	DeleteNetworkConfiguration(id string) error
	DetachVolume(id string, volumeID string) error
	Find(id string) (*servers.Server, bool, error)
	Reboot(id string) error
	SetMetadata(id string, serverMetadata Metadata) error
	UpdateNetworkConfiguration(id string, networks Networks) error
}

type AttachedVolumes []string

type Metadata map[string]interface{}

type Properties struct {
	ImageID            string
	FlavorID           string
	AvailabilityZone   string
	KeyPair            string
	DisableConfigDrive bool
	SchedulerHints     SchedulerHintsProperties
}

type SchedulerHintsProperties struct {
	Group           string
	DifferentHost   []string
	SameHost        []string
	Query           []interface{}
	TargetCell      string
	BuildNearHostIP string
}
