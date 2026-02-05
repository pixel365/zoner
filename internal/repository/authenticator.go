package repository

type Authenticator interface {
	Login(string, string) error
	Logout() error
}
