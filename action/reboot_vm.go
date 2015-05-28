package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

type RebootVM struct {
	serverService server.Service
}

func NewRebootVM(
	serverService server.Service,
) RebootVM {
	return RebootVM{
		serverService: serverService,
	}
}

func (rv RebootVM) Run(vmCID VMCID) (interface{}, error) {
	err := rv.serverService.Reboot(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Rebooting vm '%s'", vmCID)
	}

	return nil, nil
}
