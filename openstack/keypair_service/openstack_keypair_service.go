package keypair

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"

	"github.com/rackspace/gophercloud"
)

const openstackKeyPairServiceLogTag = "OpenStackKeyPairService"

type OpenStackKeyPairService struct {
	computeService *gophercloud.ServiceClient
	logger         boshlog.Logger
}

func NewOpenStackKeyPairService(
	computeService *gophercloud.ServiceClient,
	logger boshlog.Logger,
) OpenStackKeyPairService {
	return OpenStackKeyPairService{
		computeService: computeService,
		logger:         logger,
	}
}
