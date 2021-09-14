package handlers

var (
	handlers = map[string]func(data interface{}) []byte{}
)

func AddHandler(name string, callback func(data interface{}) []byte) {
	handlers[name] = callback
}

func CallHandler(name string, data interface{})[]byte {
	if v, ok := handlers[name]; ok {
		return v(data)
	}
	return nil
}

func RemoveHandler(name string) {
	delete(handlers, name)
}

func GetHandlers() map[string]func(data interface{}) []byte{
	return handlers
} 