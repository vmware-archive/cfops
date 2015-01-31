package osutils_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/gtils/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("File", func() {
	var (
		dir  string
		file string
	)

	BeforeEach(func() {
		dir, _ = ioutil.TempDir("", "cfops-osutils")
		os.MkdirAll(dir, 0777)
	})

	AfterEach(func() {
		os.RemoveAll(dir)
	})

	Describe("File creation", func() {
		Context("With a path that does not exist and has missing directories", func() {
			BeforeEach(func() {
				file = path.Join(dir, "a", "nonexistent", "directory", "with", "af.ile")
			})

			It("should create the parent directory", func() {
				base := path.Join(dir, "a", "nonexistent", "directory", "with")
				e, _ := Exists(base)
				Ω(e).To(BeFalse())
				_, err := SafeCreate(file)
				Ω(err).To(BeNil())
				e, _ = Exists(base)
				Ω(e).To(BeTrue())
			})
		})

		Context("With a path that does not exist and has missing directories", func() {
			BeforeEach(func() {
				file = path.Join(dir, "af.ile")
			})

			It("should create the parent directory", func() {
				base := path.Join(dir)
				e, _ := Exists(base)
				Ω(e).To(BeTrue())
				_, err := SafeCreate(file)
				Ω(err).To(BeNil())
				e, _ = Exists(base)
				Ω(e).To(BeTrue())
			})
		})
	})
})
