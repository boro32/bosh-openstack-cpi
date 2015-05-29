package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Network struct {
	Type           string
	IP             string
	Gateway        string
	Netmask        string
	DNS            []string
	Default        []string
	Network        string
	SecurityGroups []string
}

func (n Network) IsDynamic() bool { return n.Type == "dynamic" }

func (n Network) IsManual() bool { return n.Type == "manual" }

func (n Network) IsVip() bool { return n.Type == "vip" }

func (n Network) Validate() error {
	switch {
	case n.IsDynamic():
	case n.IsManual():
		if n.Network == "" {
			return bosherr.Error("Manual Networks must provide a Network name")
		}
	case n.IsVip():
		if n.IP == "" {
			return bosherr.Error("VIP Networks must provide an IP Address")
		}
	default:
		return bosherr.Errorf("Network type '%s' not supported", n.Type)
	}

	return nil
}

func (n Network) IPAddress() string {
	if n.IsManual() || n.IsVip() {
		return n.IP
	}

	return ""
}

func (n Network) DNSList() (dnsList []string) {
	if n.IsDynamic() || n.IsManual() {
		return n.DNS
	}

	return dnsList
}

func (n Network) NetworkName() string {
	if n.IsDynamic() || n.IsManual() {
		return n.Network
	}

	return ""
}

func (n Network) SecurityGroupsList() (securityGroupsList []string) {
	if n.IsDynamic() || n.IsManual() {
		return n.SecurityGroups
	}

	return securityGroupsList
}
