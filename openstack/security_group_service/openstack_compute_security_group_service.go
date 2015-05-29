package securitygroup

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackComputeSecurityGroupServiceLogTag = "OpenStackComputeSecurityGroupService"

type OpenStackComputeSecurityGroupService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackComputeSecurityGroupService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackComputeSecurityGroupService {
	return OpenStackComputeSecurityGroupService{
		computeService: computeService,
		logger:         logger,
	}
}
