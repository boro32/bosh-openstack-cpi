package extension

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkExtensionServiceLogTag = "OpenStackNetworkExtensionService"

type OpenStackNetworkExtensionService struct {
	networkService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkExtensionService(
	networkService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkExtensionService {
	return OpenStackNetworkExtensionService{
		networkService: networkService,
		logger:         logger,
	}
}
