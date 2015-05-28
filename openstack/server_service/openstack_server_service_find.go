package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (i OpenStackServerService) Find(id string) (*servers.Server, bool, error) {
	i.logger.Debug(openstackServerServiceLogTag, "Finding OpenStack Server '%s'", id)
	instance, err := servers.Get(i.computeService, id).Extract()
	if err != nil {
		return &servers.Server{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack Server '%s'", id)
	}

	return instance, true, nil
}
