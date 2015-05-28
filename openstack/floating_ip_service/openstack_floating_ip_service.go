package floatingip

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackFloatingIPServiceLogTag = "OpenStackFloatingIPService"

type OpenStackFloatingIPService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackFloatingIPService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackFloatingIPService {
	return OpenStackFloatingIPService{
		computeService: computeService,
		logger:         logger,
	}
}
