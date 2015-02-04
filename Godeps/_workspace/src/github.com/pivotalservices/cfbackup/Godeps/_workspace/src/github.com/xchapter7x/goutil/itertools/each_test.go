package itertools

import (
	"testing"
)

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
