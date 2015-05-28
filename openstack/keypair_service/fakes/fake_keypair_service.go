package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/keypair_service"
)

type FakeKeyPairService struct {
	FindCalled  bool
	FindFound   bool
	FindKeyPair keypair.KeyPair
	FindErr     error
}

func (k *FakeKeyPairService) Find(id string) (keypair.KeyPair, bool, error) {
	k.FindCalled = true
	return k.FindKeyPair, k.FindFound, k.FindErr
}
