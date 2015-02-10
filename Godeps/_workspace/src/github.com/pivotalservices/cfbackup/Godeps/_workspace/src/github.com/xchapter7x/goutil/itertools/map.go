package itertools

import (
	"reflect"
	"sync"
)

func CMap(iter, f interface{}) (out chan Pair){
	var wg sync.WaitGroup
	iterationCounter := 0
	out = make(chan Pair, GetIterBuffer())
	defer close(out)

	for p := range Iterate(iter) {
		wg.Add(1)

		go func(pp Pair) {
			defer wg.Done()
			args := []reflect.Value{reflect.ValueOf(pp.First), reflect.ValueOf(pp.Second)}
			functorResponseValue := reflect.ValueOf(f).Call(args)
			out <- Pair{iterationCounter, functorResponseValue}
			iterationCounter+=1
		}(p)
	}
	wg.Wait()
	return
}

func Map(iter, f interface{}) (out chan Pair){
	iterationCounter := 0
	out = make(chan Pair, GetIterBuffer())
	defer close(out)

	for p := range Iterate(iter) {
		args := []reflect.Value{reflect.ValueOf(p.First), reflect.ValueOf(p.Second)}
		functorResponseValue := reflect.ValueOf(f).Call(args)
		out <- Pair{iterationCounter, functorResponseValue}
		iterationCounter+=1
	}
	return
}
