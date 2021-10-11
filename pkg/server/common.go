package server

type Listener interface {
	WaitForInterrupt() error
	AddHandler(name string, callback func(data []byte) (interface{}, error))
	CallHandler(name string, data []byte) (interface{}, error)
}
