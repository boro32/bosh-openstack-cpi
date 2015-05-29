package securitygroup

type Service interface {
	Find(id string) (SecurityGroup, bool, error)
	FindByName(name string) (SecurityGroup, bool, error)
}
