package command

type Name string

func (n Name) String() string {
	return string(n)
}

const (
	HelloCommand Name = "hello"
)

type Command interface {
	Name() Name
	ClTRID() string
	AsBytes() []byte
}

type Response interface {
}
