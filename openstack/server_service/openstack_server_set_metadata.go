package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (s OpenStackServerService) SetMetadata(id string, serverMetadata Metadata) error {
	updateMetadataOpts := make(servers.MetadataOpts)
	for key, value := range serverMetadata {
		updateMetadataOpts[key] = value.(string)
	}

	s.logger.Debug(openstackServerServiceLogTag, "Setting metadata for OpenStack Server '%s'", id)
	if _, err := servers.UpdateMetadata(s.computeService, id, updateMetadataOpts).Extract(); err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return bosherr.WrapErrorf(err, "Failed to set metadata for OpenStack Server '%s'", id)
	}

	return nil
}
