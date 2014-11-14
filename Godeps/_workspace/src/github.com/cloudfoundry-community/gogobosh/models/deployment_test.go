package models_test

import (
	"github.com/cloudfoundry-community/gogobosh/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deployments", func() {
	It("FindByRelease", func() {
		manifest := &models.Deployments{
			{
				Name: "cf-warden",
				Releases: []models.NameVersion{
					{Name: "nagios"},
					{Name: "cf"},
				},
			},
			{
				Name: "other",
				Releases: []models.NameVersion{
					{Name: "nagios"},
					{Name: "other"},
				},
			},
		}
		deployments := manifest.FindByRelease("cf")
		Expect(len(deployments)).To(Equal(1))
	})
})
