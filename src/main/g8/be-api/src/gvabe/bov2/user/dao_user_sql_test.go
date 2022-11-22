package user

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	promsql "github.com/btnguyen2k/prom/sql"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testTimeZone = "Asia/Ho_Chi_Minh"
	testSqlTable = "test_user"
)

func sqlInitTable(sqlc *promsql.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	extraCols := map[string]string{UserColMaskUid: "VARCHAR(32)"}
	switch sqlc.GetDbFlavor() {
	case promsql.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId, Uk: [][]string{{"/" + UserColMaskUid}}}
		err = henge.InitCosmosdbCollection(sqlc, table, spec)
	case promsql.FlavorSqlite:
		err = henge.InitSqliteTable(sqlc, table, extraCols)
	case promsql.FlavorMySql:
		err = henge.InitMysqlTable(sqlc, table, extraCols)
	case promsql.FlavorPgSql:
		err = henge.InitPgsqlTable(sqlc, table, extraCols)
	}
	return err
}

func newSqlConnect(t *testing.T, testName string, driver, url, timezone string, flavor promsql.DbFlavor) (*promsql.SqlConnect, error) {
	driver = strings.Trim(driver, "\"")
	url = strings.Trim(url, "\"")
	if driver == "" || url == "" {
		t.Skipf("%s skipped", testName)
	}

	cosmosdb := cosmosdbDbName
	if flavor == promsql.FlavorCosmosDb {
		dbre := regexp.MustCompile(`(?i);db=(\w+)`)
		findResult := dbre.FindAllStringSubmatch(url, -1)
		if findResult == nil {
			url += ";Db=" + cosmosdb
		} else {
			cosmosdb = findResult[0][1]
		}
	}

	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	sqlc, err := promsql.NewSqlConnectWithFlavor(driver, url, 10000, nil, flavor)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}

	if err == nil && flavor == promsql.FlavorCosmosDb {
		sqlc.GetDB().Exec("CREATE DATABASE IF NOT EXISTS " + cosmosdb + " WITH maxru=10000")
	}

	return sqlc, err
}

func initDaoSql(sqlc *promsql.SqlConnect) UserDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return NewUserDaoCosmosdb(sqlc, testSqlTable, true)
	}
	return NewUserDaoSql(sqlc, testSqlTable, true)
}

const (
	envSqliteDriver = "SQLITE_DRIVER"
	envSqliteUrl    = "SQLITE_URL"
	envMssqlDriver  = "MSSQL_DRIVER"
	envMssqlUrl     = "MSSQL_URL"
	envMysqlDriver  = "MYSQL_DRIVER"
	envMysqlUrl     = "MYSQL_URL"
	envOracleDriver = "ORACLE_DRIVER"
	envOracleUrl    = "ORACLE_URL"
	envPgsqlDriver  = "PGSQL_DRIVER"
	envPgsqlUrl     = "PGSQL_URL"
)

type sqlDriverAndUrl struct {
	driver, url string
}

func newSqlDriverAndUrl(driver, url string) sqlDriverAndUrl {
	return sqlDriverAndUrl{driver: strings.Trim(driver, `"`), url: strings.Trim(url, `"`)}
}

func sqlGetUrlFromEnv() map[string]sqlDriverAndUrl {
	urlMap := make(map[string]sqlDriverAndUrl)
	if os.Getenv(envSqliteDriver) != "" && os.Getenv(envSqliteUrl) != "" {
		urlMap["sqlite"] = newSqlDriverAndUrl(os.Getenv(envSqliteDriver), os.Getenv(envSqliteUrl))
	}
	if os.Getenv(envMssqlDriver) != "" && os.Getenv(envMssqlUrl) != "" {
		urlMap["mssql"] = newSqlDriverAndUrl(os.Getenv(envMssqlDriver), os.Getenv(envMssqlUrl))
	}
	if os.Getenv(envMysqlDriver) != "" && os.Getenv(envMysqlUrl) != "" {
		urlMap["mysql"] = newSqlDriverAndUrl(os.Getenv(envMysqlDriver), os.Getenv(envMysqlUrl))
	}
	if os.Getenv(envOracleDriver) != "" && os.Getenv(envOracleUrl) != "" {
		urlMap["oracle"] = newSqlDriverAndUrl(os.Getenv(envOracleDriver), os.Getenv(envOracleUrl))
	}
	if os.Getenv(envPgsqlDriver) != "" && os.Getenv(envPgsqlUrl) != "" {
		urlMap["pgsql"] = newSqlDriverAndUrl(os.Getenv(envPgsqlDriver), os.Getenv(envPgsqlUrl))
	}
	if os.Getenv(envCosmosDriver) != "" && os.Getenv(envCosmosUrl) != "" {
		urlMap["cosmosdb"] = newSqlDriverAndUrl(os.Getenv(envCosmosDriver), os.Getenv(envCosmosUrl))
	}
	return urlMap
}

func initSqlConnect(t *testing.T, testName string, dbtype string, info sqlDriverAndUrl) (*promsql.SqlConnect, error) {
	switch dbtype {
	case "sqlite", "sqlite3":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorSqlite)
	case "mssql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorMsSql)
	case "mysql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorMySql)
	case "oracle":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorOracle)
	case "pgsql", "postgresql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorPgSql)
	case "cosmos", "cosmosdb":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, promsql.FlavorCosmosDb)
	default:
		t.Fatalf("%s failed: unknown database type [%s]", testName, dbtype)
	}
	return nil, nil
}

/*----------------------------------------------------------------------*/

func TestNewUserDaoSql(t *testing.T) {
	name := "TestNewUserDaoSql"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype+"/initDaoSql")
			}
		})
	}
}

func TestUserDaoSql_CreateGet(t *testing.T) {
	name := "TestUserDaoSql_CreateGet"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestUserDaoCreateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestUserDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestUserDaoSql_CreateUpdateGet"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestUserDaoCreateUpdateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestUserDaoSql_CreateDelete(t *testing.T) {
	name := "TestUserDaoSql_CreateDelete"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestUserDaoCreateDelete(t, name+"/"+dbtype, dao)
		})
	}
}

func TestUserDaoSql_GetAll(t *testing.T) {
	name := "TestUserDaoSql_GetAll"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestUserDaoGetAll(t, name+"/"+dbtype, dao)
		})
	}
}

func TestUserDaoSql_GetN(t *testing.T) {
	name := "TestUserDaoSql_GetN"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		t.Run(dbtype, func(t *testing.T) {
			sqlc, err := initSqlConnect(t, name, dbtype, info)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
			} else if sqlc == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			defer sqlc.Close()
			err = sqlInitTable(sqlc, testSqlTable)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTable/"+dbtype, err)
			}
			dao := initDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestUserDaoGetN(t, name+"/"+dbtype, dao)
		})
	}
}
