package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func (fip OpenStackNetworkFloatingIPService) Find(id string) (FloatingIP, bool, error) {
	fip.logger.Debug(openstackNetworkFloatingIPServiceLogTag, "Finding OpenStack Floating IP '%s'", id)
	floatingIPItem, err := floatingips.Get(fip.networkService, id).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return FloatingIP{}, false, nil
		}

		return FloatingIP{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Floating IP '%s'", id)
	}

	floatingIP := FloatingIP{
		ID: floatingIPItem.ID,
		IP: floatingIPItem.FloatingIP,
	}
	return floatingIP, true, nil
}
