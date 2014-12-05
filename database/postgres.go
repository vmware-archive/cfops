package database

import (
	"io"
	"os"
	"os/exec"
	"strconv"
)

var PostGresBackupCommand = func(context Context, writer io.Writer) CommandInterface {
	os.Setenv(PGPASSWORD, context.password)
	cmd := exec.Command(PGCOMMAND, "-h", context.host, "-U", context.username, "-p", strconv.Itoa(context.port), context.db)
	cmd.Stdout = writer
	cmd.Stderr = writer
	return cmd
}
