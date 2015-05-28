package snapshot

type Service interface {
	Create(volumeID string, description string) (string, error)
	Delete(id string) error
	Find(id string) (Snapshot, bool, error)
}
