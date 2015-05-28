package keypair

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/rackspace/gophercloud/rackspace/compute/v2/keypairs"
)

func (k OpenStackKeyPairService) Find(id string) (KeyPair, bool, error) {
	k.logger.Debug(openstackKeyPairServiceLogTag, "Finding OpenStack KeyPair '%s'", id)
	keypairItem, err := keypairs.Get(k.computeService, id).Extract()
	if err != nil {
		return KeyPair{}, false, bosherr.WrapErrorf(err, "Failed to find OpenStack KeyPair '%s'", id)
	}

	keypair := KeyPair{
		Name: keypairItem.Name,
	}
	return keypair, true, nil
}
