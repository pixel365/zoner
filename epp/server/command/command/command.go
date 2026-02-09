package command

type CommandName string

func (c CommandName) String() string {
	return string(c)
}

const (
	Greeting CommandName = "greeting"
	Login    CommandName = "login"
	Logout   CommandName = "logout"
	Hello    CommandName = "hello"
	Check    CommandName = "check"
	Info     CommandName = "info"
	Poll     CommandName = "poll"
	Transfer CommandName = "transfer"
	Create   CommandName = "create"
)

type Commander interface {
	Name() CommandName
	Validate() error
	NeedAuth() bool
}

type Responser interface {
	AsBytes() []byte
}
