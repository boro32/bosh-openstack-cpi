package volume

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func (v OpenStackVolumeService) Delete(id string) error {
	volume, found, err := v.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewDiskNotFoundError(id, false)
	}

	if volume.Status != openstackVolumeReadyStatus && volume.Status != openstackVolumeErrorStatus {
		return bosherr.WrapErrorf(err, "Cannot delete OpenStack Volume '%s', status is '%s'", id, volume.Status)
	}

	v.logger.Debug(openstackVolumeServiceLogTag, "Deleting OpenStack Volume '%s'", id)
	if err = volumes.Delete(v.blockstorageService, id).ExtractErr(); err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete OpenStack Volume '%s'", id)
	}

	return nil
}
