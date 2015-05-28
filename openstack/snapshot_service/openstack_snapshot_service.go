package snapshot

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/rackspace/gophercloud"
)

const openstackSnapshotServiceLogTag = "OpenStackSnapshotService"
const openstackSnapshotNamePrefix = "snapshot"
const openstackSnapshotDescription = "Snapshot managed by BOSH"
const openstackSnapshotReadyStatus = "available"
const openstackSnapshotErrorStatus = "error"

type OpenStackSnapshotService struct {
	blockstorageService *gophercloud.ServiceClient
	uuidGen             boshuuid.Generator
	logger              boshlog.Logger
}

func NewOpenStackSnapshotService(
	blockstorageService *gophercloud.ServiceClient,
	uuidGen boshuuid.Generator,
	logger boshlog.Logger,
) OpenStackSnapshotService {
	return OpenStackSnapshotService{
		blockstorageService: blockstorageService,
		uuidGen:             uuidGen,
		logger:              logger,
	}
}
