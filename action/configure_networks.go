package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"
	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"

	"github.com/frodenas/bosh-registry/client"
)

type ConfigureNetworks struct {
	serverService  server.Service
	registryClient registry.Client
}

func NewConfigureNetworks(
	serverService server.Service,
	registryClient registry.Client,
) ConfigureNetworks {
	return ConfigureNetworks{
		serverService:  serverService,
		registryClient: registryClient,
	}
}

func (rv ConfigureNetworks) Run(vmCID VMCID, networks Networks) (interface{}, error) {
	// Parse networks
	vmNetworks := networks.AsServerServiceNetworks()
	if err := vmNetworks.Validate(); err != nil {
		return "", bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update networks
	err := rv.serverService.UpdateNetworkConfiguration(string(vmCID), vmNetworks)
	if err != nil {
		if _, ok := err.(api.CloudError); ok {
			return nil, err
		}
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Read VM agent settings
	agentSettings, err := rv.registryClient.Fetch(string(vmCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	// Update VM agent settings
	agentNetworks := networks.AsRegistryNetworks()
	newAgentSettings := agentSettings.ConfigureNetworks(agentNetworks)
	err = rv.registryClient.Update(string(vmCID), newAgentSettings)
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Configuring networks for vm '%s'", vmCID)
	}

	return nil, nil
}
