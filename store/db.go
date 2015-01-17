package store

import (
	"os"
	"sync"

	"github.com/cryptix/go/logging"
	"github.com/jmoiron/modl"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is the global database.
var DB = &modl.DbMap{Dialect: modl.PostgresDialect{}}

// DBH is a modl.SqlExecutor interface to DB, the global database. It is better
// to use DBH instead of DB because it prevents you from calling methods that
// could not later be wrapped in a transaction.
var DBH modl.SqlExecutor = DB

var connectOnce sync.Once

// Connect connects to the PostgreSQL database specified by the PG* environment
// variables. It calls logging.CheckFatal if it encounters an error.
func Connect() {
	connectOnce.Do(func() {
		var err error
		DB.Dbx, err = sqlx.Open("postgres", os.Getenv("HEROKU_POSTGRESQL_PURPLE_URL"))
		logging.CheckFatal(err)
		DB.Db = DB.Dbx.DB
	})
}

// Create the database schema. It calls log.Fatal if it encounters an error.
func Create() {
	err := DB.CreateTablesIfNotExists()
	logging.CheckFatal(err)
}

// Drop the database schema.
func Drop() (err error) {
	_, err = DB.Exec(`DROP TABLE IF EXISTS "track";DROP TABLE IF EXISTS "user";`)
	return
}
