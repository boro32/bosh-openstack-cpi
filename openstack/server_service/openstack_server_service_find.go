package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) Find(id string) (*servers.Server, bool, error) {
	s.logger.Debug(openstackServerServiceLogTag, "Finding OpenStack Server '%s'", id)
	instance, err := servers.Get(s.computeService, id).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return &servers.Server{}, false, nil
		}

		return &servers.Server{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Server '%s'", id)
	}

	return instance, true, nil
}
