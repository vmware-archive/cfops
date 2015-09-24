package itertools

import (
	"errors"
	"reflect"
	"sync"
)

var (
	NotFuncError    = errors.New("not a func type")
	InvalidArgError = errors.New("invalid argument count")
)

func CEach(iter, f interface{}) {
	if err := validateEachFunction(f); err == nil {
		var wg sync.WaitGroup

		for p := range Iterate(iter) {
			wg.Add(1)

			go func(pp Pair) {
				defer wg.Done()
				runEach(f, p)
			}(p)
		}
		wg.Wait()
	}
}

func Each(iter, f interface{}) {
	if err := validateEachFunction(f); err == nil {

		for p := range Iterate(iter) {
			runEach(f, p)
		}
	}
}

func runEach(f interface{}, p Pair) {
	function := reflect.TypeOf(f)
	args := []reflect.Value{}

	switch function.NumIn() {
	case 1:
		val := reflect.ValueOf(p.Second).Convert(function.In(0))
		args = []reflect.Value{val}

	default:
		val1 := reflect.ValueOf(p.First).Convert(function.In(0))
		val2 := reflect.ValueOf(p.Second).Convert(function.In(1))
		args = []reflect.Value{val1, val2}
	}
	reflect.ValueOf(f).Call(args)
}

func validateEachFunction(f interface{}) (err error) {
	function := reflect.TypeOf(f)

	if function.Kind() != reflect.Func {
		err = NotFuncError
	}

	if function.NumIn() > 2 {
		err = InvalidArgError
	}
	return
}
