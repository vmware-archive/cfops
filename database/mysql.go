package database

import (
	"fmt"
	"io"
	"os/exec"
	"strconv"
)

const (
	MYSQLCOMMAND = "mysqldump"
)

func MysqlBackupCommand(context Context, writer io.Writer) CommandInterface {
	db := parseDB(context.db)
	cmd := exec.Command(MYSQLCOMMAND, "-h", context.host, "-U", context.username, "-p", context.password, "-P", strconv.Itoa(context.port), db)
	cmd.Stdout = writer
	cmd.Stderr = writer
	return cmd
}

func parseDB(db string) string {
	if db != "--all-databases" {
		return fmt.Sprintf("--databases %v", db)
	}
	return db
}
