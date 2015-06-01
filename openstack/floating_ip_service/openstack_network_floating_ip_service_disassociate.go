package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func (fip OpenStackNetworkFloatingIPService) Disassociate(ipAddress string, portID string) error {
	fip.logger.Debug(openstackNetworkFloatingIPServiceLogTag, "Desassociating OpenStack Floating IP '%s' from OpenStack Port '%s'", ipAddress, portID)
	updateOpts := floatingips.UpdateOpts{PortID: ""}
	_, err := floatingips.Update(fip.networkService, ipAddress, updateOpts).Extract()
	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to desassociate OpenStack Floating IP '%s' from OpenStack Port '%s'", ipAddress, portID)
	}

	return nil
}
