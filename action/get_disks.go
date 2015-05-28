package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

type GetDisks struct {
	serverService server.Service
}

func NewGetDisks(
	serverService server.Service,
) GetDisks {
	return GetDisks{
		serverService: serverService,
	}
}

func (gd GetDisks) Run(vmCID VMCID) (disks []string, err error) {
	disks, err = gd.serverService.AttachedVolumes(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Finding disks for vm '%s'", vmCID)
	}

	return disks, nil
}
