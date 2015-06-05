package extension

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack/compute/v2/extensions"
)

func (e OpenStackComputeExtensionService) Find(alias string) (string, bool, error) {
	e.logger.Debug(openstackComputeExtensionServiceLogTag, "Finding OpenStack Compute Extension '%s'", alias)
	extensionItem, err := extensions.Get(e.computeService, alias).Extract()
	if err != nil {
		errCode, _ := err.(*gophercloud.UnexpectedResponseCodeError)
		if errCode.Actual == 404 {
			return "", false, nil
		}

		return "", false, bosherr.WrapErrorf(err, "Failed to find OpenStack Compute Extension '%s'", alias)
	}

	return extensionItem.Name, true, nil
}
