package handlers 

type HandlerManager struct {
	handlers map[string]func(...interface{})
}

func (hm *HandlerManager) AddHandler(name string, callback func(...interface{})){
	hm.handlers[name] = callback
}

func (hm *HandlerManager) RemoveHandler(name string){
	delete(hm.handlers, name)
}

func New()*HandlerManager{
	return &HandlerManager{
		handlers: make(map[string]func(...interface{})),
	}
}
