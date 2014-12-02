package ssh_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/pivotalservices/cfops/ssh"
	"strings"
)

var _ = Describe("Ssh", func() {
	Describe("Dump from a reader to a writer", func() {
		Context("Dump succeeded with reader to writer", func() {
			It("dump from string to string", func() {
				r := strings.NewReader("teststring")
				var b bytes.Buffer
				dump := DumpToWriter{
					Writer: &b,
				}
				dump.Execute(r)
				Expect(b.String()).To(Equal("teststring"))
			})
		})
	})

})
