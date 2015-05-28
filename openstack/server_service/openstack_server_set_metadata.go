package server

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/servers"
)

func (i OpenStackServerService) SetMetadata(id string, serverMetadata Metadata) error {
	updateMetadataOpts := make(servers.MetadataOpts)
	for key, value := range serverMetadata {
		updateMetadataOpts[key] = value.(string)
	}

	i.logger.Debug(openstackServerServiceLogTag, "Setting metadata for OpenStack Server '%s'", id)
	_, err := servers.UpdateMetadata(i.computeService, id, updateMetadataOpts).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return bosherr.WrapErrorf(err, "OpenStack Server '%s' does not exists", id)
		}

		return bosherr.WrapErrorf(err, "Failed to set metadata for OpenStack Server '%s'", id)
	}

	return nil
}
