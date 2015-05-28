package flavor

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
	"github.com/rackspace/gophercloud/pagination"
)

func (f OpenStackFlavorService) FindByName(name string) (Flavor, bool, error) {
	var flavor Flavor

	f.logger.Debug(openstackFlavorServiceLogTag, "Finding OpenStack Flavor '%s'", name)
	pager := flavors.ListDetail(f.computeService, flavors.ListOpts{})
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		flavorList, err := flavors.ExtractFlavors(page)
		if err != nil {
			return false, err
		}

		for _, flavorItem := range flavorList {
			if flavorItem.Name == name {
				flavor = Flavor{
					ID:   flavorItem.ID,
					Name: flavorItem.Name,
				}
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return flavor, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Flavor '%s'", name)
	}

	if flavor.ID != "" {
		return flavor, true, nil
	}

	return flavor, false, nil
}
