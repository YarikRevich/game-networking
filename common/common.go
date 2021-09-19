package common

type Conn interface{
	Read() interface{}
	Write(interface{})
	Close() error 
}