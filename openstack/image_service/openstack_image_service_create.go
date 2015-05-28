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

	// TODO: gophercloud does not support creating images
	imageID := fmt.Sprintf("%s-%s", openstackImageNamePrefix, uuidStr)

	return imageID, nil
}

func (i OpenStackImageService) cleanUp(id string) {
	if err := i.Delete(id); err != nil {
		i.logger.Debug(openstackImageServiceLogTag, "Failed cleaning up OpenStack Image '%s': %#v", id, err)
	}
}
