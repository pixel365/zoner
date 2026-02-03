package command

type Command interface {
	Name() string
	ClTRID() string
	AsBytes() []byte
	Validate() error
}

type Response interface {
}
