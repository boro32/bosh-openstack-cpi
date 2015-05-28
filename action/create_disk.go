package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_type_service"
	"github.com/frodenas/bosh-openstack-cpi/util"
)

type CreateDisk struct {
	volumeService                volume.Service
	volumeTypeService            volumetype.Service
	serverService                server.Service
	ignoreServerAvailabilityZone bool
}

func NewCreateDisk(
	volumeService volume.Service,
	volumeTypeService volumetype.Service,
	serverService server.Service,
	ignoreServerAvailabilityZone bool,
) CreateDisk {
	return CreateDisk{
		volumeService:                volumeService,
		volumeTypeService:            volumeTypeService,
		serverService:                serverService,
		ignoreServerAvailabilityZone: ignoreServerAvailabilityZone,
	}
}

func (cd CreateDisk) Run(size int, cloudProps DiskCloudProperties, vmCID VMCID) (DiskCID, error) {
	availabilityZone := cloudProps.AvailabilityZone

	if !cd.ignoreServerAvailabilityZone {
		// Find the VM (if provided) so we can create the volume in the same availability_zone
		if vmCID != "" {
			_, found, err := cd.serverService.Find(string(vmCID))
			if err != nil {
				return "", bosherr.WrapError(err, "Creating disk")
			}
			if !found {
				return "", api.NewVMNotFoundError(string(vmCID))
			}

			// TODO: gophercloud does not return the availability zone
			// TODO: check cloud properties and server availability zones
			//availabilityZone = vm.AvailabilityZone
		}
	}

	// Find the Volume Type (if provided)
	if cloudProps.VolumeType != "" {
		_, found, err := cd.volumeTypeService.FindByName(cloudProps.VolumeType)
		if err != nil {
			return "", bosherr.WrapError(err, "Creating disk")
		}
		if !found {
			return "", bosherr.WrapErrorf(err, "Creating disk: Volume Type '%s' does not exists", cloudProps.VolumeType)
		}
	}

	// Create the Disk
	diskID, err := cd.volumeService.Create(util.ConvertMib2Gib(size), cloudProps.VolumeType, availabilityZone)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating disk")
	}

	return DiskCID(diskID), nil
}
