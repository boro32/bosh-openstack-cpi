package flavor

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackFlavorServiceLogTag = "OpenStackFlavorService"

type OpenStackFlavorService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackFlavorService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackFlavorService {
	return OpenStackFlavorService{
		computeService: computeService,
		logger:         logger,
	}
}
