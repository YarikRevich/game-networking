package client

import "reflect"

type parser struct{}

type IParser interface {
	Parse(src interface{}, dst reflect.Value)
}

func (p *parser) parseMap(src map[string]interface{}, dst reflect.Value) {
	for k, v := range src {
		f := dst.Elem().FieldByName(k)
		if f.Kind() == reflect.Ptr {
			f.Elem().Set(reflect.ValueOf(v))
		} else {
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
