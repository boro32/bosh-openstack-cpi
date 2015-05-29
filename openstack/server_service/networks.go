package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Networks map[string]Network

func (n Networks) Validate() error {
	var dynamicNetworks, manualNetworks, vipNetworks int

	for _, network := range n {
		if err := network.Validate(); err != nil {
			return err
		}

		switch {
		case network.IsDynamic():
			dynamicNetworks++
		case network.IsManual():
			manualNetworks++
		case network.IsVip():
			vipNetworks++
		}
	}

	if dynamicNetworks == 0 && manualNetworks == 0 {
		return bosherr.Error("At least one 'dynamic' or 'manual' network should be defined")
	}

	if vipNetworks > 1 {
		return bosherr.Error("Only one VIP network is allowed")
	}

	return nil
}

func (n Networks) DNSList() (dnsList []string) {
	nameServers := make(map[string]struct{})
	for _, network := range n {
		for _, nameServer := range network.DNSList() {
			nameServers[nameServer] = struct{}{}
		}
	}

	for dnsItem := range nameServers {
		dnsList = append(dnsList, dnsItem)
	}

	return dnsList
}

func (n Networks) SecurityGroupsList() (securityGroupsList []string) {
	securityGroups := make(map[string]struct{})
	for _, network := range n {
		for _, securityGroup := range network.SecurityGroupsList() {
			securityGroups[securityGroup] = struct{}{}
		}
	}

	for securityGroupItem := range securityGroups {
		securityGroupsList = append(securityGroupsList, securityGroupItem)
	}

	return securityGroupsList
}
