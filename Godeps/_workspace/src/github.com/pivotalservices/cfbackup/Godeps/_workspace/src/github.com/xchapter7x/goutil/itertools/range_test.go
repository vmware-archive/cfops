package itertools

import (
	"testing"
)

var (
	controlRangeCount int
)

func SetupRange() {
	controlRangeCount = 5
}

func TearDownRange() {
	SetupRange()
}

func Test_Range(t *testing.T) {
	SetupRange()
	defer TearDownRange()

	testCount := 0
	for i := range Range(1, controlRangeCount) {
		testCount++

		if testCount != i {
			t.Errorf("Range returned %d but was expecting %d", i, testCount)
		}
	}

	if testCount != controlRangeCount {
		t.Errorf("Range was not called %d times. it was called %d", controlRangeCount, testCount)
	}
}
