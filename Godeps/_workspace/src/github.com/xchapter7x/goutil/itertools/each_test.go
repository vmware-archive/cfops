package itertools

import (
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Filter", func() {
	var iterable []string = []string{
		"hi",
		"there",
		"who are",
		"you",
		"that thing",
	}
	var output int

	Context("called w/ 2 typed argument signature", func() {
		BeforeEach(func() {
			output = 0
			Each(iterable, func(v string) bool {
				output++
				return strings.HasPrefix(v, "t")
			})
		})

		It("should contain the correct number of elements", func() {
			Ω(output).Should(Equal(len(iterable)))
		})
	})

	Context("called with a non-convertable type argument func signature", func() {
		It("should not convert the type and just panic", func() {
			Ω(func() {
				Each(iterable, func(v int) {})
			}).Should(Panic())
		})
	})

	Context("called with a convertable type argument func signature", func() {
		It("should convert the type and not panic", func() {
			Ω(func() {
				Each(iterable, func(_, v string) {})
			}).ShouldNot(Panic())
		})
	})

	Context("called w/ 2 convertable typed argument signature", func() {
		var result bool

		BeforeEach(func() {
			output = 0
			Each(iterable, func(_, v string) bool {
				output++
				result = strings.HasPrefix(v, "t")
				return result
			})
		})

		It("should contain the correct number of elements", func() {
			Ω(output).Should(Equal(len(iterable)))
			Ω(result).Should(BeTrue())
		})
	})

	Context("called w/ 1 typed argument signature", func() {
		var result bool

		BeforeEach(func() {
			output = 0
			Each(iterable, func(v string) bool {
				output++
				result = strings.HasPrefix(v, "t")
				return result
			})
		})

		It("should contain the correct number of elements", func() {
			Ω(output).Should(Equal(len(iterable)))
			Ω(result).Should(BeTrue())
		})
	})
})

var (
	f_called_each int
	s_each        []string
	m_each        map[string]string
)

func SetupEach() {
	f_called_each = 0
	s_each = []string{"asdf", "asdfasdf", "geeeg", "gggggggg"}
	m_each = map[string]string{"a": "asdf", "b": "asdfasdf", "c": "geeeg", "d": "gggggggg"}
}

func TearDownEach() {
	SetupEach()
}

func Test_EachSliceArray(t *testing.T) {
	SetupEach()
	defer TearDownEach()

	Each(s, func(i int, v string) string {
		f_called += 1
		return v
	})

	if f_called != len(s) {
		t.Errorf("func f was not called %d times", len(s))
	}
}

func Test_EachEach(t *testing.T) {
	SetupEach()
	defer TearDownEach()

	Each(m, func(i, v string) string {
		f_called += 1
		return v
	})

	if f_called != len(m) {
		t.Errorf("func mf was not called %d times", len(m))
	}
}

func Test_CEachSliceArray(t *testing.T) {
	SetupEach()
	defer TearDownEach()

	CEach(s, func(i int, v string) string {
		f_called += 1
		return v
	})

	if f_called != len(s) {
		t.Errorf("func f was not called %d times", len(s))
	}
}

func Test_CEachEach(t *testing.T) {
	SetupEach()
	defer TearDownEach()

	CEach(m, func(i, v string) string {
		f_called += 1
		return v
	})

	if f_called != len(m) {
		t.Errorf("func mf was not called %d times", len(m))
	}
}
