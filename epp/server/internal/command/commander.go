package command

type Commander interface {
	Name() string
	ClTRID() string
}

type Response interface {
}
