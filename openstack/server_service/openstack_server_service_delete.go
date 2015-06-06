package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) Delete(id string) error {
	s.logger.Debug(openstackServerServiceLogTag, "Deleting OpenStack Server '%s'", id)
	if err := servers.Delete(s.computeService, id).ExtractErr(); err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return bosherr.WrapErrorf(err, "Failed to delete OpenStack Server '%s'", id)
	}

	return nil
}
