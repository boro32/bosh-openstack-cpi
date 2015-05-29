package network

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackComputeNetworkServiceLogTag = "OpenStackComputeNetworkService"

type OpenStackComputeNetworkService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackComputeNetworkService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackComputeNetworkService {
	return OpenStackComputeNetworkService{
		computeService: computeService,
		logger:         logger,
	}
}
