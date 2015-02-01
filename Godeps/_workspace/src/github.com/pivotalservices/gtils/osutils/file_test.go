package osutils_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	. "github.com/pivotalservices/gtils/osutils"
	"github.com/pkg/sftp"

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
		Context("remote file", func() {
			var client *sftpClientMock
			var fsState string = path.Join("some", "path")

			Context("check parent directories failures", func() {

				Context("Mkdir call returns error", func() {
					BeforeEach(func() {
						client = &sftpClientMock{
							FileSystemState: fsState,
							MkdirError:      true,
						}
					})

					It("should not create the directories and return non nil error", func() {
						newState := path.Join(fsState, "one", "two", "three")
						_, err := SafeCreateSSH(client, path.Join(newState, "filename.txt"))
						Ω(err).ShouldNot(BeNil())
						Ω(client.FileSystemState).ShouldNot(Equal(newState))
					})
				})

				Context("Create call returns error", func() {
					BeforeEach(func() {
						client = &sftpClientMock{
							FileSystemState: fsState,
							CreateError:     true,
						}
					})

					It("should create directories if needed but return non nil error", func() {
						newState := path.Join(fsState, "one", "two", "three")
						_, err := SafeCreateSSH(client, path.Join(newState, "filename.txt"))
						Ω(err).ShouldNot(BeNil())
						Ω(client.FileSystemState).Should(Equal(newState))
					})
				})

			})

			Context("check parent directories success", func() {
				BeforeEach(func() {
					client = &sftpClientMock{
						FileSystemState: fsState,
					}
				})

				Context("x-level parent directory does not exist", func() {
					It("should create the directories and nil error", func() {
						newState := path.Join(fsState, "one", "two", "three")
						_, err := SafeCreateSSH(client, path.Join(newState, "filename.txt"))
						Ω(err).Should(BeNil())
						Ω(client.FileSystemState).Should(Equal(newState))
					})
				})

				Context("1st-level parent directory does not exist", func() {
					It("should create the directories and nil error", func() {
						newState := path.Join(fsState, "one")
						_, err := SafeCreateSSH(client, path.Join(newState, "filename.txt"))
						Ω(err).Should(BeNil())
						Ω(client.FileSystemState).Should(Equal(newState))
					})
				})

				Context("parent directory already exists", func() {
					It("should not create any directories and nil error", func() {
						_, err := SafeCreateSSH(client, path.Join(fsState, "filename.txt"))
						Ω(err).Should(BeNil())
						Ω(client.FileSystemState).Should(Equal(fsState))
					})
				})
			})
		})

		Context("local file", func() {
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
})

type sftpClientMock struct {
	CreateError     bool
	MkdirError      bool
	FileSystemState string
}

func (s *sftpClientMock) Create(path string) (f *sftp.File, err error) {

	if s.CreateError {
		err = fmt.Errorf("create file error")

	} else {
		f = new(sftp.File)
	}
	return
}

func (s *sftpClientMock) Mkdir(p string) (err error) {

	if s.MkdirError {
		err = fmt.Errorf("mkdir failed")

	} else {
		s.FileSystemState = path.Join(p)
	}
	return
}

func (s *sftpClientMock) ReadDir(p string) (fi []os.FileInfo, err error) {
	if !strings.HasPrefix(s.FileSystemState, p) {
		err = fmt.Errorf("dir doesnt exist")
	}
	return
}
