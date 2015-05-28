package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
)

type DeleteStemcell struct {
	imageService image.Service
}

func NewDeleteStemcell(
	imageService image.Service,
) DeleteStemcell {
	return DeleteStemcell{
		imageService: imageService,
	}
}

func (ds DeleteStemcell) Run(stemcellCID StemcellCID) (interface{}, error) {
	err := ds.imageService.Delete(string(stemcellCID))
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "Deleting stemcell '%s'", stemcellCID)
	}

	return nil, nil
}
