package common

type Dialer interface {
	Call(string, interface{}, interface{}, func(error), bool) error
	Close() error
}

type Listener interface {
	WaitForInterrupt() error
	AddHandler(name string, callback func(data interface{}) ([]byte, error))
	CallHandler(name string, data interface{}) ([]byte, error)
	Close() error
}
