package blog

import (
	"errors"
	"math/rand"
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

func _dynamodbWaitForTableStatus(adc *prom.AwsDynamodbConnect, table, status string, timeout time.Duration) error {
	t := time.Now()
	for tblStatus, err := adc.GetTableStatus(nil, table); ; {
		if err != nil {
			return err
		}
		if strings.ToUpper(tblStatus) == status {
			return nil
		}
		if time.Now().Sub(t).Milliseconds() > timeout.Milliseconds() {
			return errors.New("")
		}
	}
}

func dynamodbInitTable(adc *prom.AwsDynamodbConnect, table string, spec *henge.DynamodbTablesSpec) error {
	rand.Seed(time.Now().UnixNano())
	adc.DeleteTable(nil, table)
	if err := _dynamodbWaitForTableStatus(adc, table, "", 10*time.Second); err != nil {
		return err
	}
	if spec.CreateUidxTable {
		adc.DeleteTable(nil, table+henge.AwsDynamodbUidxTableSuffix)
		if err := _dynamodbWaitForTableStatus(adc, table+henge.AwsDynamodbUidxTableSuffix, "", 10*time.Second); err != nil {
			return err
		}
	}
	return henge.InitDynamodbTables(adc, table, spec)
}

func newDynamodbConnect(t *testing.T, testName string) (*prom.AwsDynamodbConnect, error) {
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
	return prom.NewAwsDynamodbConnect(cfg, nil, nil, 10000)
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

func TestNewCommentDaoDynamodb(t *testing.T) {
	name := "TestNewCommentDaoDynamodb"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
}

func TestCommentDaoDynamodb_CreateGet(t *testing.T) {
	name := "TestCommentDaoDynamodb_CreateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
	doTestCommentDaoCreateGet(t, name, dao)
}

func TestCommentDaoDynamodb_CreateUpdateGet(t *testing.T) {
	name := "TestCommentDaoDynamodb_CreateUpdateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
	doTestCommentDaoCreateUpdateGet(t, name, dao)
}

func TestCommentDaoDynamodb_CreateDelete(t *testing.T) {
	name := "TestCommentDaoDynamodb_CreateDelete"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
	doTestCommentDaoCreateDelete(t, name, dao)
}

func TestCommentDaoDynamodb_GetAll(t *testing.T) {
	name := "TestCommentDaoDynamodb_GetAll"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
	doTestCommentDaoGetAll(t, name, dao)
}

func TestCommentDaoDynamodb_GetN(t *testing.T) {
	name := "TestCommentDaoDynamodb_GetN"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableComment, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogCommentDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoDynamodb")
	}
	defer adc.Close()
	doTestCommentDaoGetN(t, name, dao)
}

/*----------------------------------------------------------------------*/

func TestNewPostDaoDynamodb(t *testing.T) {
	name := "TestNewPostDaoDynamodb"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
}

func TestPostDaoDynamodb_CreateGet(t *testing.T) {
	name := "TestPostDaoDynamodb_CreateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoCreateGet(t, name, dao)
}

func TestPostDaoDynamodb_CreateUpdateGet(t *testing.T) {
	name := "TestPostDaoDynamodb_CreateUpdateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoCreateUpdateGet(t, name, dao)
}

func TestPostDaoDynamodb_CreateDelete(t *testing.T) {
	name := "TestPostDaoDynamodb_CreateDelete"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoCreateDelete(t, name, dao)
}

func TestPostDaoDynamodb_GetUserPostsAll(t *testing.T) {
	name := "TestPostDaoDynamodb_GetUserPostsAll"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoGetUserPostsAll(t, name, dao)
}

func TestPostDaoDynamodb_GetUserPostsN(t *testing.T) {
	name := "TestPostDaoDynamodb_GetUserPostsN"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoGetUserPostsN(t, name, dao)
}

func TestPostDaoDynamodb_GetUserFeedAll(t *testing.T) {
	name := "TestPostDaoDynamodb_GetUserFeedAll"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoGetUserFeedAll(t, name, dao)
}

func TestPostDaoDynamodb_GetUserFeedN(t *testing.T) {
	name := "TestPostDaoDynamodb_GetUserFeedN"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTablePost, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogPostDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoDynamodb")
	}
	defer adc.Close()
	doTestPostDaoGetUserFeedN(t, name, dao)
}

/*----------------------------------------------------------------------*/

func TestNewVoteDaoDynamodb(t *testing.T) {
	name := "TestNewVoteDaoDynamodb"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
}

func TestVoteDaoDynamodb_CreateGet(t *testing.T) {
	name := "TestVoteDaoDynamodb_CreateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoCreateGet(t, name, dao)
}

func TestVoteDaoDynamodb_CreateUpdateGet(t *testing.T) {
	name := "TestVoteDaoDynamodb_CreateUpdateGet"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoCreateUpdateGet(t, name, dao)
}

func TestVoteDaoDynamodb_CreateDelete(t *testing.T) {
	name := "TestVoteDaoDynamodb_CreateDelete"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoCreateDelete(t, name, dao)
}

func TestVoteDaoDynamodb_GetAll(t *testing.T) {
	name := "TestVoteDaoDynamodb_GetAll"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoGetAll(t, name, dao)
}

func TestVoteDaoDynamodb_GetN(t *testing.T) {
	name := "TestVoteDaoDynamodb_GetN"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoGetN(t, name, dao)
}

func TestVoteDaoDynamodb_GetUserVoteForTarget(t *testing.T) {
	name := "TestVoteDaoDynamodb_GetUserVoteForTarget"
	adc, err := newDynamodbConnect(t, name)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name, err)
	} else if adc == nil {
		t.Fatalf("%s failed: nil", name)
	}
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	err = dynamodbInitTable(adc, testDynamodbTableVote, spec)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/dynamodbInitTable", err)
	}
	dao := initBlogVoteDaoDynamodb(adc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoDynamodb")
	}
	defer adc.Close()
	doTestVoteDaoGetUserVoteForTarget(t, name, dao)
}
