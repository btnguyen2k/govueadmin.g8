package bo

import (
	"github.com/btnguyen2k/prom"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// NewSqliteConnection creates a new connection pool to SQLite3.
func NewSqliteConnection(dir, dbName string) *prom.SqlConnect {
	err := os.MkdirAll(dir, 0711)
	if err != nil {
		panic(err)
	}
	sqlc, err := prom.NewSqlConnect("sqlite3", dir+"/"+dbName+".db", 10000, nil)
	if err != nil {
		panic(err)
	}
	return sqlc
}
