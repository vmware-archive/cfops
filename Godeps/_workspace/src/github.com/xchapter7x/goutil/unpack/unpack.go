package unpack

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

		} else {
			err = fmt.Errorf("Incorrect pointer type %s != %s", ptrElemKind, packedValueReflectValueKind)
		}
	}
	return
}
