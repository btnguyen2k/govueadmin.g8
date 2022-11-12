package blog

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

	_ "github.com/btnguyen2k/gocosmos"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testDynamodbTableComment = "test_comment"
	testDynamodbTablePost    = "test_post"
	testDynamodbTableVote    = "test_vote"
)

func _createAwsDynamodbConnect(t *testing.T, testName string) *prom.AwsDynamodbConnect {
	awsRegion := strings.ReplaceAll(os.Getenv("AWS_REGION"), `"`, "")
	awsAccessKeyId := strings.ReplaceAll(os.Getenv("AWS_ACCESS_KEY_ID"), `"`, "")
	awsSecretAccessKey := strings.ReplaceAll(os.Getenv("AWS_SECRET_ACCESS_KEY"), `"`, "")
	if awsRegion == "" || awsAccessKeyId == "" || awsSecretAccessKey == "" {
		t.Skipf("%s skipped", testName)
	}
	cfg := &aws.Config{
		Region:      aws.String(awsRegion),
		Credentials: credentials.NewEnvCredentials(),
	}
	if awsDynamodbEndpoint := strings.ReplaceAll(os.Getenv("AWS_DYNAMODB_ENDPOINT"), `"`, ""); awsDynamodbEndpoint != "" {
		cfg.Endpoint = aws.String(awsDynamodbEndpoint)
		if strings.HasPrefix(awsDynamodbEndpoint, "http://") {
			cfg.DisableSSL = aws.Bool(true)
		}
	}
	adc, err := prom.NewAwsDynamodbConnect(cfg, nil, nil, 10000)
	if err != nil {
		t.Fatalf("%s/%s failed: %s", testName, "NewAwsDynamodbConnect", err)
	}
	return adc
}

func _adbDeleteTableAndWait(adc *prom.AwsDynamodbConnect, tableName string) error {
	if err := adc.DeleteTable(nil, tableName); err != nil {
		return err
	}
	for ok, err := adc.HasTable(nil, tableName); (err == nil && ok) || err != nil; {
		if err != nil {
			fmt.Printf("\tError: %s\n", err)
		}
		fmt.Printf("\tTable %s exists, waiting for deletion...\n", tableName)
		time.Sleep(1 * time.Second)
	}

	uidxTableName := tableName + henge.AwsDynamodbUidxTableSuffix
	if err := adc.DeleteTable(nil, uidxTableName); err != nil {
		return err
	}
	for ok, err := adc.HasTable(nil, uidxTableName); (err == nil && ok) || err != nil; {
		if err != nil {
			fmt.Printf("\tError: %s\n", err)
		}
		fmt.Printf("\tTable %s exists, waiting for deletion...\n", uidxTableName)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func initBlogCommentDaoDynamodb(adc *prom.AwsDynamodbConnect) BlogCommentDao {
	return NewBlogCommentDaoDynamodb(adc, testDynamodbTableComment)
}

func initBlogPostDaoDynamodb(adc *prom.AwsDynamodbConnect) BlogPostDao {
	return NewBlogPostDaoDynamodb(adc, testDynamodbTablePost)
}

func initBlogVoteDaoDynamodb(adc *prom.AwsDynamodbConnect) BlogVoteDao {
	return NewBlogVoteDaoDynamodb(adc, testDynamodbTableVote)
}

/*----------------------------------------------------------------------*/
var setupTestDynamodb = func(t *testing.T, testName string) {
	testAdc = _createAwsDynamodbConnect(t, testName)

	_adbDeleteTableAndWait(testAdc, testDynamodbTableComment)
	if err := InitBlogCommentTableDynamodb(testAdc, testDynamodbTableComment); err != nil {
		t.Fatalf("%s failed: error [%s]", testName+"/InitBlogCommentTableDynamodb", err)
	}
	testDaoComment = initBlogCommentDaoDynamodb(testAdc)

	_adbDeleteTableAndWait(testAdc, testDynamodbTablePost)
	if err := InitBlogPostTableDynamodb(testAdc, testDynamodbTablePost); err != nil {
		t.Fatalf("%s failed: error [%s]", testName+"/InitBlogPostTableDynamodb", err)
	}
	testDaoPost = initBlogPostDaoDynamodb(testAdc)

	_adbDeleteTableAndWait(testAdc, testDynamodbTableVote)
	if err := InitBlogVoteTableDynamodb(testAdc, testDynamodbTableVote); err != nil {
		t.Fatalf("%s failed: error [%s]", testName+"/InitBlogVoteTableDynamodb", err)
	}
	testDaoVote = initBlogVoteDaoDynamodb(testAdc)
}

var teardownTestDynamodb = func(t *testing.T, testName string) {
	if testAdc != nil {
		defer func() { testAdc = nil }()
		_adbDeleteTableAndWait(testAdc, testDynamodbTableComment)
		_adbDeleteTableAndWait(testAdc, testDynamodbTablePost)
		_adbDeleteTableAndWait(testAdc, testDynamodbTableVote)
	}
}

/*----------------------------------------------------------------------*/

func TestNewCommentDaoDynamodb(t *testing.T) {
	testName := "TestNewCommentDaoDynamodb"
	adc := _createAwsDynamodbConnect(t, testName)
	defer adc.Close()
	if err := InitBlogCommentTableDynamodb(adc, testDynamodbTableComment); err != nil {
		t.Fatalf("%s failed: error [%s]", testName+"/InitBlogCommentTableDynamodb", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", testName+"/initBlogCommentDaoDynamodb")
	}
}

func TestCommentDaoDynamodb_CreateGet(t *testing.T) {
	testName := "TestCommentDaoDynamodb_CreateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestCommentDaoCreateGet(t, testName, testDaoComment)
}

func TestCommentDaoDynamodb_CreateUpdateGet(t *testing.T) {
	testName := "TestCommentDaoDynamodb_CreateUpdateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestCommentDaoCreateUpdateGet(t, testName, testDaoComment)
}

func TestCommentDaoDynamodb_CreateDelete(t *testing.T) {
	testName := "TestCommentDaoDynamodb_CreateDelete"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestCommentDaoCreateDelete(t, testName, testDaoComment)
}

func TestCommentDaoDynamodb_GetAll(t *testing.T) {
	testName := "TestCommentDaoDynamodb_GetAll"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestCommentDaoGetAll(t, testName, testDaoComment)
}

func TestCommentDaoDynamodb_GetN(t *testing.T) {
	testName := "TestCommentDaoDynamodb_GetN"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestCommentDaoGetN(t, testName, testDaoComment)
}

/*----------------------------------------------------------------------*/

func TestNewPostDaoDynamodb(t *testing.T) {
	testName := "TestNewPostDaoDynamodb"
	adc := _createAwsDynamodbConnect(t, testName)
	defer adc.Close()
	if err := InitBlogPostTableDynamodb(adc, testDynamodbTablePost); err != nil {
		t.Fatalf("%s failed: error [%s]", testName+"/InitBlogPostTableDynamodb", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", testName+"/initBlogPostDaoDynamodb")
	}
}

func TestPostDaoDynamodb_CreateGet(t *testing.T) {
	testName := "TestPostDaoDynamodb_CreateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoCreateGet(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_CreateUpdateGet(t *testing.T) {
	testName := "TestPostDaoDynamodb_CreateUpdateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoCreateUpdateGet(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_CreateDelete(t *testing.T) {
	testName := "TestPostDaoDynamodb_CreateDelete"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoCreateDelete(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_GetUserPostsAll(t *testing.T) {
	testName := "TestPostDaoDynamodb_GetUserPostsAll"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoGetUserPostsAll(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_GetUserPostsN(t *testing.T) {
	testName := "TestPostDaoDynamodb_GetUserPostsN"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoGetUserPostsN(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_GetUserFeedAll(t *testing.T) {
	testName := "TestPostDaoDynamodb_GetUserFeedAll"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	// fmt.Println(teardownTest != nil)
	doTestPostDaoGetUserFeedAll(t, testName, testDaoPost)
}

func TestPostDaoDynamodb_GetUserFeedN(t *testing.T) {
	testName := "TestPostDaoDynamodb_GetUserFeedN"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestPostDaoGetUserFeedN(t, testName, testDaoPost)
}

/*----------------------------------------------------------------------*/

func TestNewVoteDaoDynamodb(t *testing.T) {
	name := "TestNewVoteDaoDynamodb"
	adc := _createAwsDynamodbConnect(t, name)
	defer adc.Close()
	err := InitBlogVoteTableDynamodb(adc, testDynamodbTableVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/InitBlogVoteTableDynamodb", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
}

func TestVoteDaoDynamodb_CreateGet(t *testing.T) {
	testName := "TestVoteDaoDynamodb_CreateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoCreateGet(t, testName, testDaoVote)
}

func TestVoteDaoDynamodb_CreateUpdateGet(t *testing.T) {
	testName := "TestVoteDaoDynamodb_CreateUpdateGet"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoCreateUpdateGet(t, testName, testDaoVote)
}

func TestVoteDaoDynamodb_CreateDelete(t *testing.T) {
	testName := "TestVoteDaoDynamodb_CreateDelete"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoCreateDelete(t, testName, testDaoVote)
}

func TestVoteDaoDynamodb_GetAll(t *testing.T) {
	testName := "TestVoteDaoDynamodb_GetAll"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoGetAll(t, testName, testDaoVote)
}

func TestVoteDaoDynamodb_GetN(t *testing.T) {
	testName := "TestVoteDaoDynamodb_GetN"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoGetN(t, testName, testDaoVote)
}

func TestVoteDaoDynamodb_GetUserVoteForTarget(t *testing.T) {
	testName := "TestVoteDaoDynamodb_GetUserVoteForTarget"
	teardownTest := setupTest(t, testName, setupTestDynamodb, teardownTestDynamodb)
	defer teardownTest(t)
	doTestVoteDaoGetUserVoteForTarget(t, testName, testDaoVote)
}
