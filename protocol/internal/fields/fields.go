package fields

import "reflect"

var (
	fm FieldManager
)

type FieldManager interface {
	SetField(string, interface{})
	GetField(string) (interface{}, bool)
	DeleteField(string)

	SetProtocolValues(interface{})
}

type FM struct {
	protocolValues map[string]interface{}
}

func (f *FM) SetField(name string, value interface{}) {
	f.protocolValues[name] = value
}

func (f *FM) GetField(name string) (interface{}, bool) {
	v, ok := f.protocolValues[name]
	return v, ok
}

func (f *FM) DeleteField(name string) {
	delete(f.protocolValues, name)
}

func (f *FM) SetProtocolValues(t interface{}){
	val := reflect.TypeOf(t).Elem()
	for i := 0; i < val.NumField(); i++ {
		fi := val.Field(i)
		f.protocolValues[fi.Name] = fi.Type.Name()
	}
}

func UseFieldManager() FieldManager {
	if fm == nil {
		fm = &FM{
			protocolValues: make(map[string]interface{}),
		}
	}
	return fm
}
