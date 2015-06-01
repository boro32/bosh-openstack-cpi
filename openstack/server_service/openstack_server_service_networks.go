package server

import (
	"reflect"
	"sort"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (i OpenStackServerService) AddNetworkConfiguration(id string, networks Networks) error {
	server, found, err := i.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	err = i.associateFloatingIP(server, networks)
	if err != nil {
		return err
	}

	return nil
}

func (i OpenStackServerService) DeleteNetworkConfiguration(id string) error {
	server, found, err := i.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	err = i.disassociateFloatingIP(server)
	if err != nil {
		return err
	}

	return nil
}

func (i OpenStackServerService) UpdateNetworkConfiguration(id string, networks Networks) error {
	server, found, err := i.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}
	i.logger.Debug(openstackServerServiceLogTag, "OpenStack Server %#v", server)

	// TODO: Compare networks/private ip addresses

	err = i.updateSecurityGroups(server, networks)
	if err != nil {
		return err
	}

	err = i.updateFloatingIP(server, networks)
	if err != nil {
		return err
	}

	return nil
}

func (i OpenStackServerService) floatingIP(server *servers.Server) (floatingIP string) {
	for _, addresses := range server.Addresses {
		for _, addressItem := range addresses.([]interface{}) {
			address := addressItem.(map[string]interface{})
			if address["OS-EXT-IPS:type"] == "floating" {
				floatingIP = address["addr"].(string)
				return
			}
		}
	}

	return
}

func (i OpenStackServerService) associateFloatingIP(server *servers.Server, networks Networks) error {
	networkFloatingIP := networks.FloatingIP()
	if networkFloatingIP == "" {
		return nil
	}

	_, found, err := i.floatingIPService.FindByIP(networkFloatingIP)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("OpenStack Floating IP '%s' not found", networkFloatingIP)
	}

	if err := i.floatingIPService.Associate(networkFloatingIP, server.ID); err != nil {
		return err
	}

	return nil
}

func (i OpenStackServerService) updateFloatingIP(server *servers.Server, networks Networks) error {
	networkFloatingIP := networks.FloatingIP()
	serverFloatingIP := i.floatingIP(server)

	if networkFloatingIP == serverFloatingIP {
		return nil
	}

	if serverFloatingIP != "" {
		if err := i.floatingIPService.Disassociate(serverFloatingIP, server.ID); err != nil {
			return err
		}
	}

	if networkFloatingIP != "" {
		_, found, err := i.floatingIPService.FindByIP(networkFloatingIP)
		if err != nil {
			return err
		}
		if !found {
			return bosherr.Errorf("OpenStack Floating IP '%s' not found", networkFloatingIP)
		}

		if err := i.floatingIPService.Associate(networkFloatingIP, server.ID); err != nil {
			return err
		}
	}

	return nil
}

func (i OpenStackServerService) disassociateFloatingIP(server *servers.Server) error {
	serverFloatingIP := i.floatingIP(server)
	if serverFloatingIP == "" {
		return nil
	}

	if err := i.floatingIPService.Disassociate(serverFloatingIP, server.ID); err != nil {
		return err
	}

	return nil
}

func (i OpenStackServerService) updateSecurityGroups(server *servers.Server, networks Networks) error {
	var serverSecurityGroups []string
	for _, secgrp := range server.SecurityGroups {
		serverSecurityGroups = append(serverSecurityGroups, secgrp["name"].(string))
	}

	networkSecurityGroups := networks.SecurityGroupsList()
	if len(networkSecurityGroups) == 0 {
		networkSecurityGroups = i.defaultSecurityGroups
	}

	sort.Strings(networkSecurityGroups)
	sort.Strings(serverSecurityGroups)
	if reflect.DeepEqual(networkSecurityGroups, serverSecurityGroups) {
		return nil
	}

	i.logger.Debug(openstackServerServiceLogTag, "Changing OpenStack security groups for OpenStack server '%s' not supported", server.Name)
	i.logger.Debug(openstackServerServiceLogTag, "Server Security Groups: %#v", serverSecurityGroups)
	i.logger.Debug(openstackServerServiceLogTag, "Network Security Groups: %#v", networkSecurityGroups)
	return api.NotSupportedError{}
}
