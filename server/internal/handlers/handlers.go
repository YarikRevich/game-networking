package handlers

var (
	handlers map[string]func() []byte
)

func AddHandler(name string, callback func() []byte) {
	handlers[name] = callback
}

func CallHandler(name string)[]byte {
	if v, ok := handlers[name]; ok {
		return v()
	}
	return nil
}

func RemoveHandler(name string) {
	delete(handlers, name)
}

func GetHandlers() map[string]func() []byte{
	return handlers
} 