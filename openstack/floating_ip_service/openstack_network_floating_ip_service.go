package floatingip

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkFloatingIPServiceLogTag = "OpenStackNetworkFloatingIPService"

type OpenStackNetworkFloatingIPService struct {
	networkService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkFloatingIPService(
	networkService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkFloatingIPService {
	return OpenStackNetworkFloatingIPService{
		networkService: networkService,
		logger:         logger,
	}
}
