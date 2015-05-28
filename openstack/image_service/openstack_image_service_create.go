package image

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

func (i OpenStackImageService) Create(imagePath string, description string) (string, error) {
	uuidStr, err := i.uuidGen.Generate()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Generating random OpenStack Image name")
	}

	imageID := fmt.Sprintf("%s-%s", openstackImageNamePrefix, uuidStr)
	// TODO: gophercloud does not support creating images

	return imageID, nil
}

func (i OpenStackImageService) cleanUp(id string) {
	err := i.Delete(id)
	if err != nil {
		i.logger.Debug(openstackImageServiceLogTag, "Failed cleaning up OpenStack Image '%s': %#v", id, err)
	}
}
