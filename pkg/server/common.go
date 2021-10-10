package server

type Listener interface {
	WaitForInterrupt() error
	AddHandler(name string, callback func(data interface{}) ([]byte, error))
	CallHandler(name string, data interface{}) ([]byte, error)
	Close() error
}
