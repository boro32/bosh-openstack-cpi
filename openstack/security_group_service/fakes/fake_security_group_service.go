package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/security_group_service"
)

type FakeSecurityGroupService struct {
	FindCalled        bool
	FindFound         bool
	FindSecurityGroup securitygroup.SecurityGroup
	FindErr           error

	FindByNameCalled        bool
	FindByNameFound         bool
	FindByNameSecurityGroup securitygroup.SecurityGroup
	FindByNameErr           error
}

func (sg *FakeSecurityGroupService) Find(id string) (securitygroup.VSecurityGroup, bool, error) {
	sg.FindCalled = true
	return sg.FindVSecurityGroup, sg.FindFound, sg.FindErr
}

func (sg *FakeSecurityGroupService) FindByName(name string) (securitygroup.SecurityGroup, bool, error) {
	sg.FindByNameCalled = true
	return sg.FindByNameSecurityGroup, sg.FindByNameFound, sg.FindByNameErr
}
