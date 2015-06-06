package server

import (
	"encoding/json"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/bootfromvolume"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) Create(serverProps *Properties, networks Networks, registryEndpoint string) (string, error) {
	var createOpts servers.CreateOptsBuilder

	uuidStr, err := s.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random OpenStack Server name")
	}

	serverName := fmt.Sprintf("%s-%s", openstackServerNamePrefix, uuidStr)
	networksParams, err := s.createNetworksParams(networks)
	if err != nil {
		return "", err
	}
	securityGroupsParams, err := s.createSecurityGroupsParams(networks)
	if err != nil {
		return "", err
	}

	userdataParams, err := s.createUserdataParams(serverName, registryEndpoint, networks)
	if err != nil {
		return "", err
	}

	configDrive := true
	if s.disableConfigDrive {
		configDrive = false
	}

	createOpts = &servers.CreateOpts{
		Name:             serverName,
		ImageRef:         serverProps.ImageID,
		FlavorRef:        serverProps.FlavorID,
		AvailabilityZone: serverProps.AvailabilityZone,
		Networks:         networksParams,
		SecurityGroups:   securityGroupsParams,
		UserData:         userdataParams,
		ConfigDrive:      configDrive,
		AdminPass:        "TODO",
	}
	createOpts = s.addBootFromVolumeParams(createOpts, serverProps.ImageID, serverProps.RootDiskSizeGb)
	createOpts = s.addKeyPairParams(createOpts, serverProps.KeyPair)
	createOpts = s.addSchedulerHintsParams(createOpts, serverProps.SchedulerHints)

	serverOpts, _ := createOpts.ToServerCreateMap()
	s.logger.Debug(openstackServerServiceLogTag, "Creating OpenStack Server with params: %#v", serverOpts)
	server, err := servers.Create(s.computeService, createOpts).Extract()
	if err != nil {
		s.logger.Debug(openstackServerServiceLogTag, "Failed to create OpenStack Server: %#v", err)
		return "", api.NewVMCreationFailedError(true)
	}

	return server.ID, nil
}

func (s OpenStackServerService) CleanUp(id string) {
	if err := s.Delete(id); err != nil {
		s.logger.Debug(openstackServerServiceLogTag, "Failed cleaning up OpenStack Server '%s': %#v", id, err)
	}
}

func (s OpenStackServerService) addBootFromVolumeParams(
	createOpts servers.CreateOptsBuilder,
	imageID string,
	rootDiskSizeGb int,
) *bootfromvolume.CreateOptsExt {
	blockDevice := bootfromvolume.BlockDevice{
		BootIndex:           0,
		DeleteOnTermination: true,
		DestinationType:     "volume",
		SourceType:          bootfromvolume.SourceType("image"),
		UUID:                imageID,
	}
	if rootDiskSizeGb > 0 {
		blockDevice.VolumeSize = rootDiskSizeGb
	}

	return &bootfromvolume.CreateOptsExt{
		createOpts,
		[]bootfromvolume.BlockDevice{blockDevice},
	}
}

func (s OpenStackServerService) addKeyPairParams(
	createOpts servers.CreateOptsBuilder,
	keypair string,
) *keypairs.CreateOptsExt {
	return &keypairs.CreateOptsExt{
		createOpts,
		keypair,
	}
}

func (s OpenStackServerService) addSchedulerHintsParams(
	createOpts servers.CreateOptsBuilder,
	schedulerHintsProperties SchedulerHintsProperties,
) *schedulerhints.CreateOptsExt {
	var schedulerHints schedulerhints.SchedulerHints

	if schedulerHintsProperties.Group != "" {
		schedulerHints.Group = schedulerHintsProperties.Group
	}

	if schedulerHintsProperties.DifferentHost != nil {
		schedulerHints.DifferentHost = schedulerHintsProperties.DifferentHost
	}

	if schedulerHintsProperties.SameHost != nil {
		schedulerHints.SameHost = schedulerHintsProperties.SameHost
	}

	if schedulerHintsProperties.Query != nil {
		schedulerHints.Query = schedulerHintsProperties.Query
	}

	if schedulerHintsProperties.TargetCell != "" {
		schedulerHints.TargetCell = schedulerHintsProperties.TargetCell
	}

	if schedulerHintsProperties.BuildNearHostIP != "" {
		schedulerHints.BuildNearHostIP = schedulerHintsProperties.BuildNearHostIP
	}

	return &schedulerhints.CreateOptsExt{
		createOpts,
		schedulerHints,
	}
}

func (s OpenStackServerService) createNetworksParams(networks Networks) ([]servers.Network, error) {
	var networksParams []servers.Network

	if s.disableNeutron {
		return networksParams, nil
	}

	for _, network := range networks {
		if networkName := network.NetworkName(); networkName != "" {
			net, found, err := s.networkService.FindByName(networkName)
			if err != nil {
				return networksParams, err
			}
			if !found {
				return networksParams, bosherr.Errorf("OpenStack Network '%s' not found", networkName)
			}
			serverNetwork := servers.Network{UUID: net.ID}

			if ipAddress := network.IPAddress(); ipAddress != "" {
				serverNetwork.FixedIP = ipAddress
			}

			networksParams = append(networksParams, serverNetwork)
		}
	}

	return networksParams, nil
}

func (s OpenStackServerService) createSecurityGroupsParams(networks Networks) ([]string, error) {
	var securityGroupsParams []string

	for _, securityGroup := range networks.SecurityGroupsList() {
		_, found, err := s.securityGroupService.FindByName(securityGroup)
		if err != nil {
			return securityGroupsParams, err
		}
		if !found {
			return securityGroupsParams, bosherr.Errorf("OpenStack Security Group '%s' not found", securityGroup)
		}

		securityGroupsParams = append(securityGroupsParams, securityGroup)
	}

	if len(securityGroupsParams) > 0 {
		return securityGroupsParams, nil
	}

	return s.defaultSecurityGroups, nil
}

func (s OpenStackServerService) createUserdataParams(name string, registryEndpoint string, networks Networks) ([]byte, error) {
	openstackServerName := OpenStackUserDataServerName{Name: name}
	openstackRegistryEndpoint := OpenStackUserDataRegistryEndpoint{Endpoint: registryEndpoint}
	openstackUserData := OpenStackUserData{Server: openstackServerName, Registry: openstackRegistryEndpoint}

	if networkDNS := networks.DNSList(); len(networkDNS) > 0 {
		openstackUserData.DNS = OpenStackUserDataDNSItems{NameServer: networkDNS}
	}

	userData, err := json.Marshal(openstackUserData)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Marshalling user data")
	}

	return userData, nil
}
