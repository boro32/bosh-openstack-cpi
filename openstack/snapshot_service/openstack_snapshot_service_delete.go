package snapshot

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

func (s OpenStackSnapshotService) Delete(id string) error {
	snapshot, found, err := s.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.WrapErrorf(err, "OpenStack Snapshot '%s' does not exists", id)
	}

	if snapshot.Status != openstackSnapshotReadyStatus && snapshot.Status != openstackSnapshotErrorStatus {
		return bosherr.WrapErrorf(err, "Cannot delete OpenStack Snapshot '%s', status is '%s'", id, snapshot.Status)
	}

	s.logger.Debug(openstackSnapshotServiceLogTag, "Deleting OpenStack Snapshot '%s'", id)
	err = snapshots.Delete(s.blockstorageService, id).ExtractErr()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete OpenStack Snapshot '%s'", id)
	}

	return nil
}
