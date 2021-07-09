package blog

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
	testTimeZone        = "Asia/Ho_Chi_Minh"
	testSqlTableComment = "test_comment"
	testSqlTablePost    = "test_post"
	testSqlTableVote    = "test_vote"
)

func sqlInitTableComment(sqlc *prom.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		_, err = sqlc.GetDB().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s WITH MAXRU=10000", cosmosdbDbName))
		if err != nil {
			return err
		}
	}
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
		err = henge.InitCosmosdbCollection(sqlc, table, spec)
	case prom.FlavorPgSql:
		err = henge.InitPgsqlTable(sqlc, table, map[string]string{
			CommentColParentId: "VARCHAR(32)",
			CommentColPostId:   "VARCHAR(32)",
			CommentColOwnerId:  "VARCHAR(32)",
		})
	}
	return err
}

func sqlInitTablePost(sqlc *prom.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		_, err = sqlc.GetDB().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s WITH MAXRU=10000", cosmosdbDbName))
		if err != nil {
			return err
		}
	}
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
		err = henge.InitCosmosdbCollection(sqlc, table, spec)
	case prom.FlavorPgSql:
		err = henge.InitPgsqlTable(sqlc, table, map[string]string{PostColOwnerId: "VARCHAR(32)", PostColIsPublic: "INT"})
	}
	return err
}

func sqlInitTableVote(sqlc *prom.SqlConnect, table string) error {
	rand.Seed(time.Now().UnixNano())
	var err error
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		_, err = sqlc.GetDB().Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s WITH MAXRU=10000", cosmosdbDbName))
		if err != nil {
			return err
		}
	}
	sqlc.GetDB().Exec(fmt.Sprintf("DROP TABLE %s", table))
	switch sqlc.GetDbFlavor() {
	case prom.FlavorCosmosDb:
		spec := &henge.CosmosdbCollectionSpec{Pk: henge.CosmosdbColId}
		err = henge.InitCosmosdbCollection(sqlc, table, spec)
	case prom.FlavorPgSql:
		err = henge.InitPgsqlTable(sqlc, table, map[string]string{VoteColOwnerId: "VARCHAR(32)", VoteColTargetId: "VARCHAR(32)", VoteColValue: "INT"})
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

func initBlogCommentDaoSql(sqlc *prom.SqlConnect) BlogCommentDao {
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		return NewBlogCommentDaoCosmosdb(sqlc, testSqlTableComment, true)
	}
	return NewBlogCommentDaoSql(sqlc, testSqlTableComment, true)
}

func initBlogPostDaoSql(sqlc *prom.SqlConnect) BlogPostDao {
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
		return NewBlogPostDaoCosmosdb(sqlc, testSqlTablePost, true)
	}
	return NewBlogPostDaoSql(sqlc, testSqlTablePost, true)
}

func initBlogVoteDaoSql(sqlc *prom.SqlConnect) BlogVoteDao {
	if sqlc.GetDbFlavor() == prom.FlavorCosmosDb {
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

func TestNewCommentDaoSql(t *testing.T) {
	name := "TestNewCommentDaoSql"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoSql")
		}
		sqlc.Close()
	}
}

func TestCommentDaoSql_CreateGet(t *testing.T) {
	name := "TestCommentDaoSql_CreateGet"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		doTestCommentDaoCreateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestCommentDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestCommentDaoSql_CreateUpdateGet"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		doTestCommentDaoCreateUpdateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestCommentDaoSql_CreateDelete(t *testing.T) {
	name := "TestCommentDaoSql_CreateDelete"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		doTestCommentDaoCreateDelete(t, name, dao)
		sqlc.Close()
	}
}

func TestCommentDaoSql_GetAll(t *testing.T) {
	name := "TestCommentDaoSql_GetAll"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		doTestCommentDaoGetAll(t, name, dao)
		sqlc.Close()
	}
}

func TestCommentDaoSql_GetN(t *testing.T) {
	name := "TestCommentDaoSql_GetN"
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
		err = sqlInitTableComment(sqlc, testSqlTableComment)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableComment/"+dbtype, err)
		}
		dao := initBlogCommentDaoSql(sqlc)
		doTestCommentDaoGetN(t, name, dao)
		sqlc.Close()
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
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name+"/initBlogPostDaoSql")
		}
		sqlc.Close()
	}
}

func TestPostDaoSql_CreateGet(t *testing.T) {
	name := "TestPostDaoSql_CreateGet"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoCreateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestPostDaoSql_CreateUpdateGet"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoCreateUpdateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_CreateDelete(t *testing.T) {
	name := "TestPostDaoSql_CreateDelete"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoCreateDelete(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_GetUserPostsAll(t *testing.T) {
	name := "TestPostDaoSql_GetUserPostsAll"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoGetUserPostsAll(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_GetUserPostsN(t *testing.T) {
	name := "TestPostDaoSql_GetUserPostsN"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoGetUserPostsN(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_GetUserFeedAll(t *testing.T) {
	name := "TestPostDaoSql_GetUserFeedAll"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoGetUserFeedAll(t, name, dao)
		sqlc.Close()
	}
}

func TestPostDaoSql_GetUserFeedN(t *testing.T) {
	name := "TestPostDaoSql_GetUserFeedN"
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
		err = sqlInitTablePost(sqlc, testSqlTablePost)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTablePost/"+dbtype, err)
		}
		dao := initBlogPostDaoSql(sqlc)
		doTestPostDaoGetUserFeedN(t, name, dao)
		sqlc.Close()
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
		sqlc, err := initSqlConnect(t, name, dbtype, info)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/"+dbtype, err)
		} else if sqlc == nil {
			t.Fatalf("%s failed: nil", name+"/"+dbtype)
		}
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		if dao == nil {
			t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoSql")
		}
		sqlc.Close()
	}
}

func TestVoteDaoSql_CreateGet(t *testing.T) {
	name := "TestVoteDaoSql_CreateGet"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoCreateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestVoteDaoSql_CreateUpdateGet(t *testing.T) {
	name := "TestVoteDaoSql_CreateUpdateGet"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoCreateUpdateGet(t, name, dao)
		sqlc.Close()
	}
}

func TestVoteDaoSql_CreateDelete(t *testing.T) {
	name := "TestVoteDaoSql_CreateDelete"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoCreateDelete(t, name, dao)
		sqlc.Close()
	}
}

func TestVoteDaoSql_GetAll(t *testing.T) {
	name := "TestVoteDaoSql_GetAll"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoGetAll(t, name, dao)
		sqlc.Close()
	}
}

func TestVoteDaoSql_GetN(t *testing.T) {
	name := "TestVoteDaoSql_GetN"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoGetN(t, name, dao)
		sqlc.Close()
	}
}

func TestVoteDaoSql_GetUserVoteForTarget(t *testing.T) {
	name := "TestVoteDaoSql_GetUserVoteForTarget"
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
		err = sqlInitTableVote(sqlc, testSqlTableVote)
		if err != nil {
			t.Fatalf("%s failed: error [%s]", name+"/sqlInitTableVote/"+dbtype, err)
		}
		dao := initBlogVoteDaoSql(sqlc)
		doTestVoteDaoGetUserVoteForTarget(t, name, dao)
		sqlc.Close()
	}
}
