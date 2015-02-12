package goutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/goutil"
)

type testGetSet struct {
	GetSet
	SampleIntField1    int
	SampleStringField2 string
}

func (s *testGetSet) Get(n string) interface{} {
	return s.GetSet.Get(s, n)
}

func (s *testGetSet) Set(n string, val interface{}) {
	s.GetSet.Set(s, n, val)
}

var _ = Describe("GetSet", func() {
	Describe("Struct composed of GetSet", func() {
		var (
			testGS        *testGetSet
			intDefault    int    = 100
			stringDefault string = "hello"
		)

		BeforeEach(func() {
			testGS = &testGetSet{
				SampleIntField1:    intDefault,
				SampleStringField2: stringDefault,
			}
		})

		Context("Get method", func() {
			It("Should provide allow to get the value in a structs field ", func() {
				fld := testGS.GetSet.Get(testGS, "SampleIntField1")
				Ω(fld.(int)).Should(Equal(intDefault))
			})

			It("Should provide allow to get the value, regardless of type, in a structs field ", func() {
				fld := testGS.GetSet.Get(testGS, "SampleStringField2")
				Ω(fld.(string)).Should(Equal(stringDefault))
			})
		})

		Context("Set Method", func() {
			It("Should allow us to set the structs field value", func() {
				controlValue := "some other string"
				testGS.GetSet.Set(testGS, "SampleStringField2", controlValue)
				fld := testGS.GetSet.Get(testGS, "SampleStringField2")
				Ω(fld.(string)).ShouldNot(Equal(stringDefault))
				Ω(fld.(string)).Should(Equal(controlValue))
			})
		})

		Describe("Using a interface type", func() {
			var interfaceObject GetSetter

			BeforeEach(func() {
				interfaceObject = func() GetSetter {
					return testGS
				}()
			})

			It("Should allow us to get a field value from the underlying struct, but from the interface object", func() {
				v := interfaceObject.Get("SampleStringField2")
				Ω(v.(string)).Should(Equal(testGS.SampleStringField2))
			})

			It("Should allow us to set a field value on the underlying struct, but from the interface object", func() {
				control := "some text"
				interfaceObject.Set("SampleStringField2", control)
				v := interfaceObject.Get("SampleStringField2")
				Ω(control).Should(Equal(testGS.SampleStringField2))
				Ω(v.(string)).Should(Equal(testGS.SampleStringField2))

			})

		})
	})
})
