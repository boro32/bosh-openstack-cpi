package image

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/rackspace/gophercloud"
)

const openstackImageServiceLogTag = "OpenStackImageService"
const openstackImageNamePrefix = "stemcell"
const openstackImageDescription = "Image managed by BOSH"
const openstackImageReadyStatus = "READY"

type OpenStackImageService struct {
	computeService *gophercloud.ServiceClient
	uuidGen        boshuuid.Generator
	logger         boshlog.Logger
}

func NewOpenStackImageService(
	computeService *gophercloud.ServiceClient,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) OpenStackImageService {
	return OpenStackImageService{
		computeService: computeService,
		uuidGen:        uuidGen,
		logger:         logger,
	}
}
