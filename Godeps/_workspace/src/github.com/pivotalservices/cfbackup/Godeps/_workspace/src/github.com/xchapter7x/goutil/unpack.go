package goutil

import (
	"fmt"
	"reflect"
)

func Unpack(packedValues []interface{}, unpackedPointers ...interface{}) (err error) {
	return UnpackArray(packedValues, unpackedPointers)
}

func UnpackArray(packedValues []interface{}, unpackedPointers []interface{}) (err error) {
	packedValuesLen := len(packedValues)
	unpackedPointersLen := len(unpackedPointers)

	if packedValuesLen == unpackedPointersLen {
		err = mapPackedValuesToUnpackedPointers(packedValues, unpackedPointers)

	} else {
		err = fmt.Errorf("Incorrect argument count: pointers dont match response element count %s != %s", packedValuesLen, unpackedPointersLen)
	}
	return
}

type unpackEmpty struct{}

func (s unpackEmpty) Empty() {}

func Empty() *unpackEmpty {
	return &unpackEmpty{}
}

func mapPackedValuesToUnpackedPointers(packedValues []interface{}, unpackedPointers []interface{}) (err error) {
	for i, packedValue := range packedValues {
		ptrVal := reflect.ValueOf((unpackedPointers)[i])
		ptrElem := ptrVal.Elem()
		ptrElemKind := ptrElem.Kind()
		packedValueReflectValue := reflect.ValueOf(packedValue)
		packedValueReflectValueKind := packedValueReflectValue.Kind()

		if ptrElemKind == packedValueReflectValueKind {
			ptrElem.Set(packedValueReflectValue)
			(unpackedPointers)[i] = ptrElem.Interface()

		} else if !packedValueReflectValue.IsValid() {

			if packedValue != nil {
				err = fmt.Errorf("invalid packed value %s", packedValue)
			}

		} else if packedValueReflectValue.Type() == reflect.ValueOf(fmt.Errorf("")).Type() {
			e := fmt.Sprintf("%s", packedValue)
			*((unpackedPointers)[i].(*error)) = fmt.Errorf(e)
			(unpackedPointers)[i] = fmt.Errorf(e)

		} else if ptrVal.Type() == reflect.ValueOf(Empty()).Type() {
			//do nothing

		} else if packedValueReflectValue.IsValid() {
			err = fmt.Errorf("Incorrect pointer type %s != %s at index %s %s %s %s", ptrElemKind, packedValueReflectValueKind, i, ptrVal.Type(), ptrElem, packedValueReflectValue.Type())
		}
	}
	return
}
