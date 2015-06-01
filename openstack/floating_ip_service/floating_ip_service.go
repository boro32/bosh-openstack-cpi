package floatingip

type Service interface {
	Associate(ipAddress string, serverID string) error
	Disassociate(ipAddress string, serverID string) error
	Find(id string) (FloatingIP, bool, error)
	FindByIP(ipAddress string) (FloatingIP, bool, error)
}
