package auth

type Auth struct{}

func (a *Auth) Login(_, _ string) error {
	return nil
}

func (a *Auth) Logout() error {
	return nil
}

func NewAuth() *Auth {
	return &Auth{}
}
