package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkNetworkServiceLogTag = "OpenStackNetworkNetworkService"

type OpenStackNetworkNetworkService struct {
	networkService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkNetworkService(
	networkService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkNetworkService {
	return OpenStackNetworkNetworkService{
		networkService: networkService,
		logger:         logger,
	}
}
