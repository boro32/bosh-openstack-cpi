package floatingip

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions/floatingip"
)

func (fip OpenStackComputeFloatingIPService) Disassociate(ipAddress string, serverID string) error {
	fip.logger.Debug(openstackComputeFloatingIPServiceLogTag, "Desassociating OpenStack Floating IP '%s' from OpenStack Server '%s'", ipAddress, serverID)
	if err := floatingip.Disassociate(fip.computeService, serverID, ipAddress).ExtractErr(); err != nil {
		return bosherr.WrapErrorf(err, "Failed to desassociate OpenStack Floating IP '%s' from OpenStack Server '%s'", ipAddress, serverID)
	}

	return nil
}
