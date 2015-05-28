package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkServiceLogTag = "OpenStackNetworkService"

type OpenStackNetworkService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkService {
	return OpenStackNetworkService{
		computeService: computeService,
		logger:         logger,
	}
}
