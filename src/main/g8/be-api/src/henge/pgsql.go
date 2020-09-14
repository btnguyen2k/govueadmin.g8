package henge

import (
	"time"

	"github.com/btnguyen2k/prom"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// NewPgsqlConnection creates a new connection pool to PostgreSQL.
func NewPgsqlConnection(url, timezone string) *prom.SqlConnect {
	driver := "pgx"
	sqlConnect, err := prom.NewSqlConnect(driver, url, 10000, nil)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation(timezone)
	sqlConnect.SetLocation(loc).SetDbFlavor(prom.FlavorPgSql)
	return sqlConnect
}

// InitPgsqlTable initializes database table to store bo
func InitPgsqlTable(sqlc *prom.SqlConnect, tableName string, extraCols map[string]string) {
	colDef := map[string]string{
		ColId:          "VARCHAR(64)",
		ColData:        "JSONB",
		ColChecksum:    "VARCHAR(32)",
		ColTimeCreated: "TIMESTAMP WITH TIME ZONE",
		ColTimeUpdated: "TIMESTAMP WITH TIME ZONE",
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
