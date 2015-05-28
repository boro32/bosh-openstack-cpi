package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
)

func (fip OpenStackFloatingIPService) Find(id string) (FloatingIP, bool, error) {
	fip.logger.Debug(openstackFloatingIPServiceLogTag, "Finding OpenStack Floating IP '%s'", id)
	floatingIPItem, err := floatingip.Get(fip.computeService, id).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return FloatingIP{}, false, nil
		}

		return FloatingIP{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Floating IP '%s'", id)
	}

	floatingIP := FloatingIP{
		ID: floatingIPItem.ID,
		IP: floatingIPItem.IP,
	}
	return floatingIP, true, nil
}
