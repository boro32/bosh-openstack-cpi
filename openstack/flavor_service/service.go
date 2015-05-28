package flavor

type Service interface {
	Find(id string) (Flavor, bool, error)
	FindByName(name string) (Flavor, bool, error)
}
