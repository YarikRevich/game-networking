package client

import (
	// "fmt"
	"reflect"
)

type parser struct{}

type IParser interface {
	Parse(src interface{}, dst reflect.Value)
}

func (p *parser) parseMap(src map[string]interface{}, dst reflect.Value) {
	for k, v := range src {
		var f reflect.Value
		if dst.Kind() == reflect.Ptr{
			f = dst.Elem().FieldByName(k)
		}else{
			f = dst.FieldByName(k) 
		}

		switch f.Kind() {
		case reflect.Struct:
			p.parseMap(v.(map[string]interface{}), f)
		case reflect.Map:
			p.parseMap(v.(map[string]interface{}), f)
		case reflect.Ptr:
			f.Elem().Set(reflect.ValueOf(v).Convert(f.Type()))

		default:
			// vf := reflect.ValueOf(v)
			// if f.Kind() != vf.Kind(){
			// 	vf.Convert(f.Type())
			// }
			
			// reflect.TypeOf()
			f.Set(reflect.ValueOf(v))
		}
	}
}

func (p *parser) parseSlice(src []interface{}, dst reflect.Value, dstType reflect.Type) {
	for _, v := range src {
		switch x := v.(type) {
		case map[string]interface{}:
			save := reflect.New(dstType)
			p.parseMap(x, save)
			dst.Elem().Set(reflect.Append(dst.Elem(), save.Elem()))
		case []interface{}:
			save := reflect.New(dstType)
			dstType = dstType.Elem()
			p.parseSlice(x, save, dstType)
			dst.Elem().Set(reflect.Append(dst.Elem(), save.Elem()))
		}
	}
}

func (p *parser) Parse(src interface{}, dst reflect.Value) {
	dstType := reflect.TypeOf(dst.Interface()).Elem()

	if dstType.Kind() == reflect.Slice {
		dstType = dstType.Elem()
	}

	switch v := src.(type) {
	case map[string]interface{}:
		p.parseMap(v, dst)
	case []interface{}:
		p.parseSlice(v, dst, dstType)
	}
}

func NewParser() IParser {
	return new(parser)
}
