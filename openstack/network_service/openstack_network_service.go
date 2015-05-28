package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkServiceLogTag = "OpenStackNetworkService"

type OpenStackNetworkService struct {
	networkService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkService(
	networkService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkService {
	return OpenStackNetworkService{
		networkService: networkService,
		logger:         logger,
	}
}
