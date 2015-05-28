package network

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
)

func (n OpenStackNetworkService) Find(id string) (Network, bool, error) {
	n.logger.Debug(openstackNetworkServiceLogTag, "Finding OpenStack Network '%s'", id)
	networkItem, err := tenantnetworks.Get(n.computeService, id).Extract()
	if err != nil {
		return Network{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Network '%s'", id)
	}

	network := Network{
		ID:   networkItem.ID,
		Name: networkItem.Name,
	}
	return network, true, nil
}
