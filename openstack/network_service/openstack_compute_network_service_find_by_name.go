package network

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/tenantnetworks"
	"github.com/rackspace/gophercloud/pagination"
)

func (n OpenStackComputeNetworkService) FindByName(name string) (Network, bool, error) {
	var network Network

	n.logger.Debug(openstackComputeNetworkServiceLogTag, "Finding OpenStack Network '%s'", name)
	pager := tenantnetworks.List(n.computeService)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := tenantnetworks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, networkItem := range networkList {
			if networkItem.Name == name {
				network = Network{
					ID:   networkItem.ID,
					Name: networkItem.Name,
				}
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return network, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Network '%s'", name)
	}

	if network.ID != "" {
		return network, true, nil
	}

	return network, false, nil
}
