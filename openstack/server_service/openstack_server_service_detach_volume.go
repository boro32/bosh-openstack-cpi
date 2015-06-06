package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/rackspace/gophercloud/pagination"
)

func (s OpenStackServerService) DetachVolume(id string, volumeID string) error {
	// Look up for the device name
	var volumeAttachmentID string
	pager := volumeattach.List(s.computeService, id)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		volumeAttachments, err := volumeattach.ExtractVolumeAttachments(page)
		if err != nil {
			return false, err
		}

		for _, volumeAttachment := range volumeAttachments {
			if volumeID == volumeAttachment.VolumeID {
				volumeAttachmentID = volumeAttachment.ID
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to find OpenStack Volumes attached to OpenStack Server '%s'", id)
	}

	if volumeAttachmentID == "" {
		return api.NewDiskNotAttachedError(id, volumeID, false)
	}

	// Detach the volume
	s.logger.Debug(openstackServerServiceLogTag, "Detaching OpenStack Volume '%s' from OpenStack Server '%s'", volumeID, id)
	if err = volumeattach.Delete(s.computeService, id, volumeAttachmentID).ExtractErr(); err != nil {
		return bosherr.WrapErrorf(err, "Failed to detach OpenStack Volume '%s' from OpenStack Server '%s'", volumeID, id)
	}

	return nil
}
