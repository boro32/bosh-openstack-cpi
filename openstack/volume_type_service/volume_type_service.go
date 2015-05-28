package volumetype

type Service interface {
	Find(id string) (VolumeType, bool, error)
	FindByName(name string) (VolumeType, bool, error)
}
