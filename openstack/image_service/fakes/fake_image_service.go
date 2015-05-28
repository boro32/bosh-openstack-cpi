package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/image_service"
)

type FakeImageService struct {
	CreateCalled      bool
	CreateErr         error
	CreateID          string
	CreateImagePath   string
	CreateDescription string

	DeleteCalled bool
	DeleteErr    error

	FindCalled bool
	FindFound  bool
	FindImage  image.Image
	FindErr    error
}

func (i *FakeImageService) Create(imagePath string, description string) (string, error) {
	i.CreateCalled = true
	i.CreateImagePath = imagePath
	i.CreateDescription = description
	return i.CreateID, i.CreateErr
}

func (i *FakeImageService) Delete(id string) error {
	i.DeleteCalled = true
	return i.DeleteErr
}

func (i *FakeImageService) Find(id string) (image.Image, bool, error) {
	i.FindCalled = true
	return i.FindImage, i.FindFound, i.FindErr
}
