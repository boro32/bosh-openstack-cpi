package action

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/openstack/snapshot_service"
)

type SnapshotDisk struct {
	snapshotService snapshot.Service
}

func NewSnapshotDisk(
	snapshotService snapshot.Service,
) SnapshotDisk {
	return SnapshotDisk{
		snapshotService: snapshotService,
	}
}

func (sd SnapshotDisk) Run(diskCID DiskCID, metadata SnapshotMetadata) (SnapshotCID, error) {
	var description string
	if metadata.Deployment != "" && metadata.Job != "" && metadata.Index != "" {
		description = fmt.Sprintf("%s/%s/%s", metadata.Deployment, metadata.Job, metadata.Index)
	}

	snapshotID, err := sd.snapshotService.Create(string(diskCID), description)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Creating snapshot of disk '%s'", diskCID)
	}

	return SnapshotCID(snapshotID), nil
}
