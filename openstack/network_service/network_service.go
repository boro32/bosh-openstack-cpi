package network

type Service interface {
	Find(id string) (Network, bool, error)
	FindByName(name string) (Network, bool, error)
}
