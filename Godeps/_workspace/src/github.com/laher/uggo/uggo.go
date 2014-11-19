package uggo

import (
	"strings"
)

const (
	VERSION = "0.3.2"
)

// gnuify a list of flags
// 'Gnuification' in this case refers to the Gnu-like expansion of single-hyphen multi-letter flags such as `-la` into separate flags `-l -a`
// note: allow '-help' to be used as single-hyphen (to assist the unitiated)
func Gnuify(call []string) []string {
	return GnuifyWithExceptions(call, []string{"-help"})
}

// simple slice helper
func contains(slice []string, subject string) bool {
	for _, item := range slice {
		if item == subject {
			return true
		}
	}
	return false
}

// 'Gnuify' a slice of flags, all bar for a list of exceptions.
// 'Gnuification' in this case refers to the Gnu-like expansion of single-hyphen multi-letter flags such as `-la` into separate flags `-l -a`
// 'Execeptions' in this case means exceptions to the rule (not the language construct)
func GnuifyWithExceptions(call, exceptions []string) []string {
	splut := []string{}
	for _, item := range call {
		if strings.HasPrefix(item, "-") && !strings.HasPrefix(item, "--") && !contains(exceptions, item) {
			for _, letter := range item[1:] {
				splut = append(splut, "-"+string(letter))
			}
		} else {
			splut = append(splut, item)
		}
	}
	return splut
}
