package models_test

import (
	"github.com/cloudfoundry-community/gogobosh/models"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DeploymentManifest", func() {
	It("FindByJobTemplates", func() {
		manifest := &models.DeploymentManifest{
			Jobs: []*models.ManifestJob{
				{Name: "job1", JobTemplates: []*models.ManifestJobTemplate{{Name: "common"}}},
				{Name: "job2", JobTemplates: []*models.ManifestJobTemplate{{Name: "common"}}},
				{Name: "other", JobTemplates: []*models.ManifestJobTemplate{{Name: "other"}}},
			},
		}
		jobs := manifest.FindByJobTemplates("common")
		Expect(len(jobs)).To(Equal(2))
	})
})
