package backup

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBackupInternals(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Backup Suite (Internal)")
}
