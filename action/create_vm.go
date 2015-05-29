package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/flavor_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
	"github.com/frodenas/bosh-openstack-cpi/openstack/volume_service"

	"github.com/frodenas/bosh-registry/client"
)

type CreateVM struct {
	serverService   server.Service
	flavorService   flavor.Service
	imageService    image.Service
	keypairService  keypair.Service
	volumeService   volume.Service
	registryClient  registry.Client
	registryOptions registry.ClientOptions
	agentOptions    registry.AgentOptions
	defaultKeyPair  string
}

func NewCreateVM(
	serverService server.Service,
	flavorService flavor.Service,
	imageService image.Service,
	keypairService keypair.Service,
	volumeService volume.Service,
	registryClient registry.Client,
	registryOptions registry.ClientOptions,
	agentOptions registry.AgentOptions,
	defaultKeyPair string,
) CreateVM {
	return CreateVM{
		serverService:   serverService,
		flavorService:   flavorService,
		imageService:    imageService,
		keypairService:  keypairService,
		volumeService:   volumeService,
		registryClient:  registryClient,
		registryOptions: registryOptions,
		agentOptions:    agentOptions,
		defaultKeyPair:  defaultKeyPair,
	}
}

func (cv CreateVM) Run(agentID string, stemcellCID StemcellCID, cloudProps VMCloudProperties, networks Networks, disks []DiskCID, env Environment) (VMCID, error) {
	// Find all affinity zones
	availabilityZones := make(map[string]struct{})
	if cloudProps.AvailabilityZone != "" {
		availabilityZones[cloudProps.AvailabilityZone] = struct{}{}
	}
	for _, diskCID := range disks {
		volume, found, err := cv.volumeService.Find(string(diskCID))
		if err != nil {
			return "", bosherr.WrapError(err, "Creating vm")
		}
		if !found {
			return "", api.NewDiskNotFoundError(string(diskCID), false)
		}
		availabilityZones[volume.AvailabilityZone] = struct{}{}
	}
	if len(availabilityZones) > 1 {
		return "", bosherr.Errorf("Creating vm: can't use multiple availability zones: '%v'", availabilityZones)
	}

	// Determine availability zone
	availabilityZone := ""
	for k := range availabilityZones {
		availabilityZone = k
		break
	}

	// Find image
	image, found, err := cv.imageService.Find(string(stemcellCID))
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: Stemcell '%s' does not exists", stemcellCID)
	}

	// Find flavor
	if cloudProps.Flavor == "" {
		return "", bosherr.WrapError(err, "Creating vm: 'flavor' must be provided")
	}
	flavor, found, err := cv.flavorService.FindByName(cloudProps.Flavor)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: Flavor '%s' does not exists", cloudProps.Flavor)
	}

	// Find keypair
	keyPair := cv.defaultKeyPair
	if cloudProps.KeyPair != "" {
		keyPair = cloudProps.KeyPair
	}
	_, found, err = cv.keypairService.Find(keyPair)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating vm")
	}
	if !found {
		return "", bosherr.WrapErrorf(err, "Creating vm: KeyPair '%s' does not exists", keyPair)
	}

	// Parse VM networks
	serverNetworks := networks.AsServerServiceNetworks()
	if err = serverNetworks.Validate(); err != nil {
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Parse VM properties
	serverProps := &server.Properties{
		ImageID:          image.ID,
		FlavorID:         flavor.ID,
		AvailabilityZone: availabilityZone,
		KeyPair:          keyPair,
		SchedulerHints: server.SchedulerHintsProperties{
			Group:           cloudProps.SchedulerHints.Group,
			DifferentHost:   cloudProps.SchedulerHints.DifferentHost,
			SameHost:        cloudProps.SchedulerHints.SameHost,
			Query:           cloudProps.SchedulerHints.Query,
			TargetCell:      cloudProps.SchedulerHints.TargetCell,
			BuildNearHostIP: cloudProps.SchedulerHints.BuildNearHostIP,
		},
	}

	// Create VM
	vmID, err := cv.serverService.Create(serverProps, serverNetworks, cv.registryOptions.Endpoint())
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// If any of the below code fails, we must delete the created vm
	defer func() {
		if err != nil {
			cv.serverService.CleanUp(vmID)
		}
	}()

	// Configure VM networks
	if err = cv.serverService.AddNetworkConfiguration(vmID, serverNetworks); err != nil {
		if _, ok := err.(api.CloudError); ok {
			return "", err
		}
		return "", bosherr.WrapError(err, "Creating VM")
	}

	// Create VM settings
	agentNetworks := networks.AsRegistryNetworks()
	agentSettings := registry.NewAgentSettings(agentID, vmID, agentNetworks, registry.EnvSettings(env), cv.agentOptions)
	if err = cv.registryClient.Update(vmID, agentSettings); err != nil {
		return "", bosherr.WrapErrorf(err, "Creating VM")
	}

	return VMCID(vmID), nil
}
