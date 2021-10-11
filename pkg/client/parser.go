package client

import "reflect"

func ParseToDst(src interface{}, dst reflect.Value) {
	if srcMap, ok := src.(map[string]interface{}); ok {
		for k, v := range srcMap {
			f := dst.Elem().FieldByName(k)
			if f.Kind() == reflect.Ptr {
				f.Elem().Set(reflect.ValueOf(v))
			} else {
				f.Set(reflect.ValueOf(v))
			}
		}
	} else {
		dst.Elem().Set(reflect.ValueOf(src))
	}
}
