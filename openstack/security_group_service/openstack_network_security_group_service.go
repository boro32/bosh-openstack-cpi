package securitygroup

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackNetworkSecurityGroupServiceLogTag = "OpenStackNetworkSecurityGroupService"

type OpenStackNetworkSecurityGroupService struct {
	networkService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackNetworkSecurityGroupService(
	networkService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackNetworkSecurityGroupService {
	return OpenStackNetworkSecurityGroupService{
		networkService: networkService,
		logger:         logger,
	}
}
