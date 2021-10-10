package client

type Dialer interface {
	Call(string, interface{}, interface{}, func(error), bool) error
	Close() error
}