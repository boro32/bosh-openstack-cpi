package extension

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackComputeExtensionServiceLogTag = "OpenStackComputeExtensionService"

type OpenStackComputeExtensionService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackComputeExtensionService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackComputeExtensionService {
	return OpenStackComputeExtensionService{
		computeService: computeService,
		logger:         logger,
	}
}
