package blog

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
	testTimeZone        = "Asia/Ho_Chi_Minh"
	testSqlTableComment = "test_comment"
	testSqlTablePost    = "test_post"
	testSqlTableVote    = "test_vote"
)

func sqlInitTableComment(sqlc *promsql.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	extraCols := map[string]string{
		CommentColParentId: "VARCHAR(32)",
		CommentColPostId:   "VARCHAR(32)",
		CommentColOwnerId:  "VARCHAR(32)",
	}
	switch sqlc.GetDbFlavor() {
	case promsql.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
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

func sqlInitTablePost(sqlc *promsql.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	extraCols := map[string]string{PostColOwnerId: "VARCHAR(32)", PostColIsPublic: "INT"}
	switch sqlc.GetDbFlavor() {
	case promsql.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
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

func sqlInitTableVote(sqlc *promsql.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	extraCols := map[string]string{VoteColOwnerId: "VARCHAR(32)", VoteColTargetId: "VARCHAR(32)", VoteColValue: "INT"}
	switch sqlc.GetDbFlavor() {
	case promsql.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
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

func initBlogCommentDaoSql(sqlc *promsql.SqlConnect) BlogCommentDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return NewBlogCommentDaoCosmosdb(sqlc, testSqlTableComment, true)
	}
	return NewBlogCommentDaoSql(sqlc, testSqlTableComment, true)
}

func initBlogPostDaoSql(sqlc *promsql.SqlConnect) BlogPostDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return NewBlogPostDaoCosmosdb(sqlc, testSqlTablePost, true)
	}
	return NewBlogPostDaoSql(sqlc, testSqlTablePost, true)
}

func initBlogVoteDaoSql(sqlc *promsql.SqlConnect) BlogVoteDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return NewBlogVoteDaoCosmosdb(sqlc, testSqlTableVote, true)
	}
	return NewBlogVoteDaoSql(sqlc, testSqlTableVote, true)
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

func TestNewCommentDaoSql(t *testing.T) {
	name := "TestNewCommentDaoSql"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype+"/initBlogCommentDaoSql")
			}
		})
	}
}

func TestCommentDaoSql_CreateGet(t *testing.T) {
	name := "TestCommentDaoSql_CreateGet"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestCommentDaoCreateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestCommentDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestCommentDaoSql_CreateUpdateGet"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestCommentDaoCreateUpdateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestCommentDaoSql_CreateDelete(t *testing.T) {
	name := "TestCommentDaoSql_CreateDelete"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestCommentDaoCreateDelete(t, name+"/"+dbtype, dao)
		})
	}
}

func TestCommentDaoSql_GetAll(t *testing.T) {
	name := "TestCommentDaoSql_GetAll"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestCommentDaoGetAll(t, name+"/"+dbtype, dao)
		})
	}
}

func TestCommentDaoSql_GetN(t *testing.T) {
	name := "TestCommentDaoSql_GetN"
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
			err = sqlInitTableComment(sqlc, testSqlTableComment)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableComment/"+dbtype, err)
			}
			dao := initBlogCommentDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestCommentDaoGetN(t, name+"/"+dbtype, dao)
		})
	}
}

/*----------------------------------------------------------------------*/

func TestNewPostDaoSql(t *testing.T) {
	name := "TestNewPostDaoSql"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype+"/initBlogPostDaoSql")
			}
		})
	}
}

func TestPostDaoSql_CreateGet(t *testing.T) {
	name := "TestPostDaoSql_CreateGet"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoCreateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestPostDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestPostDaoSql_CreateUpdateGet"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoCreateUpdateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestPostDaoSql_CreateDelete(t *testing.T) {
	name := "TestPostDaoSql_CreateDelete"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoCreateDelete(t, name+"/"+dbtype, dao)
		})
	}
}

func TestPostDaoSql_GetUserPostsAll(t *testing.T) {
	name := "TestPostDaoSql_GetUserPostsAll"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoGetUserPostsAll(t, name+"/"+dbtype, dao)

		})
	}
}

func TestPostDaoSql_GetUserPostsN(t *testing.T) {
	name := "TestPostDaoSql_GetUserPostsN"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoGetUserPostsN(t, name+"/"+dbtype, dao)
		})
	}
}

func TestPostDaoSql_GetUserFeedAll(t *testing.T) {
	name := "TestPostDaoSql_GetUserFeedAll"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoGetUserFeedAll(t, name+"/"+dbtype, dao)
		})
	}
}

func TestPostDaoSql_GetUserFeedN(t *testing.T) {
	name := "TestPostDaoSql_GetUserFeedN"
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
			err = sqlInitTablePost(sqlc, testSqlTablePost)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTablePost/"+dbtype, err)
			}
			dao := initBlogPostDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestPostDaoGetUserFeedN(t, name+"/"+dbtype, dao)
		})
	}
}

/*----------------------------------------------------------------------*/

func TestNewVoteDaoSql(t *testing.T) {
	name := "TestNewVoteDaoSql"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype+"/initBlogVoteDaoSql")
			}
		})
	}
}

func TestVoteDaoSql_CreateGet(t *testing.T) {
	name := "TestVoteDaoSql_CreateGet"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoCreateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestVoteDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestVoteDaoSql_CreateUpdateGet"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoCreateUpdateGet(t, name+"/"+dbtype, dao)
		})
	}
}

func TestVoteDaoSql_CreateDelete(t *testing.T) {
	name := "TestVoteDaoSql_CreateDelete"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoCreateDelete(t, name+"/"+dbtype, dao)
		})
	}
}

func TestVoteDaoSql_GetAll(t *testing.T) {
	name := "TestVoteDaoSql_GetAll"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoGetAll(t, name+"/"+dbtype, dao)
		})
	}
}

func TestVoteDaoSql_GetN(t *testing.T) {
	name := "TestVoteDaoSql_GetN"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoGetN(t, name+"/"+dbtype, dao)
		})
	}
}

func TestVoteDaoSql_GetUserVoteForTarget(t *testing.T) {
	name := "TestVoteDaoSql_GetUserVoteForTarget"
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
			err = sqlInitTableVote(sqlc, testSqlTableVote)
			if err != nil {
				t.Fatalf("%s failed: error [%s]", name+"/"+dbtype+"/sqlInitTableVote/"+dbtype, err)
			}
			dao := initBlogVoteDaoSql(sqlc)
			if dao == nil {
				t.Fatalf("%s failed: nil", name+"/"+dbtype)
			}
			doTestVoteDaoGetUserVoteForTarget(t, name+"/"+dbtype, dao)
		})
	}
}
