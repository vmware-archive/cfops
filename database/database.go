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

func (db *Database) Backup(writer io.Writer) error {
	command := db.backup(db.context, writer)
	return command.Run()
}
