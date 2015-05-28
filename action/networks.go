package action

import (
	"github.com/frodenas/bosh-registry/client"

	"github.com/frodenas/bosh-openstack-cpi/openstack/server_service"
)

type Networks map[string]Network

type Network struct {
	Type            string                 `json:"type,omitempty"`
	IP              string                 `json:"ip,omitempty"`
	Gateway         string                 `json:"gateway,omitempty"`
	Netmask         string                 `json:"netmask,omitempty"`
	DNS             []string               `json:"dns,omitempty"`
	Default         []string               `json:"default,omitempty"`
	CloudProperties NetworkCloudProperties `json:"cloud_properties,omitempty"`
}

func (ns Networks) AsServerServiceNetworks() server.Networks {
	networks := server.Networks{}

	for netName, network := range ns {
		networks[netName] = server.Network{
			Type:           network.Type,
			IP:             network.IP,
			Gateway:        network.Gateway,
			Netmask:        network.Netmask,
			DNS:            network.DNS,
			Default:        network.Default,
			Network:        network.CloudProperties.Network,
			SecurityGroups: network.CloudProperties.SecurityGroups,
		}
	}

	return networks
}

func (ns Networks) AsRegistryNetworks() registry.NetworksSettings {
	networksSettings := registry.NetworksSettings{}

	for netName, network := range ns {
		networksSettings[netName] = registry.NetworkSettings{
			Type:    network.Type,
			IP:      network.IP,
			Gateway: network.Gateway,
			Netmask: network.Netmask,
			DNS:     network.DNS,
			Default: network.Default,
		}
	}

	return networksSettings
}
