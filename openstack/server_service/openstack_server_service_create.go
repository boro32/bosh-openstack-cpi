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
	userdataParams, err := i.createUserdataParams(serverName, registryEndpoint, networks)
	if err != nil {
		return "", err
	}

	configDrive := true
	if serverProps.DisableConfigDrive {
		configDrive = false
	}

	createOpts = &servers.CreateOpts{
		Name:             serverName,
		ImageRef:         serverProps.ImageID,
		FlavorRef:        serverProps.FlavorID,
		SecurityGroups:   []string{},
		UserData:         userdataParams,
		AvailabilityZone: serverProps.AvailabilityZone,
		//Networks: "",
		ConfigDrive: configDrive,
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
	err := i.Delete(id)
	if err != nil {
		i.logger.Debug(openstackServerServiceLogTag, "Failed cleaning up OpenStack Server '%s': %#v", id, err)
	}
}

func (i OpenStackServerService) createUserdataParams(name string, registryEndpoint string, networks Networks) ([]byte, error) {
	openstackServerName := OpenStackUserDataServerName{Name: name}
	openstackRegistryEndpoint := OpenStackUserDataRegistryEndpoint{Endpoint: registryEndpoint}
	openstackUserData := OpenStackUserData{Server: openstackServerName, Registry: openstackRegistryEndpoint}

	if networkDNS := networks.DNS(); len(networkDNS) > 0 {
		openstackUserData.DNS = OpenStackUserDataDNSItems{NameServer: networkDNS}
	}

	userData, err := json.Marshal(openstackUserData)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Marshalling user data")
	}

	return userData, nil
}
