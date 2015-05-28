package flavor

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/flavors"
)

func (f OpenStackFlavorService) Find(id string) (Flavor, bool, error) {
	f.logger.Debug(openstackFlavorServiceLogTag, "Finding OpenStack Flavor '%s'", id)
	flavorItem, err := flavors.Get(f.computeService, id).Extract()
	if err != nil {
		return Flavor{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Flavor '%s'", id)
	}

	flavor := Flavor{
		ID:   flavorItem.ID,
		Name: flavorItem.Name,
	}
	return flavor, true, nil
}
