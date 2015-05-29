package server

import (
	"encoding/json"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/keypairs"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/schedulerhints"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (i OpenStackServerService) Create(serverProps *Properties, networks Networks, registryEndpoint string) (string, error) {
	var createOpts servers.CreateOptsBuilder

	uuidStr, err := i.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random OpenStack Server name")
	}

	serverName := fmt.Sprintf("%s-%s", openstackServerNamePrefix, uuidStr)
	networksParams, err := i.createNetworksParams(networks)
	if err != nil {
		return "", err
	}
	securityGroupsParams, err := i.createSecurityGroupsParams(networks)
	if err != nil {
		return "", err
	}

	userdataParams, err := i.createUserdataParams(serverName, registryEndpoint, networks)
	if err != nil {
		return "", err
	}

	configDrive := true
	if i.disableConfigDrive {
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
	}

	createOpts = &keypairs.CreateOptsExt{
		createOpts,
		serverProps.KeyPair,
	}

	createOpts = &schedulerhints.CreateOptsExt{
		createOpts,
		schedulerhints.SchedulerHints{
			Group:           serverProps.SchedulerHints.Group,
			DifferentHost:   serverProps.SchedulerHints.DifferentHost,
			SameHost:        serverProps.SchedulerHints.SameHost,
			Query:           serverProps.SchedulerHints.Query,
			TargetCell:      serverProps.SchedulerHints.TargetCell,
			BuildNearHostIP: serverProps.SchedulerHints.BuildNearHostIP,
		},
	}

	serverOpts, _ := createOpts.ToServerCreateMap()
	i.logger.Debug(openstackServerServiceLogTag, "Creating OpenStack Server with params: %#v", serverOpts)
	server, err := servers.Create(i.computeService, createOpts).Extract()
	if err != nil {
		i.logger.Debug(openstackServerServiceLogTag, "Failed to create OpenStack Server: %#v", err)
		return "", api.NewVMCreationFailedError(true)
	}

	return server.ID, nil
}

func (i OpenStackServerService) CleanUp(id string) {
	if err := i.Delete(id); err != nil {
		i.logger.Debug(openstackServerServiceLogTag, "Failed cleaning up OpenStack Server '%s': %#v", id, err)
	}
}

func (i OpenStackServerService) createNetworksParams(networks Networks) ([]servers.Network, error) {
	var networksParams []servers.Network

	if i.disableNeutron {
		return networksParams, nil
	}

	for _, network := range networks {
		if networkName := network.NetworkName(); networkName != "" {
			net, found, err := i.networkService.FindByName(networkName)
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

func (i OpenStackServerService) createSecurityGroupsParams(networks Networks) ([]string, error) {
	var securityGroupsParams []string

	for _, securityGroup := range networks.SecurityGroupsList() {
		_, found, err := i.securityGroupService.FindByName(securityGroup)
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

	return i.defaultSecurityGroups, nil
}

func (i OpenStackServerService) createUserdataParams(name string, registryEndpoint string, networks Networks) ([]byte, error) {
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
