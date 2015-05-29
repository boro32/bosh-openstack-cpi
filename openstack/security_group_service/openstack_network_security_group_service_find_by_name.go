package securitygroup

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/rackspace/gophercloud/pagination"
)

func (sg OpenStackNetworkSecurityGroupService) FindByName(name string) (SecurityGroup, bool, error) {
	var securityGroup SecurityGroup

	sg.logger.Debug(openstackNetworkSecurityGroupServiceLogTag, "Finding OpenStack Security Group '%s'", name)
	listOpts := groups.ListOpts{Name: name}
	pager := groups.List(sg.networkService, listOpts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		securityGroupList, err := groups.ExtractGroups(page)
		if err != nil {
			return false, err
		}

		for _, securityGroupItem := range securityGroupList {
			if securityGroupItem.Name == name {
				securityGroup = SecurityGroup{
					ID:   securityGroupItem.ID,
					Name: securityGroupItem.Name,
				}
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return securityGroup, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Security Group '%s'", name)
	}

	if securityGroup.ID != "" {
		return securityGroup, true, nil
	}

	return securityGroup, false, nil
}
