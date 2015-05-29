package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/rackspace/gophercloud/pagination"
)

func (fip OpenStackNetworkFloatingIPService) FindByIP(ipAddress string) (FloatingIP, bool, error) {
	var floatingIP FloatingIP

	fip.logger.Debug(openstackNetworkFloatingIPServiceLogTag, "Finding OpenStack Floating IP Address '%s'", ipAddress)
	listOpts := floatingips.ListOpts{FloatingIP: ipAddress}
	pager := floatingips.List(fip.networkService, listOpts)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		floatingIPList, err := floatingips.ExtractFloatingIPs(page)
		if err != nil {
			return false, err
		}

		for _, floatingIPItem := range floatingIPList {
			if floatingIPItem.FloatingIP == ipAddress {
				floatingIP = FloatingIP{
					ID: floatingIPItem.ID,
					IP: floatingIPItem.FloatingIP,
				}
				return false, nil
			}
		}

		return true, nil
	})
	if err != nil {
		return floatingIP, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Floating IP Address '%s'", ipAddress)
	}

	if floatingIP.ID != "" {
		return floatingIP, true, nil
	}

	return floatingIP, false, nil
}
