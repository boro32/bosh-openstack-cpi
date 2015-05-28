package image

type Service interface {
	Create(imagePath string, description string) (string, error)
	Delete(id string) error
	Find(id string) (Image, bool, error)
}
