package itertools

import (
	"reflect"
	"sync"
)

func Filter(iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	out = filter(passThrough, iter, f)
	return
}

func CFilter(iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	out = cFilter(passThrough, iter, f)
	return
}

func FilterFalse(iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	out = filter(falsify, iter, f)
	return
}

func CFilterFalse(iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	out = cFilter(falsify, iter, f)
	return
}

func falsify(in bool) bool {
	return !in
}

func passThrough(in bool) bool {
	return in
}

func pipeToFilterChannel(p Pair, out chan Pair, f interface{}, functor func(bool) bool) {
	args := []reflect.Value{reflect.ValueOf(p.First), reflect.ValueOf(p.Second)}

	if functor(reflect.ValueOf(f).Call(args)[0].Bool()) {
		out <- p
	}
}

func filter(functor func(bool) bool, iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	var wg sync.WaitGroup
	out = make(chan Pair, GetIterBuffer())
	wg.Add(1)

	go func() {
		defer close(out)
		defer wg.Done()

		for p := range Iterate(iter) {
			pipeToFilterChannel(p, out, f, functor)
		}
	}()
	wg.Wait()
	return
}

func cFilter(functor func(bool) bool, iter interface{}, f func(first, second interface{}) bool) (out chan Pair) {
	var wg1 sync.WaitGroup
	out = make(chan Pair, GetIterBuffer())
	wg1.Add(1)

	go func() {
		defer close(out)
		defer wg1.Done()
		var wg2 sync.WaitGroup

		for p := range Iterate(iter) {
			wg2.Add(1)

			go func(pp Pair) {
				defer wg2.Done()
				pipeToFilterChannel(pp, out, f, functor)
			}(p)
		}
		wg2.Wait()
	}()
	wg1.Wait()
	return
}
