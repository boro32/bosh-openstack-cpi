package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
)

func (fip OpenStackNetworkFloatingIPService) Associate(ipAddress string, serverID string) error {
	fip.logger.Debug(openstackNetworkFloatingIPServiceLogTag, "Associating OpenStack Floating IP '%s' to OpenStack Server '%s'", ipAddress, serverID)
	if err := floatingip.Associate(fip.computeService, serverID, ipAddress).ExtractErr(); err != nil {
		return bosherr.WrapErrorf(err, "Failed to associate OpenStack Floating IP '%s' to OpenStack Server '%s'", ipAddress, serverID)
	}

	return nil
}
