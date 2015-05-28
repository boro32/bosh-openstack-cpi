package volumetype

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackVolumeTypeServiceLogTag = "OpenStackVolumeTypeService"

type OpenStackVolumeTypeService struct {
	blockstorageService *gophercloud.ServiceClient
	logger              boshlog.Logger
}

func NewOpenStackVolumeTypeService(
	blockstorageService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackVolumeTypeService {
	return OpenStackVolumeTypeService{
		blockstorageService: blockstorageService,
		logger:              logger,
	}
}
