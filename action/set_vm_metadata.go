package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

type SetVMMetadata struct {
	serverService server.Service
}

func NewSetVMMetadata(
	serverService server.Service,
) SetVMMetadata {
	return SetVMMetadata{
		serverService: serverService,
	}
}

func (svm SetVMMetadata) Run(vmCID VMCID, vmMetadata VMMetadata) (interface{}, error) {
	if err := svm.serverService.SetMetadata(string(vmCID), server.Metadata(vmMetadata)); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Setting metadata for vm '%s'", vmCID)
	}

	return nil, nil
}
