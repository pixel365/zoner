package command

type Command interface {
	Name() string
	Validate() error
}

type Responser interface {
	AsBytes() []byte
}
