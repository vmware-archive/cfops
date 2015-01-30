package goutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/goutil"
)

var _ = Describe("unpack package", func() {
	controlANew := "hi"
	controlAOld := "good"
	controlBNew := "there"
	controlBOld := "bye"

	Describe("unpack args function", func() {
		It("Should assign the value in the array to the associated pointer given", func() {
			internalA := controlAOld
			internalB := controlBOld
			arr := []interface{}{controlANew, controlBNew}
			err := Unpack(arr, &internalA, &internalB)
			Ω(err).Should(BeNil())
			Expect(internalA).NotTo(Equal(controlAOld))
			Expect(internalA).To(Equal(controlANew))
			Expect(internalB).NotTo(Equal(controlBOld))
			Expect(internalB).To(Equal(controlBNew))
		})

		It("Should return error if there the argument lengths dont match", func() {
			internalA := controlAOld
			arr := []interface{}{controlANew, controlBNew}
			err := Unpack(arr, &internalA)
			Ω(err).ShouldNot(BeNil())
		})

		It("Should return error if there the arguments of non matching types", func() {
			internalA := []string{"hi there"}
			arr := []interface{}{controlANew}
			err := Unpack(arr, &internalA)
			Ω(err).ShouldNot(BeNil())
		})

		It("Should not panic if called with incorrect arg count", func() {
			internalA := controlAOld
			arr := []interface{}{controlANew, controlBNew}
			Ω(func() { Unpack(arr, &internalA) }).ShouldNot(Panic())
		})

		It("Should not panic if called with invalid arg types", func() {
			internalA := []string{"hi there"}
			arr := []interface{}{controlANew}
			Ω(func() { Unpack(arr, &internalA) }).ShouldNot(Panic())
		})

		It("Should not panic if called with empty", func() {
			arr := []interface{}{controlANew}
			Ω(func() { Unpack(arr, Empty()) }).ShouldNot(Panic())
		})

		It("Should not error if called with empty", func() {
			arr := []interface{}{controlANew}
			err := Unpack(arr, Empty())
			Ω(err).Should(BeNil())
		})

	})
})
