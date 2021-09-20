package common

type Conn interface{
	Call(string) interface{}
	Send(string, interface{})

	Error() error
	
	// Read() ([]byte, error)
	// Write(interface{}) error
	Close() error 
}