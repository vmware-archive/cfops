package itertools

import (
	"testing"
)

var (
	firstPairTest   int
	secondPairTest  int
	controlFirst    int
	controlSecond   int
	controlPairTest Pair
)

func SetupPair() {
	firstPairTest = 0
	secondPairTest = 0
	controlFirst = 5
	controlSecond = 2
	controlPairTest = Pair{controlFirst, controlSecond}
}

func TearDownPair() {
	SetupPair()
}

func Test_PairUnPack(t *testing.T) {
	SetupPair()
	defer TearDownPair()

	PairUnPack(controlPairTest, &firstPairTest, &secondPairTest)

	if firstPairTest != controlFirst {
		t.Errorf("%d should is not equal to %d", controlFirst, firstPairTest)
	}

	if secondPairTest != controlSecond {
		t.Errorf("%d should is not equal to %d", controlSecond, secondPairTest)
	}
}
