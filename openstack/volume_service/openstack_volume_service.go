package volume

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/rackspace/gophercloud"
)

const openstackVolumeServiceLogTag = "OpenStackVolumeService"
const openstackVolumeNamePrefix = "volume"
const openstackVolumeDescription = "Volume managed by BOSH"
const openstackVolumeReadyStatus = "available"

type OpenStackVolumeService struct {
	blockstorageService *gophercloud.ServiceClient
	uuidGen             boshuuid.Generator
	logger              boshlog.Logger
}

func NewOpenStackVolumeService(
	blockstorageService *gophercloud.ServiceClient,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) OpenStackVolumeService {
	return OpenStackVolumeService{
		blockstorageService: blockstorageService,
		uuidGen:             uuidGen,
		logger:              logger,
	}
}
