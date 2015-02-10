package itertools

import (
	"testing"
)

var (
	controlFill    string
	controlSample1 string
	controlSample2 string
	controlSample3 string
	controlZipped  [][]string
)

func SetupZip() {
	controlFill = "-"
	controlSample1 = "abcdefghijk"
	controlSample2 = "ABCde"
	controlSample3 = "aBCDefg"
	controlZipped = [][]string{
		[]string{"a", "A", "a"},
		[]string{"b", "B", "B"},
		[]string{"c", "C", "C"},
		[]string{"d", "d", "D"},
		[]string{"e", "e", "e"},
		[]string{"f", "-", "f"},
		[]string{"g", "-", "g"},
		[]string{"h", "-", "-"},
		[]string{"i", "-", "-"},
		[]string{"j", "-", "-"},
		[]string{"k", "-", "-"}}
}

func TearDownZip() {
	controlFill = ""
	controlSample1 = ""
	controlSample2 = ""
	controlSample3 = ""
}

func Test_ZipLogest(t *testing.T) {
	SetupZip()
	defer TearDownZip()
	count := 0

	for z := range ZipLongest(controlFill, controlSample1, controlSample2, controlSample3) {

		for i := range controlZipped[count] {

			if controlZipped[count][i] != z[i].(string) {
				t.Errorf("Error: %s should match the control zipped dataset string w/ %s ", z[i].(string), controlZipped[count][i])
			}
		}
		count++
	}

	if count != len(controlSample1) {
		t.Errorf("Error: %d should match the longest dataset w/ %d ", count, len(controlSample1))
	}
}

func Test_Zip(t *testing.T) {
	SetupZip()
	defer TearDownZip()
	count := 0

	for z := range Zip(controlFill, controlSample1, controlSample2, controlSample3) {

		for i := range controlZipped[count] {

			if controlZipped[count][i] != z[i].(string) {
				t.Errorf("Error: %s should match the control zipped dataset string w/ %s ", z[i].(string), controlZipped[count][i])
			}
		}
		count++
	}

	if count != len(controlSample2) {
		t.Errorf("Error: %d should match the shortest dataset w/ %d ", count, len(controlSample1))
	}
}
