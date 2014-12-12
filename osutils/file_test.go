package osutils_test

import (
	"io/ioutil"
	"os"
	"path"

	. "github.com/pivotalservices/cfops/osutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	content string = "test content"
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

	FDescribe("File creation", func() {
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

	FDescribe("File opening", func() {
		Context("With a path that does not exist and has missing directories", func() {
			BeforeEach(func() {
				file = path.Join(dir, "a", "nonexistent", "directory", "with", "af.ile")
			})

			It("should create the parent directory", func() {
				base := path.Join(dir, "a", "nonexistent", "directory", "with")
				e, _ := Exists(base)
				Ω(e).To(BeFalse())
				f, _ := SafeCreate(file)
				f.WriteString(content)

				rf, err := OpenFile(file)
				Ω(err).To(BeNil())
				e, _ = Exists(base)
				Ω(e).To(BeTrue())

				b1 := make([]byte, 20)
				rf.Read(b1)

				output := string(b1)

				Ω(output).Should(Equal(content))
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
				SafeCreate(file)
				_, err := OpenFile(file)
				Ω(err).To(BeNil())
				e, _ = Exists(base)
				Ω(e).To(BeTrue())
			})
		})

	})
})
