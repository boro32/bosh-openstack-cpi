package volumetype

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumetypes"
	"github.com/rackspace/gophercloud/pagination"
)

func (vt OpenStackVolumeTypeService) FindByName(name string) (VolumeType, bool, error) {
	var volumeType VolumeType

	vt.logger.Debug(openstackVolumeTypeServiceLogTag, "Finding OpenStack Volume Type '%s'", name)
	pager := volumetypes.List(vt.blockstorageService)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		volumeTypeList, err := volumetypes.ExtractVolumeTypes(page)
		if err != nil {
			return false, err
		}

		for _, volumeTypeItem := range volumeTypeList {
			if volumeTypeItem.Name == name {
				volumeType = VolumeType{
					ID:   volumeTypeItem.ID,
					Name: volumeTypeItem.Name,
				}
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return volumeType, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Volume Type '%s'", name)
	}

	if volumeType.ID != "" {
		return volumeType, true, nil
	}

	return volumeType, false, nil
}
