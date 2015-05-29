package network

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/networking/v2/networks"
	"github.com/rackspace/gophercloud/pagination"
)

func (n OpenStackNetworkNetworkService) FindByName(name string) (Network, bool, error) {
	var network Network

	n.logger.Debug(openstackNetworkNetworkServiceLogTag, "Finding OpenStack Network '%s'", name)
	listOpts := networks.ListOpts{Name: name}
	pager := networks.List(n.networkService, listOpts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
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
