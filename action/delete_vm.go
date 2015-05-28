package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"

	"github.com/frodenas/bosh-registry/client"
)

type DeleteVM struct {
	serverService  server.Service
	registryClient registry.Client
}

func NewDeleteVM(
	serverService server.Service,
	registryClient registry.Client,
) DeleteVM {
	return DeleteVM{
		serverService:  serverService,
		registryClient: registryClient,
	}
}

func (dv DeleteVM) Run(vmCID VMCID) (interface{}, error) {
	// Delete VM networks
	err := dv.serverService.DeleteNetworkConfiguration(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM
	err = dv.serverService.Delete(string(vmCID))
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	// Delete the VM agent settings
	err = dv.registryClient.Delete(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting vm '%s'", vmCID)
	}

	return nil, nil
}
