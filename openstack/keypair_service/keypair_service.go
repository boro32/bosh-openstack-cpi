package keypair

type Service interface {
	Find(id string) (KeyPair, bool, error)
}
