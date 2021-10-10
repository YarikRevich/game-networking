package server

type Listener interface {
	WaitForInterrupt() error
	AddHandler(name string, callback func(data []byte) ([]byte, error))
	CallHandler(name string, data []byte) ([]byte, error)
}
