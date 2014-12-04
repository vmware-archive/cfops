package localengine_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/toggle/engines/localengine"
)

var controlSuccessStatus string = "true"

func successGetenvMock(fs string) (status string) {
	status = controlSuccessStatus
	return
}

func failureGetenvMock(fs string) (status string) {
	status = ""
	return
}

var _ = Describe("localengine package", func() {
	Describe("LocalEngine struct", func() {
		Describe("GetFeatureStatusValue function", func() {
			It("Should return the result of getenv and have nil error on success", func() {
				engine := &localengine.LocalEngine{
					Getenv: successGetenvMock,
				}
				res, err := engine.GetFeatureStatusValue("")
				Expect(res).To(Equal(controlSuccessStatus))
				Ω(err).Should(BeNil())
			})

			It("Should return non nil err on failed call", func() {
				engine := &localengine.LocalEngine{
					Getenv: failureGetenvMock,
				}
				_, err := engine.GetFeatureStatusValue("")
				Ω(err).ShouldNot(BeNil())
			})
		})
	})
})
