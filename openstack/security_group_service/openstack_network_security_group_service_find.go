package securitygroup

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/security/groups"
)

func (sg OpenStackNetworkSecurityGroupService) Find(id string) (SecurityGroup, bool, error) {
	sg.logger.Debug(openstackNetworkSecurityGroupServiceLogTag, "Finding OpenStack Security Group '%s'", id)
	securityGroupItem, err := groups.Get(sg.networkService, id).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return SecurityGroup{}, false, nil
		}

		return SecurityGroup{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Security Group '%s'", id)
	}

	securityGroup := SecurityGroup{
		ID:   securityGroupItem.ID,
		Name: securityGroupItem.Name,
	}
	return securityGroup, true, nil
}
