package command_test

import (
	"bytes"
	"errors"
	"io"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/command"
)

type mockSession struct {
	runSuccess    bool
	StdOutSuccess bool
}

func (session mockSession) Run(command string) (err error) {
	if !session.runSuccess {
		err = errors.New("")
	}
	return
}

func (session mockSession) StdoutPipe() (reader io.Reader, err error) {
	if !session.StdOutSuccess {
		err = errors.New("")
		return nil, err
	}
	reader = strings.NewReader("mocksession")
	return
}

var _ = Describe("Ssh", func() {
	var (
		session mockSession
	)

	BeforeEach(func() {
		session = mockSession{}
	})

	Describe("Session Run success", func() {
		Context("With Good stdpipeline", func() {
			It("should write to the writer", func() {
				var writer bytes.Buffer
				session.runSuccess = true
				session.StdOutSuccess = true
				copier := NewCopier(session)
				err := copier.Execute(&writer, "command")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(writer.String()).Should(Equal("mocksession"))
			})
		})

		Context("With bad stdpipeline", func() {
			It("should failed to write to the writer", func() {
				var writer bytes.Buffer
				session.runSuccess = true
				session.StdOutSuccess = false
				copier := NewCopier(session)
				err := copier.Execute(&writer, "command")
				Ω(err).Should(HaveOccurred())
			})
		})
	})
	Describe("Session Run failed", func() {
		Context("With Good stdpipeline", func() {
			It("should write to the writer", func() {
				var writer bytes.Buffer
				session.runSuccess = false
				session.StdOutSuccess = true
				copier := NewCopier(session)
				err := copier.Execute(&writer, "command")
				Ω(err).Should(HaveOccurred())
			})
		})
	})

})
