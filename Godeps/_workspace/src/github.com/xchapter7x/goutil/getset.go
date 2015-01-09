package goutil

import (
	"reflect"
)

type GetSetter interface {
	Get(n string) interface{}
	Set(n string, val interface{})
}

type GetSet struct{}

func (s *GetSet) Get(i interface{}, n string) interface{} {
	var field reflect.Value
	v := reflect.ValueOf(i)

	if reflect.TypeOf(i).Kind() == reflect.Struct {
		field = v.FieldByName(n)
	} else {
		field = v.Elem().FieldByName(n)
	}
	return field.Interface()
}

func (s *GetSet) Set(i interface{}, n string, val interface{}) {
	var field reflect.Value
	v := reflect.ValueOf(i)

	if reflect.TypeOf(i).Kind() == reflect.Struct {
		field = v.FieldByName(n)
	} else {
		field = v.Elem().FieldByName(n)
	}
	va := reflect.ValueOf(val)
	field.Set(va)
}
