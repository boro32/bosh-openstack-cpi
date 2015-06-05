package extension

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/openstack/networking/v2/extensions"
	"github.com/rackspace/gophercloud/pagination"
)

func (e OpenStackNetworkExtensionService) List() (OpenStackNetworkExtensions, error) {
	availableExtensions := make(map[string]struct{})

	e.logger.Debug(openstackNetworkExtensionServiceLogTag, "Finding OpenStack Network Extensions")
	pager := extensions.List(e.networkService)
	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		extensionList, err := extensions.ExtractExtensions(page)
		if err != nil {
			return false, err
		}

		for _, extension := range extensionList {
			availableExtensions[extension.Alias] = struct{}{}
		}

		return true, nil
	})
	if err != nil {
		return availableExtensions, bosherr.WrapError(err, "Failed to find OpenStack Network Extensions")
	}

	return availableExtensions, nil
}
