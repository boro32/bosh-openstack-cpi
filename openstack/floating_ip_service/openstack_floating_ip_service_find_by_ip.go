package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
	"github.com/rackspace/gophercloud/pagination"
)

func (fip OpenStackFloatingIPService) FindByIP(ipAddress string) (FloatingIP, bool, error) {
	var floatingIP FloatingIP

	fip.logger.Debug(openstackFloatingIPServiceLogTag, "Finding OpenStack Floating IP Address '%s'", ipAddress)
	pager := floatingip.List(fip.computeService)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		floatingIPList, err := floatingip.ExtractFloatingIPs(page)
		if err != nil {
			return false, err
		}

		for _, floatingIPItem := range floatingIPList {
			if floatingIPItem.IP == ipAddress {
				floatingIP = FloatingIP{
					ID: floatingIPItem.ID,
					IP: floatingIPItem.IP,
				}
			}
			return false, nil
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
