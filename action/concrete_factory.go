package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshuuid "github.com/cloudfoundry/bosh-utils/uuid"

	"github.com/frodenas/bosh-openstack-cpi/openstack/client"
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/floating_ip_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/network_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/security_group_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/snapshot_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_type_service"

	"github.com/frodenas/bosh-registry/client"
)

type ConcreteFactory struct {
	availableActions map[string]Action
}

func NewConcreteFactory(
	openstackClient client.OpenStackClient,
	uuidGen boshuuid.Generator,
	options ConcreteFactoryOptions,
	logger boshlog.Logger,
) ConcreteFactory {
	flavorService := flavor.NewOpenStackFlavorService(
		openstackClient.ComputeService(),
		logger,
	)

	var floatingIPService floatingip.Service
	if openstackClient.DisableNeutron() {
		floatingIPService = floatingip.NewOpenStackComputeFloatingIPService(
			openstackClient.ComputeService(),
			logger,
		)
	} else {
		floatingIPService = floatingip.NewOpenStackNetworkFloatingIPService(
			openstackClient.NetworkService(),
			logger,
		)
	}

	imageService := image.NewOpenStackImageService(
		openstackClient.ComputeService(),
		uuidGen,
		logger,
	)

	keypairService := keypair.NewOpenStackKeyPairService(
		openstackClient.ComputeService(),
		logger,
	)

	var networkService network.Service
	if openstackClient.DisableNeutron() {
		networkService = network.NewOpenStackComputeNetworkService(
			openstackClient.ComputeService(),
			logger,
		)
	} else {
		networkService = network.NewOpenStackNetworkNetworkService(
			openstackClient.NetworkService(),
			logger,
		)
	}

	registryClient := registry.NewHTTPClient(
		options.Registry,
		logger,
	)

	var securityGroupService securitygroup.Service
	if openstackClient.DisableNeutron() {
		securityGroupService = securitygroup.NewOpenStackComputeSecurityGroupService(
			openstackClient.ComputeService(),
			logger,
		)
	} else {
		securityGroupService = securitygroup.NewOpenStackNetworkSecurityGroupService(
			openstackClient.NetworkService(),
			logger,
		)
	}

	serverService := server.NewOpenStackServerService(
		openstackClient.ComputeService(),
		floatingIPService,
		networkService,
		securityGroupService,
		openstackClient.DefaultSecurityGroups(),
		openstackClient.DisableConfigDrive(),
		openstackClient.DisableNeutron(),
		uuidGen,
		logger,
	)

	snapshotService := snapshot.NewOpenStackSnapshotService(
		openstackClient.BlockStorageService(),
		uuidGen,
		logger,
	)

	volumeService := volume.NewOpenStackVolumeService(
		openstackClient.BlockStorageService(),
		uuidGen,
		logger,
	)

	volumeTypeService := volumetype.NewOpenStackVolumeTypeService(
		openstackClient.BlockStorageService(),
		logger,
	)

	return ConcreteFactory{
		availableActions: map[string]Action{
			// Disk management
			"create_disk": NewCreateDisk(
				volumeService,
				volumeTypeService,
				serverService,
				openstackClient.IgnoreServerAvailabilityZone(),
			),
			"delete_disk": NewDeleteDisk(volumeService),

			// Snapshot management
			"snapshot_disk":   NewSnapshotDisk(snapshotService),
			"delete_snapshot": NewDeleteSnapshot(snapshotService),

			// Stemcell management
			"create_stemcell": NewCreateStemcell(imageService),
			"delete_stemcell": NewDeleteStemcell(imageService),

			// VM management
			"create_vm": NewCreateVM(
				serverService,
				flavorService,
				imageService,
				keypairService,
				volumeService,
				registryClient,
				options.Registry,
				options.Agent,
				openstackClient.DefaultKeyPair(),
			),
			"configure_networks": NewConfigureNetworks(serverService, registryClient),
			"delete_vm":          NewDeleteVM(serverService, registryClient),
			"reboot_vm":          NewRebootVM(serverService),
			"set_vm_metadata":    NewSetVMMetadata(serverService),
			"has_vm":             NewHasVM(serverService),
			"attach_disk":        NewAttachDisk(serverService, registryClient),
			"detach_disk":        NewDetachDisk(serverService, registryClient),
			"get_disks":          NewGetDisks(serverService),

			// Others:
			"ping": NewPing(),

			// Not implemented:
			// current_vm_id
		},
	}
}

func (f ConcreteFactory) Create(method string) (Action, error) {
	action, found := f.availableActions[method]
	if !found {
		return nil, bosherr.Errorf("Could not create action with method %s", method)
	}

	return action, nil
}
