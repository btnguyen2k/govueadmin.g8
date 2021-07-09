package user

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

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

func sqlInitTable(sqlc *prom.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		_, err = sqlc.GetDB().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s WITH MAXRU=10000", cosmosdbDbName))
		if err != nil {
			return err
		}
	}
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	// _, err = sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	// if err != nil {
	// 	fmt.Printf("WARNING: %s\n", err)
	// }
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId, Uk: [][]string{{"/" + UserColMaskUid}}}
		err = henge.InitCosmosdbCollection(sqlc, table, spec)
	case prom.FlavorPgSql:
		err = henge.InitPgsqlTable(sqlc, table, map[string]string{UserColMaskUid: "VARCHAR(32)"})
	}
	return err
}

func newSqlConnect(t *testing.T, testName string, driver, url, timezone string, flavor prom.DbFlavor) (*prom.SqlConnect, error) {
	driver = strings.Trim(driver, "\"")
	url = strings.Trim(url, "\"")
	if driver == "" || url == "" {
		t.Skipf("%s skipped", testName)
	}

	urlTimezone := strings.ReplaceAll(timezone, "/", "%2f")
	url = strings.ReplaceAll(url, "${loc}", urlTimezone)
	url = strings.ReplaceAll(url, "${tz}", urlTimezone)
	url = strings.ReplaceAll(url, "${timezone}", urlTimezone)
	sqlc, err := prom.NewSqlConnectWithFlavor(driver, url, 10000, nil, flavor)
	if err == nil && sqlc != nil {
		loc, _ := time.LoadLocation(timezone)
		sqlc.SetLocation(loc)
	}
	return sqlc, err
}

func initDaoSql(sqlc *prom.SqlConnect) UserDao {
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
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

func initSqlConnect(t *testing.T, testName string, dbtype string, info sqlDriverAndUrl) (*prom.SqlConnect, error) {
	switch dbtype {
	case "sqlite", "sqlite3":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorSqlite)
	case "mssql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorMsSql)
	case "mysql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorMySql)
	case "oracle":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorOracle)
	case "pgsql", "postgresql":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorPgSql)
	case "cosmos", "cosmosdb":
		return newSqlConnect(t, testName, info.driver, info.url, testTimeZone, prom.FlavorCosmosDb)
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
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name+"/initDaoSql")
		}
		sqlc.Close()
	}
}

func TestUserDaoSql_CreateGet(t *testing.T) {
	name := "TestUserDaoSql_CreateGet"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name)
		}
		doTestUserDaoCreateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestUserDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestUserDaoSql_CreateUpdateGet"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name)
		}
		doTestUserDaoCreateUpdateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestUserDaoSql_CreateDelete(t *testing.T) {
	name := "TestUserDaoSql_CreateDelete"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name)
		}
		doTestUserDaoCreateDelete(t, name, dao)
		sqlc.Close()
	}
}

func TestUserDaoSql_GetAll(t *testing.T) {
	name := "TestUserDaoSql_GetAll"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name)
		}
		doTestUserDaoGetAll(t, name, dao)
		sqlc.Close()
	}
}

func TestUserDaoSql_GetN(t *testing.T) {
	name := "TestUserDaoSql_GetN"
	urlMap := sqlGetUrlFromEnv()
	if len(urlMap) == 0 {
		t.Skipf("%s skipped", name)
	}
	for dbtype, info := range urlMap {
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTable(sqlc, testSqlTable)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTable/"+dbtype, err)
		}
		dao := initDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name)
		}
		doTestUserDaoGetN(t, name, dao)
		sqlc.Close()
	}
}
