package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Networks map[string]Network

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

func (n Network) validateDynamic() error {

	return nil
}

func (n Network) IsVip() bool { return n.Type == "vip" }

func (n Network) validateVip() error {
	if n.IP == "" {
		return bosherr.Error("VIP Network must have an IP address")
	}

	return nil
}

func (n Networks) Validate() error {
	var dnet, vnet bool

	for _, net := range n {
		if net.IsDynamic() {
			if dnet {
				return bosherr.Error("Only one dynamic network is allowed")
			}

			err := net.validateDynamic()
			if err != nil {
				return err
			}

			dnet = true
		}

		if net.IsVip() {
			if vnet {
				return bosherr.Error("Only one VIP network is allowed")
			}

			err := net.validateVip()
			if err != nil {
				return err
			}

			vnet = true
		}
	}

	if !dnet {
		return bosherr.Error("At least one 'dynamic' network should be defined")
	}

	return nil
}

func (n Networks) DynamicNetwork() Network {
	for _, net := range n {
		if net.IsDynamic() {
			// There can only be 1 dynamic network
			return net
		}
	}

	return Network{}
}

func (n Networks) VipNetwork() Network {
	for _, net := range n {
		if net.IsVip() {
			// There can only be 1 vip network
			return net
		}
	}

	return Network{}
}

func (n Networks) DNS() []string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.DNS
}

func (n Networks) Network() string {
	dynamicNetwork := n.DynamicNetwork()

	return dynamicNetwork.Network
}
