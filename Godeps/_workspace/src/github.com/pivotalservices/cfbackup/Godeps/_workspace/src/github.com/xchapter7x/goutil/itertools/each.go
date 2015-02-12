package itertools

import (
	"reflect"
	"sync"
)

func CEach(iter, f interface{}) {
	var wg sync.WaitGroup

	for p := range Iterate(iter) {
		wg.Add(1)

		go func(pp Pair) {
			defer wg.Done()
			args := []reflect.Value{reflect.ValueOf(pp.First), reflect.ValueOf(pp.Second)}
			reflect.ValueOf(f).Call(args)
		}(p)
	}
	wg.Wait()
}

func Each(iter, f interface{}) {
	for p := range Iterate(iter) {
		args := []reflect.Value{reflect.ValueOf(p.First), reflect.ValueOf(p.Second)}
		reflect.ValueOf(f).Call(args)
	}
}
