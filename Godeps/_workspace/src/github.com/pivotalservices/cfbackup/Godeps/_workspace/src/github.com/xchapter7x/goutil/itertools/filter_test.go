package itertools

import (
	"testing"
)

var filterTestData []string
var filterWhiteList []string

func SetupFilter() {
	filterTestData = []string{"asdf", "asdfasdf", "geeeg", "gggggggg"}
	filterWhiteList = []string{"asdfasdf", "geeeg"}
}

func TearDownFilter() {
	filterTestData = []string{}
	filterWhiteList = []string{}
}

func Test_CFilter(t *testing.T) {
	SetupFilter()
	defer TearDownFilter()

	f := CFilter(filterTestData, func(i, v interface{}) bool {
		return findInStringArray(v.(string), filterWhiteList)
	})

	for r := range f {

		if !findInStringArray(r.Second.(string), filterWhiteList) {
			t.Errorf("Error: %s should have been filtered, but it was not ", r)
		}
	}
}

func Test_Filter(t *testing.T) {
	SetupFilter()
	defer TearDownFilter()

	f := Filter(filterTestData, func(i, v interface{}) bool {
		return findInStringArray(v.(string), filterWhiteList)
	})

	for r := range f {

		if !findInStringArray(r.Second.(string), filterWhiteList) {
			t.Errorf("Error: %s should have been filtered, but it was not ", r)
		}
	}
}

func Test_CFilterFalse(t *testing.T) {
	SetupFilter()
	defer TearDownFilter()

	f := CFilterFalse(filterTestData, func(i, v interface{}) bool {
		return findInStringArray(v.(string), filterWhiteList)
	})

	for r := range f {

		if findInStringArray(r.Second.(string), filterWhiteList) {
			t.Errorf("Error: %s should have been filtered, but it was not ", r)
		}
	}
}

func Test_FilterFalse(t *testing.T) {
	SetupFilter()
	defer TearDownFilter()

	f := FilterFalse(filterTestData, func(i, v interface{}) bool {
		return findInStringArray(v.(string), filterWhiteList)
	})

	for r := range f {

		if findInStringArray(r.Second.(string), filterWhiteList) {
			t.Errorf("Error: %s should have been filtered, but it was not ", r)
		}
	}
}

func findInStringArray(v string, a []string) (r bool) {
	r = false

	for _, i := range a {

		if v == i {
			r = true
		}
	}
	return
}
