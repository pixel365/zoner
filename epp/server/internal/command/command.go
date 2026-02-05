package command

type Command interface {
	Name() string
	Validate() error
}

type Response interface {
	AsBytes() []byte
}
