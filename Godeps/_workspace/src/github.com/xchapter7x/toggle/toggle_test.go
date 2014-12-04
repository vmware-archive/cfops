package toggle_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/xchapter7x/toggle"
)

var _ = Describe("toggle package", func() {
	controlNamespace := "hi"

	BeforeEach(func() {
		toggle.Init(controlNamespace, nil)
	})

	AfterEach(func() {
		toggle.Close()
	})

	Describe("RegisterFeature Function", func() {
		It("Should inject a new feature and return nil error", func() {
			initialFeatureCount := len(toggle.ShowFeatures())
			featureName := "sampleFeature"
			err := toggle.RegisterFeature(featureName)
			currentFeatureCount := len(toggle.ShowFeatures())
			Expect(initialFeatureCount).NotTo(Equal(currentFeatureCount))
			Ω(err).Should(BeNil())
		})

		It("Should add feature record for referencing", func() {
			initialFeatureList := toggle.ShowFeatures()
			featureName := "sampleFeature"
			_, controlExists := initialFeatureList[featureName]
			toggle.RegisterFeature(featureName)
			currentFeatureList := toggle.ShowFeatures()
			_, currentExists := currentFeatureList[toggle.GetFullFeatureSignature(featureName)]
			Expect(controlExists).NotTo(Equal(currentExists))
		})

		It("Should ignore duplicate register calls and return non nil error", func() {
			featureName := "sampleFeature"
			toggle.RegisterFeature(featureName)
			initialFeatureCount := len(toggle.ShowFeatures())
			err := toggle.RegisterFeature(featureName)
			currentFeatureCount := len(toggle.ShowFeatures())
			Expect(initialFeatureCount).To(Equal(currentFeatureCount))
			Ω(err).ShouldNot(BeNil())
		})

	})

	Describe("IsActive function", func() {
		flagName := "bogusFlag"

		It("Should return false if given unregistered flag", func() {
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(false))
		})

		It("Should return false if given flag that is FEATURE_OFF status ", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_OFF)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(false))
		})

		It("Should return true if given flag that is FEATURE_ON status ", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(true))
		})
	})

	Describe("SetFeatureStatus function", func() {
		flagName := "bogusFlag"

		It("Should return false if setting FEATURE_OFF status from default", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_OFF)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(false))
		})

		It("Should return true if setting FEATURE_ON status from default", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(true))
		})

		It("Should return false if setting FEATURE_OFF status updating existing value", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_OFF)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(false))
		})

		It("Should return true if setting FEATURE_ON status updating existing value", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_OFF)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			response := toggle.IsActive(flagName)
			Expect(response).To(Equal(true))
		})

		It("Should return a non nil error if flagName not valid", func() {
			response := toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			Ω(response).ShouldNot(BeNil())
		})

		It("Should return a nil error if flagName is valid", func() {
			toggle.RegisterFeature(flagName)
			response := toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			Ω(response).Should(BeNil())
		})

	})

	Describe("Flip function", func() {
		flagName := "bogusFlag"
		controlDefault := "default"
		controlNew := "new"

		It("Should take arguments and pass them into the selected function", func() {
			toggle.RegisterFeature(flagName)
			status := ""
			toggle.Flip(flagName, func(arg1 string) {
				status = arg1
			}, func(arg1 string) {
				status = arg1
			}, controlDefault)
			Expect(status).To(Equal(controlDefault))
		})

		It("Should return arguments from the flip function in the form of an interface array", func() {
			toggle.RegisterFeature(flagName)
			status := ""
			var response []interface{}
			argumentCount := 1
			controlLen := len(response) + argumentCount
			response = toggle.Flip(flagName, func() (r string) {
				status = controlDefault
				return status
			}, func() (r string) {
				status = controlNew
				return status
			})
			Expect(len(response)).To(Equal(controlLen))
			Expect(response[0]).To(Equal(status))
		})

		It("Should return arguments from the flip function in the form of an interface array when multiple response values", func() {
			toggle.RegisterFeature(flagName)
			status := ""
			var response []interface{}
			argumentCount := 2
			controlLen := len(response) + argumentCount
			response = toggle.Flip(flagName, func() (r, j string) {
				status = controlDefault
				return status, status
			}, func() (r, j string) {
				status = controlNew
				return status, status
			})
			Expect(len(response)).To(Equal(controlLen))
			Expect(response[0]).To(Equal(status))
			Expect(response[1]).To(Equal(status))
		})

		It("Should select the default feature function to run when flag is default", func() {
			toggle.RegisterFeature(flagName)
			status := ""
			toggle.Flip(flagName, func() {
				status = controlDefault
			}, func() {
				status = controlNew
			})
			Expect(status).To(Equal(controlDefault))
		})

		It("Should select the new feature function to run when flag is set to inactive", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_OFF)
			status := ""
			toggle.Flip(flagName, func() {
				status = controlDefault
			}, func() {
				status = controlNew
			})
			Expect(status).To(Equal(controlDefault))
		})

		It("Should select the new feature function to run when flag is set to active", func() {
			toggle.RegisterFeature(flagName)
			toggle.SetFeatureStatus(flagName, toggle.FEATURE_ON)
			status := ""
			toggle.Flip(flagName, func() {
				status = controlDefault
			}, func() {
				status = controlNew
			})
			Expect(status).To(Equal(controlNew))
		})
	})
})
