package snapshot

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

func (s OpenStackSnapshotService) Create(volumeID string, description string) (string, error) {
	uuidStr, err := s.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random OpenStack Snapshot name")
	}

	if description == "" {
		description = openstackSnapshotDescription
	}

	createOpts := &snapshots.CreateOpts{
		Description: description,
		Force:       true,
		Name:        fmt.Sprintf("%s-%s", openstackSnapshotNamePrefix, uuidStr),
		VolumeID:    volumeID,
	}

	s.logger.Debug(openstackSnapshotServiceLogTag, "Creating OpenStack Snapshot with params: %#v", createOpts)
	snapshot, err := snapshots.Create(s.blockstorageService, createOpts).Extract()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create OpenStack Snapshot")
	}

	return snapshot.ID, nil
}

func (s OpenStackSnapshotService) cleanUp(id string) {
	if err := s.Delete(id); err != nil {
		s.logger.Debug(openstackSnapshotServiceLogTag, "Failed cleaning up OpenStack Snapshot '%s': %#v", id, err)
	}
}
