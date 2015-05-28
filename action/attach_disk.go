package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"

	"github.com/frodenas/bosh-registry/client"
)

type AttachDisk struct {
	serverService  server.Service
	registryClient registry.Client
}

func NewAttachDisk(
	serverService server.Service,
	registryClient registry.Client,
) AttachDisk {
	return AttachDisk{
		serverService:  serverService,
		registryClient: registryClient,
	}
}

func (ad AttachDisk) Run(vmCID VMCID, diskCID DiskCID) (interface{}, error) {
	// Atach the Disk to the VM
	deviceName, devicePath, err := ad.serverService.AttachVolume(string(vmCID), string(diskCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Read VM agent settings
	agentSettings, err := ad.registryClient.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	// Update VM agent settings
	newAgentSettings := agentSettings.AttachPersistentDisk(string(diskCID), deviceName, devicePath)
	err = ad.registryClient.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Attaching disk '%s' to vm '%s'", diskCID, vmCID)
	}

	return nil, nil
}
