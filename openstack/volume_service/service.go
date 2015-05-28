package volume

type Service interface {
	Create(size int, volumeType string, zone string) (string, error)
	Delete(id string) error
	Find(id string) (Volume, bool, error)
}
