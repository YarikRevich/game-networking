package client

type Dialer interface {
	Call(string, interface{}, interface{}, func(error), bool)
	Close() error
}