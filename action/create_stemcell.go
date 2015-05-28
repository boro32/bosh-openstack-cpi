package action

import (
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
)

const openstackInfrastructure = "openstack"

type CreateStemcell struct {
	imageService image.Service
}

func NewCreateStemcell(
	imageService image.Service,
) CreateStemcell {
	return CreateStemcell{
		imageService: imageService,
	}
}

func (cs CreateStemcell) Run(stemcellPath string, cloudProps StemcellCloudProperties) (StemcellCID, error) {
	var err error
	var description, stemcellID string

	if cloudProps.Infrastructure != openstackInfrastructure {
		return "", bosherr.Errorf("Creating stemcell: Invalid '%s' infrastructure", cloudProps.Infrastructure)
	}

	if cloudProps.ImageUUID != "" {
		return StemcellCID(cloudProps.ImageUUID), nil
	}

	if cloudProps.Name != "" && cloudProps.Version != "" {
		description = fmt.Sprintf("%s/%s", cloudProps.Name, cloudProps.Version)
	}

	stemcellID, err = cs.imageService.Create(stemcellPath, description)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating stemcell")
	}

	return StemcellCID(stemcellID), nil
}
