package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/network_service"
)

type FakeNetworkService struct {
	FindCalled  bool
	FindFound   bool
	FindNetwork network.Network
	FindErr     error
}

func (n *FakeNetworkService) Find(id string) (network.Network, bool, error) {
	n.FindCalled = true
	return n.FindNetwork, n.FindFound, n.FindErr
}
