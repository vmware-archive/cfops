package database

import (
	"io"
	"os"
	"os/exec"
	"strconv"
)

const (
	PGPASSWORD = "PGPASSWORD"
	PGCOMMAND  = "pg_dump"
)

type CommandInterface interface {
	Run() error
}

type Command func(Context, io.Writer) CommandInterface

type Context struct {
	username string
	password string
	host     string
	port     int
	db       string
}

func NewContext(username, password, db, host string, port int) Context {
	return Context{
		username: username,
		password: password,
		host:     host,
		port:     port,
		db:       db,
	}
}

type Database struct {
	context Context
	backup  Command
}

func NewDatabase(context Context, fn Command) Database {
	return Database{
		context: context,
		backup:  fn,
	}
}

// Need integration tests on the postgres backup command

var PostGresBackupCommand = func(context Context, writer io.Writer) CommandInterface {
	os.Setenv(PGPASSWORD, context.password)
	cmd := exec.Command(PGCOMMAND, "-h", context.host, "-U", context.username, "-p", strconv.Itoa(context.port), context.db)
	cmd.Stdout = writer
	cmd.Stderr = writer
	return cmd
}

func (db *Database) Backup(writer io.Writer) error {
	command := db.backup(db.context, writer)
	return command.Run()
}
