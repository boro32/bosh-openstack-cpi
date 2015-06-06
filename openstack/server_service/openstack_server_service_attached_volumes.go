package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/rackspace/gophercloud/pagination"
)

func (s OpenStackServerService) AttachedVolumes(id string) (AttachedVolumes, error) {
	var volumes AttachedVolumes

	s.logger.Debug(openstackServerServiceLogTag, "Finding OpenStack Volumes attached to OpenStack Server '%s'", id)
	pager := volumeattach.List(s.computeService, id)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		volumeAttachments, err := volumeattach.ExtractVolumeAttachments(page)
		if err != nil {
			return false, err
		}

		for _, volumeAttachment := range volumeAttachments {
			volumes = append(volumes, volumeAttachment.VolumeID)
		}

		return true, nil
	})
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return volumes, bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return volumes, bosherr.WrapErrorf(err, "Failed to find OpenStack Volumes attached to OpenStack Server '%s'", id)
	}

	return volumes, nil
}
