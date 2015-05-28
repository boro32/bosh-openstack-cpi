package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

type HasVM struct {
	serverService server.Service
}

func NewHasVM(
	serverService server.Service,
) HasVM {
	return HasVM{
		serverService: serverService,
	}
}

func (hv HasVM) Run(vmCID VMCID) (bool, error) {
	_, found, err := hv.serverService.Find(string(vmCID))
	if err != nil {
		return false, bosherr.WrapErrorf(err, "Finding vm '%s'", vmCID)
	}

	return found, nil
}
