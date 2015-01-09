package goutil_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/goutil"
)

var _ = Describe("CallChain", func() {
	var callCount int = 0
	var sampleSuccessReturn string = "success"
	var sampleFailureReturn string = "failure"
	var controlError error = fmt.Errorf(sampleFailureReturn)
	BeforeEach(func() {
		callCount = 0
	})

	AfterEach(func() {
		callCount = 0
	})

	var successMultiReturn func(string) (string, error) = func(s string) (string, error) {
		callCount++
		return sampleSuccessReturn, nil
	}

	var failMultiReturn func(string) (string, error) = func(s string) (string, error) {
		callCount++
		return sampleFailureReturn, controlError
	}

	var successNoArgMultiReturn func() (string, error) = func() (string, error) {
		return sampleSuccessReturn, nil
	}

	var failNoArgMultiReturn func() (string, error) = func() (string, error) {
		return sampleFailureReturn, controlError
	}

	var runNoError func() string = func() string {
		return sampleSuccessReturn
	}

	Context("CallChainP function", func() {
		Context("with a nil chained error", func() {
			It("Should swap the values at the given pointers with the return values of the function and return nil error on sucess", func() {
				var stres string
				var errres error
				c := NewChain(nil)
				err := c.CallP(c.Returns(&stres, &errres), successMultiReturn, "random")
				Ω(sampleSuccessReturn).Should(Equal(stres))
				Ω(err).Should(BeNil())
			})

			It("Should not swap the values at the given pointers with the return values of the function and return non-nil error on failure", func() {
				var stres string
				var errres error
				response := []interface{}{
					&stres,
					&errres,
				}
				c := NewChain(nil)
				err := c.CallP(response, failMultiReturn, "random")
				Ω(sampleFailureReturn).Should(Equal(stres))
				Ω(err).ShouldNot(BeNil())
			})

		})
	})

	Context("CallChain function", func() {
		Context("with a non nil chained error", func() {
			It("Should return a error equal to the chained error", func() {
				e := fmt.Errorf("new error")
				c := NewChain(e)
				_, err := c.Call(failNoArgMultiReturn, "testing_error")
				Ω(err).ShouldNot(BeNil())
				Ω(err).Should(Equal(e))
			})

			It("Should return an error, skipping any call including failed calls, and return the original error", func() {
				e := fmt.Errorf("new error")
				c := NewChain(e)
				_, err := c.Call(successMultiReturn, "testing_error")
				Ω(err).ShouldNot(Equal(controlError))
			})

			It("Should skip the function if passed an error - even a success call function", func() {
				e := fmt.Errorf("new error")
				c := NewChain(e)
				c.Call(successMultiReturn, "testing_error")
				Ω(callCount).Should(Equal(0))
			})

			It("Should skip the function if passed an error - even a failed call function", func() {
				e := fmt.Errorf("new error")
				c := NewChain(e)
				c.Call(successMultiReturn, "testing_error")
				Ω(callCount).Should(Equal(0))
			})

			It("Should skip the function if passed an error - even a failed call function", func() {
				e := fmt.Errorf("new error")
				c := NewChain(e)
				t := "string"
				c.CallP([]interface{}{&t, &e}, successMultiReturn, "testing_error")
				Ω(callCount).Should(Equal(0))
				Ω(c.Error).ShouldNot(BeNil())
			})

		})

		Context("with a nil chained error", func() {
			Context("on success", func() {
				It("Should return a nil error w/ multiple return functions and arguments", func() {
					c := NewChain(nil)
					_, err := c.Call(successMultiReturn, "testing_error")
					Ω(err).Should(BeNil())
				})

				It("Should return a nil error w/ no arguments", func() {
					c := NewChain(nil)
					_, err := c.Call(successNoArgMultiReturn)
					Ω(err).Should(BeNil())
				})

				It("Should return a nil error w/ error values returned", func() {
					c := NewChain(nil)
					_, err := c.Call(runNoError)
					Ω(err).Should(BeNil())
				})
			})

			Context("on failure", func() {
				It("Should return a non nil error w/ multiple return values", func() {
					c := NewChain(nil)
					_, err := c.Call(failMultiReturn, "testing_error")
					Ω(err).ShouldNot(BeNil())
					Ω(err).Should(Equal(controlError))
					Ω(c.Error).ShouldNot(BeNil())
				})

				It("Should return a non nil error w/ no args", func() {
					c := NewChain(nil)
					_, err := c.Call(failNoArgMultiReturn)
					Ω(err).ShouldNot(BeNil())
					Ω(err).Should(Equal(controlError))
				})

				It("Should have the correct error status in the chain Error value", func() {
					var e error
					c := NewChain(nil)
					t := "string"
					c.CallP([]interface{}{&t, &e}, failMultiReturn, "testing_error")
					Ω(c.Error).ShouldNot(BeNil())
				})

				It("Should return proper value & error to given pointers", func() {
					var e error = nil
					c := NewChain(nil)
					t := "original value"
					err := c.CallP([]interface{}{&t, &e}, failMultiReturn, "testing_error")
					Ω(err).ShouldNot(BeNil())
					Ω(t).Should(Equal(sampleFailureReturn))
					Ω(e).ShouldNot(BeNil())
				})
			})
		})
	})
})
