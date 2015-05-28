package snapshot

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/snapshots"
)

func (s OpenStackSnapshotService) Find(id string) (Snapshot, bool, error) {
	s.logger.Debug(openstackSnapshotServiceLogTag, "Finding OpenStack Snapshot '%s'", id)
	snapshotItem, err := snapshots.Get(s.blockstorageService, id).Extract()
	if err != nil {
		return Snapshot{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Snapshot '%s'", id)
	}

	snapshot := Snapshot{
		ID:     snapshotItem.ID,
		Name:   snapshotItem.Name,
		Status: snapshotItem.Status,
	}
	return snapshot, true, nil
}
