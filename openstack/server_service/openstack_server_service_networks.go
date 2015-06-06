package server

import (
	"reflect"
	"sort"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/api"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) AddNetworkConfiguration(id string, networks Networks) error {
	server, found, err := s.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	err = s.associateFloatingIP(server, networks)
	if err != nil {
		return err
	}

	return nil
}

func (s OpenStackServerService) DeleteNetworkConfiguration(id string) error {
	server, found, err := s.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}

	err = s.disassociateFloatingIP(server)
	if err != nil {
		return err
	}

	return nil
}

func (s OpenStackServerService) UpdateNetworkConfiguration(id string, networks Networks) error {
	server, found, err := s.Find(id)
	if err != nil {
		return err
	}
	if !found {
		return api.NewVMNotFoundError(id)
	}
	s.logger.Debug(openstackServerServiceLogTag, "OpenStack Server %#v", server)

	// TODO: Compare networks/private ip addresses

	err = s.updateSecurityGroups(server, networks)
	if err != nil {
		return err
	}

	err = s.updateFloatingIP(server, networks)
	if err != nil {
		return err
	}

	return nil
}

func (s OpenStackServerService) floatingIP(server *servers.Server) (floatingIP string) {
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

func (s OpenStackServerService) associateFloatingIP(server *servers.Server, networks Networks) error {
	networkFloatingIP := networks.FloatingIP()
	if networkFloatingIP == "" {
		return nil
	}

	_, found, err := s.floatingIPService.FindByIP(networkFloatingIP)
	if err != nil {
		return err
	}
	if !found {
		return bosherr.Errorf("OpenStack Floating IP '%s' not found", networkFloatingIP)
	}

	if err := s.floatingIPService.Associate(networkFloatingIP, server.ID); err != nil {
		return err
	}

	return nil
}

func (s OpenStackServerService) updateFloatingIP(server *servers.Server, networks Networks) error {
	networkFloatingIP := networks.FloatingIP()
	serverFloatingIP := s.floatingIP(server)

	if networkFloatingIP == serverFloatingIP {
		return nil
	}

	if serverFloatingIP != "" {
		if err := s.floatingIPService.Disassociate(serverFloatingIP, server.ID); err != nil {
			return err
		}
	}

	if networkFloatingIP != "" {
		_, found, err := s.floatingIPService.FindByIP(networkFloatingIP)
		if err != nil {
			return err
		}
		if !found {
			return bosherr.Errorf("OpenStack Floating IP '%s' not found", networkFloatingIP)
		}

		if err := s.floatingIPService.Associate(networkFloatingIP, server.ID); err != nil {
			return err
		}
	}

	return nil
}

func (s OpenStackServerService) disassociateFloatingIP(server *servers.Server) error {
	serverFloatingIP := s.floatingIP(server)
	if serverFloatingIP == "" {
		return nil
	}

	if err := s.floatingIPService.Disassociate(serverFloatingIP, server.ID); err != nil {
		return err
	}

	return nil
}

func (s OpenStackServerService) updateSecurityGroups(server *servers.Server, networks Networks) error {
	var serverSecurityGroups []string
	for _, secgrp := range server.SecurityGroups {
		serverSecurityGroups = append(serverSecurityGroups, secgrp["name"].(string))
	}

	networkSecurityGroups := networks.SecurityGroupsList()
	if len(networkSecurityGroups) == 0 {
		networkSecurityGroups = s.defaultSecurityGroups
	}

	sort.Strings(networkSecurityGroups)
	sort.Strings(serverSecurityGroups)
	if reflect.DeepEqual(networkSecurityGroups, serverSecurityGroups) {
		return nil
	}

	s.logger.Debug(openstackServerServiceLogTag, "Changing OpenStack security groups for OpenStack server '%s' not supported", server.Name)
	s.logger.Debug(openstackServerServiceLogTag, "Server Security Groups: %#v", serverSecurityGroups)
	s.logger.Debug(openstackServerServiceLogTag, "Network Security Groups: %#v", networkSecurityGroups)
	return api.NotSupportedError{}
}
