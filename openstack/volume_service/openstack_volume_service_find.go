package volume

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func (v OpenStackVolumeService) Find(id string) (Volume, bool, error) {
	v.logger.Debug(openstackVolumeServiceLogTag, "Finding OpenStack Volume '%s'", id)
	volumeItem, err := volumes.Get(v.blockstorageService, id).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return Volume{}, false, nil
		}

		return Volume{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Volume '%s'", id)
	}

	volume := Volume{
		ID:               volumeItem.ID,
		Name:             volumeItem.Name,
		Status:           volumeItem.Status,
		AvailabilityZone: volumeItem.AvailabilityZone,
	}
	return volume, true, nil
}
