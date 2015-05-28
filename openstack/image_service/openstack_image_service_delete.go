package image

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
)

func (i OpenStackImageService) Delete(id string) error {
	image, found, err := i.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "OpenStack Image '%s' does not exists", id)
	}

	if image.Status != openstackImageReadyStatus && image.Status != openstackImageErrorStatus {
		return bosherr.WrapErrorf(err, "Cannot delete OpenStack Image '%s', status is '%s'", id, image.Status)
	}

	i.logger.Debug(openstackImageServiceLogTag, "Deleting OpenStack Image '%s'", id)
	if err = images.Delete(i.computeService, id).ExtractErr(); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete OpenStack Image '%s'", id)
	}

	return nil
}
