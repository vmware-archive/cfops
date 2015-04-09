package itertools_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/xchapter7x/goutil/itertools"
)

var _ = Describe("Find", func() {
	var iterable []string = []string{
		"hi",
		"there",
		"who are",
		"you",
		"that thing",
	}
	var output Pair

	Context("traversing iterator w/ multiple matches", func() {

		BeforeEach(func() {
			output = Find(iterable, func(pair Pair) bool {
				val := pair.Second.(string)
				return strings.HasPrefix(val, "t")
			})
		})

		It("Should return only the first match", func() {
			Ω(output.Second).Should(Equal("there"))
		})
	})
})
