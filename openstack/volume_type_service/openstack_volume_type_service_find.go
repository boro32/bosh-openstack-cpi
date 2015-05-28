package volumetype

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumetypes"
)

func (vt OpenStackVolumeTypeService) Find(id string) (VolumeType, bool, error) {
	vt.logger.Debug(openstackVolumeTypeServiceLogTag, "Finding OpenStack Volume Type '%s'", id)
	volumeTypeItem, err := volumetypes.Get(vt.blockstorageService, id).Extract()
	if err != nil {
		return VolumeType{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Volume Type '%s'", id)
	}

	volumeType := VolumeType{
		ID:   volumeTypeItem.ID,
		Name: volumeTypeItem.Name,
	}
	return volumeType, true, nil
}
