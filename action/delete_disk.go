package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"
)

type DeleteDisk struct {
	volumeService volume.Service
}

func NewDeleteDisk(
	volumeService volume.Service,
) DeleteDisk {
	return DeleteDisk{
		volumeService: volumeService,
	}
}

func (dd DeleteDisk) Run(diskCID DiskCID) (interface{}, error) {
	err := dd.volumeService.Delete(string(diskCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting disk '%s'", diskCID)
	}

	return nil, nil
}
