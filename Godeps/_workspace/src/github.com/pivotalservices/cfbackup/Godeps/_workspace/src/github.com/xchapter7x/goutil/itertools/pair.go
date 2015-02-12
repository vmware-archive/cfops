package itertools

import (
	"reflect"
)

type Pair struct {
	First  interface{}
	Second interface{}
}

func unpack(pairVal, ptr interface{}) {
	ptrVal := reflect.ValueOf(ptr)
	ptrElem := ptrVal.Elem()
	ptrElem.Set(reflect.ValueOf(pairVal))
	ptr = ptrElem.Interface()
}

func PairUnPack(pair Pair, first, second interface{}) {
	unpack(pair.First, first)
	unpack(pair.Second, second)
}
