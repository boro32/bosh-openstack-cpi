package volume

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/blockstorage/v1/volumes"
)

func (v OpenStackVolumeService) Create(size int, volumeType string, availabilityZone string) (string, error) {
	uuidStr, err := v.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random OpenStack Volume name")
	}

	createOpts := &volumes.CreateOpts{
		Availability: availabilityZone,
		Description:  openstackVolumeDescription,
		Name:         fmt.Sprintf("%s-%s", openstackVolumeNamePrefix, uuidStr),
		Size:         size,
		VolumeType:   volumeType,
	}

	v.logger.Debug(openstackVolumeServiceLogTag, "Creating OpenStack Volume with params: %#v", createOpts)
	volume, err := volumes.Create(v.blockstorageService, createOpts).Extract()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create OpenStack Volume")
	}

	return volume.ID, nil
}

func (v OpenStackVolumeService) cleanUp(id string) {
	if err := v.Delete(id); err != nil {
		v.logger.Debug(openstackVolumeServiceLogTag, "Failed cleaning up OpenStack Volume '%s': %#v", id, err)
	}
}
