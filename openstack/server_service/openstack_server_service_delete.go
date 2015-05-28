package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (i OpenStackServerService) Delete(id string) error {
	i.logger.Debug(openstackServerServiceLogTag, "Deleting OpenStack Server '%s'", id)
	err := servers.Delete(i.computeService, id).ExtractErr()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return bosherr.WrapErrorf(err, "Failed to delete OpenStack Server '%s'", id)
	}

	return nil
}
