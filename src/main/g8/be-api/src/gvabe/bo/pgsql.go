package bo

import (
	"github.com/btnguyen2k/prom"
	_ "github.com/lib/pq"
	"time"
)

// NewPgsqlConnection creates a new connection pool to PostgreSQL.
func NewPgsqlConnection(url, timezone string) *prom.SqlConnect {
	driver := "postgres"
	sqlConnect, err := prom.NewSqlConnect(driver, url, 10000, nil)
	if err != nil {
		panic(err)
	}
	loc, _ := time.LoadLocation(timezone)
	sqlConnect.SetLocation(loc)
	return sqlConnect
}
