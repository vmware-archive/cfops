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

type mockClient struct {
	session SSHSession
}

func (c *mockClient) NewSession() (SSHSession, error) {
	return c.session, nil
}

type mockSession struct {
	StartSuccess  bool
	StdOutSuccess bool
	WaitSuccess   bool
	CloseSuccess  bool
}

func (session *mockSession) Start(command string) (err error) {
	if !session.StartSuccess {
		err = errors.New("")
	}
	return
}

func (session *mockSession) Close() (err error) {
	if !session.CloseSuccess {
		err = errors.New("")
	}
	return
}
func (session *mockSession) Wait() (err error) {
	if !session.WaitSuccess {
		err = errors.New("")
	}
	return
}

func (session *mockSession) StdoutPipe() (reader io.Reader, err error) {
	if !session.StdOutSuccess {
		err = errors.New("")
		return nil, err
	}
	reader = strings.NewReader("mocksession")
	return
}

var _ = Describe("Ssh", func() {
	var (
		session *mockSession
		client  *mockClient
	)

	BeforeEach(func() {
		session = &mockSession{StartSuccess: true,
			StdOutSuccess: true,
			WaitSuccess:   true,
			CloseSuccess:  true}
		client = &mockClient{session: session}

	})

	Describe("Session Run success", func() {
		Context("Everything is fine", func() {
			It("should write to the writer", func() {
				var writer bytes.Buffer
				copier := NewCopier(client)
				err := copier.Execute(&writer, "command")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(writer.String()).Should(Equal("mocksession"))
			})
		})

	})
	Describe("Session Run failed", func() {

		Context("With bad stdpipeline", func() {
			It("should fail to write to the writer", func() {
				var writer bytes.Buffer
				copier := NewCopier(client)
				session.StdOutSuccess = false
				err := copier.Execute(&writer, "command")
				Ω(err).Should(HaveOccurred())
			})
		})
		Context("With bad command start", func() {
			It("should fail to write to the writer", func() {
				var writer bytes.Buffer
				copier := NewCopier(client)
				session.StartSuccess = false
				err := copier.Execute(&writer, "command")
				Ω(err).Should(HaveOccurred())
			})
		})
	})

})
