package goutil

import "reflect"

func createReflectValueArgsArray(iargs []interface{}) (args []reflect.Value) {
	for _, arg := range iargs {
		args = append(args, reflect.ValueOf(arg))
	}
	return
}

func createInterfaceArrayFromValuesArray(responseValuesArray []reflect.Value) (responseInterfaceArray []interface{}) {
	for _, ri := range responseValuesArray {
		responseInterfaceArray = append(responseInterfaceArray, ri.Interface())
	}
	return
}

func findErrorValue(responseInterfaceArray []interface{}) (err error) {
	for _, res := range responseInterfaceArray {
		if e, ok := res.(error); ok {
			err = e
		}
	}
	return
}

func NewChain(err error) (chain *Chain) {
	return &Chain{
		Error: err,
	}
}

type Chain struct {
	Error error
}

func (s *Chain) Returns(args ...interface{}) []interface{} {
	return args
}

func (s *Chain) Call(functor interface{}, iargs ...interface{}) (responseInterfaceArray []interface{}, err error) {
	responseInterfaceArray, err = CallChain(s.Error, functor, iargs...)
	s.Error = err
	return
}

func (s *Chain) CallP(responseInterfaceArray []interface{}, functor interface{}, iargs ...interface{}) (err error) {
	err = CallChainP(s.Error, responseInterfaceArray, functor, iargs...)
	s.Error = err
	return
}

func CallChain(preverr error, functor interface{}, iargs ...interface{}) (responseInterfaceArray []interface{}, err error) {
	if err = preverr; err == nil {
		args := createReflectValueArgsArray(iargs)
		responseValuesArray := reflect.ValueOf(functor).Call(args)
		responseInterfaceArray = createInterfaceArrayFromValuesArray(responseValuesArray)
		err = findErrorValue(responseInterfaceArray)
	}
	return
}

func CallChainP(preverr error, responseInterfaceArray []interface{}, functor interface{}, iargs ...interface{}) (err error) {
	var res []interface{}
	res, err = CallChain(preverr, functor, iargs...)

	if e := UnpackArray(res, responseInterfaceArray); err == nil {
		err = e
	}
	return
}
