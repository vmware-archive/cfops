package itertools

import (
	"container/list"
	"container/ring"
	"strings"
	"testing"
)

func Test_IterateMap(t *testing.T) {
	controlSlice := make(map[int]string)
	for i, v := range strings.Split("this is a test", "") {
		controlSlice[i] = v
	}
	controlIndex := 0

	for i := range Iterate(controlSlice) {
		var first int
		var second string
		PairUnPack(i, &first, &second)

		if second != controlSlice[first] {
			t.Errorf("Index of iterate %d is not equal to control index of %d", second, controlSlice[first])
		}

		controlIndex++
	}

	if controlIndex != len(controlSlice) {
		t.Errorf("Number of iteratesion %d should be equal to string length %d", controlIndex, len(controlSlice))
	}
}

func Test_IterateRing(t *testing.T) {
	controlLow := 1
	controlHigh := 10
	r := ring.New(controlHigh)
	z := controlLow
	r.Value = z

	for p := r.Next(); p != r; p = p.Next() {
		z++
		p.Value = z
	}

	for p := controlLow; p != controlHigh+1; p++ {
		if r.Value != p {
			t.Errorf("the value of this ring record is %d but should be %d", r.Value, p)
		}
		r = r.Next()
	}
}

func Test_IterateArray(t *testing.T) {
	controlSlice := strings.Split("this is a test", "")
	controlIndex := 0

	for i := range Iterate(controlSlice) {

		if i.First != controlIndex {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.First, controlIndex)
		}

		if i.Second != controlSlice[controlIndex] {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.Second, controlSlice[controlIndex])
		}
		controlIndex++
	}

	if controlIndex != len(controlSlice) {
		t.Errorf("Number of iteratesion %d should be equal to string length %d", controlIndex, len(controlSlice))
	}
}

func Test_IterateChan(t *testing.T) {
	l := make(chan int)
	controlIndex := 0
	controlValue := 6

	go func(l chan int) {
		for v := 6; v > 0; v-- {
			l <- v
		}
	}(l)

	ci := controlIndex
	cv := controlValue

	for i := range Iterate(l) {

		if i.First != ci {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.First, ci)
		}

		if i.Second != cv {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.Second, cv)
		}
		ci++
		cv--

		if ci == controlValue {
			close(l)
		}
	}

	if ci != controlValue {
		t.Errorf("Number of iteratesion %d should be equal to string length %d", ci, controlValue)
	}
}

func Test_IterateString(t *testing.T) {
	controlString := "this is a test"
	controlSlice := strings.Split(controlString, "")
	controlIndex := 0

	for i := range Iterate(controlString) {

		if i.First != controlIndex {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.First, controlIndex)
		}

		if i.Second != controlSlice[controlIndex] {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.Second, controlSlice[controlIndex])
		}
		controlIndex++
	}

	if controlIndex != len(controlString) {
		t.Errorf("Number of iteratesion %d should be equal to string length %d", controlIndex, len(controlString))
	}
}

func Test_IterateList(t *testing.T) {
	l := list.New()
	controlIndex := 0
	controlValue := 6

	for i := 1; i <= controlValue; i++ {
		l.PushFront(i)
	}

	for i := range Iterate(l) {

		if i.First != controlIndex {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.First, controlIndex)
		}

		if i.Second != controlValue {
			t.Errorf("Index of iterate %d is not equal to control index of %d", i.Second, controlValue)
		}
		controlIndex++
		controlValue--
	}

	if controlIndex != l.Len() {
		t.Errorf("Number of iteratesion %d should be equal to string length %d", controlIndex, l.Len())
	}
}
