package floatingip

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkFloatingIPServiceLogTag = "OpenStackNetworkFloatingIPService"

type OpenStackNetworkFloatingIPService struct {
	networkService *gophercloud.ServiceClient
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkFloatingIPService(
	networkService *gophercloud.ServiceClient,
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkFloatingIPService {
	return OpenStackNetworkFloatingIPService{
		networkService: networkService,
		computeService: computeService,
		logger:         logger,
	}
}
