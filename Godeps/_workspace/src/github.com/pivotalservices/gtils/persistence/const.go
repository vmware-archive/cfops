package persistence

const (
	//PGDmpDropCmd - command to drop a psql schema
	PGDmpDropCmd = "drop schema public cascade;"
	//PGDmpCreateCmd - create schema psql command
	PGDmpCreateCmd = "create schema public;"

	//MySQLDmpConnectCmd --
	MySQLDmpConnectCmd = "%s -u %s -h %s --password=%s"
	//MySQLDmpCreateCmd --
	MySQLDmpCreateCmd = "%s < %s"
	//MySQLDmpFlushCmd --
	MySQLDmpFlushCmd = "%s > flush privileges"
	//MySQLDmpDumpCmd --
	MySQLDmpDumpCmd = "%s --all-databases"
)

var (
	//PGDmpDumpBin - location of pg_dump binary
	PGDmpDumpBin = "/var/vcap/packages/postgres/bin/pg_dump"
	//PGDmpRestoreBin - location of pg_restore binary
	PGDmpRestoreBin = "/var/vcap/packages/postgres/bin/pg_restore"

	//MySQLDmpDumpBin --
	MySQLDmpDumpBin = "/var/vcap/packages/mariadb/bin/mysqldump"
	//MySQLDmpSQLBin --
	MySQLDmpSQLBin = "/var/vcap/packages/mariadb/bin/mysql"
)
