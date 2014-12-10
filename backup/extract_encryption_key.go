package backup

import (
	"fmt"
	"os"
	"path"

	"github.com/pivotalservices/cfops/command"
	"github.com/pivotalservices/cfops/osutils"
)

var ExtractEncryptionKey = func(backupDir, deploymentDir string, exec command.CmdExecuter) (err error) {
	backupFileName := path.Join(backupDir, "cc_db_encryption_key.txt")
	os.Remove(backupFileName)
	b, err := osutils.SafeCreate(backupFileName)
	defer b.Close()

	if err == nil {
		formatString := `grep -E 'db_encryption_key' %s/cf-*.yml | cut -d ':' -f 2 | sort -u | tr -d ' '`
		cmd := fmt.Sprintf(formatString, deploymentDir)
		err = exec.Execute(b, cmd)
	}
	return
}
