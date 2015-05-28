package floatingip

type Service interface {
	Find(id string) (FloatingIP, bool, error)
	FindByIP(ipAddress string) (FloatingIP, bool, error)
}
