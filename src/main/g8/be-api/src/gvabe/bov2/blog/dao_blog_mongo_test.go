package blog

import (
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

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
	testMongoCollectionComment = "test_comment"
	testMongoCollectionPost    = "test_post"
	testMongoCollectionVote    = "test_vote"
)

func mongoInitCollection(mc *prom.MongoConnect, collection string) error {
	rand.Seed(time.Now().UnixNano())
	mc.GetCollection(collection).Drop(nil)
	return henge.InitMongoCollection(mc, collection)
}

func newMongoConnect(t *testing.T, testName string, db, url string) (*prom.MongoConnect, error) {
	db = strings.Trim(db, "\"")
	url = strings.Trim(url, "\"")
	if db == "" || url == "" {
		t.Skipf("%s skipped", testName)
	}

	return prom.NewMongoConnect(url, db, 10000)
}

func initBlogCommentDaoMongo(mc *prom.MongoConnect) BlogCommentDao {
	url := mc.GetUrl()
	txModeOnWrite := strings.Index(url, "replicaSet=") >= 0
	return NewBlogCommentDaoMongo(mc, testMongoCollectionComment, txModeOnWrite)
}

func initBlogPostDaoMongo(mc *prom.MongoConnect) BlogPostDao {
	url := mc.GetUrl()
	txModeOnWrite := strings.Index(url, "replicaSet=") >= 0
	return NewBlogPostDaoMongo(mc, testMongoCollectionPost, txModeOnWrite)
}

func initBlogVoteDaoMongo(mc *prom.MongoConnect) BlogVoteDao {
	url := mc.GetUrl()
	txModeOnWrite := strings.Index(url, "replicaSet=") >= 0
	return NewBlogVoteDaoMongo(mc, testMongoCollectionVote, txModeOnWrite)
}

const (
	envMongoDb  = "MONGO_DB"
	envMongoUrl = "MONGO_URL"
)

/*----------------------------------------------------------------------*/

func TestNewCommentDaoMongo(t *testing.T) {
	name := "TestNewCommentDaoMongo"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogCommentDaoMongo")
	}
}

func TestCommentDaoMongo_CreateGet(t *testing.T) {
	name := "TestCommentDaoMongo_CreateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	doTestCommentDaoCreateGet(t, name, dao)
}

func TestCommentDaoMongo_CreateUpdateGet(t *testing.T) {
	name := "TestCommentDaoMongo_CreateUpdateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	doTestCommentDaoCreateUpdateGet(t, name, dao)
}

func TestCommentDaoMongo_CreateDelete(t *testing.T) {
	name := "TestCommentDaoMongo_CreateDelete"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	doTestCommentDaoCreateDelete(t, name, dao)
}

func TestCommentDaoMongo_GetAll(t *testing.T) {
	name := "TestCommentDaoMongo_GetAll"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	doTestCommentDaoGetAll(t, name, dao)
}

func TestCommentDaoMongo_GetN(t *testing.T) {
	name := "TestCommentDaoMongo_GetN"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionComment)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogCommentDaoMongo(mc)
	doTestCommentDaoGetN(t, name, dao)
}

/*----------------------------------------------------------------------*/

func TestNewPostDaoMongo(t *testing.T) {
	name := "TestNewPostDaoMongo"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
}

func TestPostDaoMongo_CreateGet(t *testing.T) {
	name := "TestPostDaoMongo_CreateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoCreateGet(t, name, dao)
}

func TestPostDaoMongo_CreateUpdateGet(t *testing.T) {
	name := "TestPostDaoMongo_CreateUpdateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoCreateUpdateGet(t, name, dao)
}

func TestPostDaoMongo_CreateDelete(t *testing.T) {
	name := "TestPostDaoMongo_CreateDelete"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoCreateDelete(t, name, dao)
}

func TestPostDaoMongo_GetUserPostsAll(t *testing.T) {
	name := "TestPostDaoMongo_GetUserPostsAll"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoGetUserPostsAll(t, name, dao)
}

func TestPostDaoMongo_GetUserPostsN(t *testing.T) {
	name := "TestPostDaoMongo_GetUserPostsN"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoGetUserPostsN(t, name, dao)
}

func TestPostDaoMongo_GetUserFeedAll(t *testing.T) {
	name := "TestPostDaoMongo_GetUserFeedAll"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoGetUserFeedAll(t, name, dao)
}

func TestPostDaoMongo_GetUserFeedN(t *testing.T) {
	name := "TestPostDaoMongo_GetUserFeedN"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionPost)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogPostDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogPostDaoMongo")
	}
	doTestPostDaoGetUserFeedN(t, name, dao)
}

/*----------------------------------------------------------------------*/

func TestNewVoteDaoMongo(t *testing.T) {
	name := "TestNewVoteDaoMongo"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
}

func TestVoteDaoMongo_CreateGet(t *testing.T) {
	name := "TestVoteDaoMongo_CreateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoCreateGet(t, name, dao)
}

func TestVoteDaoMongo_CreateUpdateGet(t *testing.T) {
	name := "TestVoteDaoMongo_CreateUpdateGet"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoCreateUpdateGet(t, name, dao)
}

func TestVoteDaoMongo_CreateDelete(t *testing.T) {
	name := "TestVoteDaoMongo_CreateDelete"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoCreateDelete(t, name, dao)
}

func TestVoteDaoMongo_GetAll(t *testing.T) {
	name := "TestVoteDaoMongo_GetAll"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoGetAll(t, name, dao)
}

func TestVoteDaoMongo_GetN(t *testing.T) {
	name := "TestVoteDaoMongo_GetN"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoGetN(t, name, dao)
}

func TestVoteDaoMongo_GetUserVoteForTarget(t *testing.T) {
	name := "TestVoteDaoMongo_GetUserVoteForTarget"
	db := os.Getenv(envMongoDb)
	url := os.Getenv(envMongoUrl)
	mc, err := newMongoConnect(t, name, db, url)
	if err != nil {
		t.Fatalf("%s faied: %s", name, err)
	}
	defer mc.Close(nil)
	err = mongoInitCollection(mc, testMongoCollectionVote)
	if err != nil {
		t.Fatalf("%s failed: error [%s]", name+"/mongoInitCollection", err)
	}
	dao := initBlogVoteDaoMongo(mc)
	if dao == nil {
		t.Fatalf("%s failed: nil", name+"/initBlogVoteDaoMongo")
	}
	doTestVoteDaoGetUserVoteForTarget(t, name, dao)
}
