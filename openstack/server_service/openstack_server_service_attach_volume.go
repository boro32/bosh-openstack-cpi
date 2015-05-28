package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func (i OpenStackServerService) AttachVolume(id string, volumeID string) (string, string, error) {
	var deviceName, devicePath string

	createOps := &volumeattach.CreateOpts{
		VolumeID: volumeID,
	}

	i.logger.Debug(openstackServerServiceLogTag, "Attaching OpenStack Volume '%s' to OpenStack Server '%s'", volumeID, id)
	volumeAttachment, err := volumeattach.Create(i.computeService, id, createOps).Extract()
	if err != nil {
		return deviceName, devicePath, bosherr.WrapErrorf(err, "Failed to attach OpenStack Volume '%s' to OpenStack Server '%s'", volumeID, id)
	}

	return volumeAttachment.Device, devicePath, nil
}
