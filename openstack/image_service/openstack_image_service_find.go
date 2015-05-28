package image

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/images"
)

func (i OpenStackImageService) Find(id string) (Image, bool, error) {
	i.logger.Debug(openstackImageServiceLogTag, "Finding OpenStack Image '%s'", id)
	imageItem, err := images.Get(i.computeService, id).Extract()
	if err != nil {
		return Image{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Image '%s'", id)
	}

	image := Image{
		ID:     imageItem.ID,
		Name:   imageItem.Name,
		Status: imageItem.Status,
	}
	return image, true, nil
}
