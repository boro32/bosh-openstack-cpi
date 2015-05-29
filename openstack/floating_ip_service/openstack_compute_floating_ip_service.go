package floatingip

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackComputeFloatingIPServiceLogTag = "OpenStackComputeFloatingIPService"

type OpenStackComputeFloatingIPService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackComputeFloatingIPService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackComputeFloatingIPService {
	return OpenStackComputeFloatingIPService{
		computeService: computeService,
		logger:         logger,
	}
}
