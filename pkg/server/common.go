package server

type Listener interface {
	WaitForInterrupt() error
	AddHandler(name string, callback func(data interface{}) (interface{}, error))
	CallHandler(name string, data interface{}) (interface{}, error)
}
