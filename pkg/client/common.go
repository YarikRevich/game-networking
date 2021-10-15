package client

type Dialer interface {
	Call(string, interface{}, interface{})
	Close() error
	IsConnected() bool
}