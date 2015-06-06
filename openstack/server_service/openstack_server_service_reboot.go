package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) Reboot(id string) error {
	s.logger.Debug(openstackServerServiceLogTag, "Rebooting OpenStack Server '%s'", id)
	if err := servers.Reboot(s.computeService, id, servers.SoftReboot).ExtractErr(); err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return bosherr.WrapErrorf(err, "Failed to reboot OpenStack Server '%s'", id)
	}

	return nil
}
