package database

import "io"

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

func (db Context) Exec(writer io.Writer, getCommand Command) error {
	command := getCommand(db, writer)
	return command.Run()
}
