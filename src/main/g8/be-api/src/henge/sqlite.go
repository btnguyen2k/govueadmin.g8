package henge

import (
	"os"

	"github.com/btnguyen2k/prom"
	_ "github.com/mattn/go-sqlite3"
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
	sqlc.SetDbFlavor(prom.FlavorDefault)
	return sqlc
}

// InitSqliteTable initializes database table to store bo
func InitSqliteTable(sqlc *prom.SqlConnect, tableName string, extraCols map[string]string) {
	colDef := map[string]string{
		ColId:          "VARCHAR(64)",
		ColData:        "VARCHAR(255)",
		ColChecksum:    "VARCHAR(32)",
		ColTimeCreated: "TIMESTAMP",
		ColTimeUpdated: "TIMESTAMP",
		ColAppVersion:  "BIGINT",
	}
	colNames := []string{ColId, ColData, ColChecksum, ColTimeCreated, ColTimeUpdated, ColAppVersion}
	for k, v := range extraCols {
		colDef[k] = v
		colNames = append(colNames, k)
	}
	pk := []string{ColId}
	if err := CreateTable(sqlc, tableName, true, colDef, colNames, pk); err != nil {
		panic(err)
	}
}
