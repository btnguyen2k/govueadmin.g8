package gvabe

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/src/goapi"
	"main/src/gvabe/bov2/blog"
	"main/src/gvabe/bov2/user"
	"main/src/utils"

	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func _createDynamodbConnect(dbtype string) *prom.AwsDynamodbConnect {
	var adc *prom.AwsDynamodbConnect = nil
	var err error
	switch dbtype {
	case "dynamo", "dynamodb", "awsdynamo", "awsdynamodb":
		region := goapi.AppConfig.GetString("gvabe.db.dynamodb.region")
		region = strings.ReplaceAll(region, `"`, "")
		cfg := &aws.Config{
			Region:      aws.String(region),
			Credentials: credentials.NewEnvCredentials(),
		}
		endpoint := goapi.AppConfig.GetString("gvabe.db.dynamodb.endpoint")
		endpoint = strings.ReplaceAll(endpoint, `"`, "")
		if endpoint != "" {
			cfg.Endpoint = aws.String(endpoint)
			if strings.HasPrefix(strings.ToLower(endpoint), "http://") {
				cfg.DisableSSL = aws.Bool(true)
			}
		}
		adc, err = prom.NewAwsDynamodbConnect(cfg, nil, nil, 10000)
	}
	if err != nil {
		panic(err)
	}
	return adc
}

func _createSqlConnect(dbtype string) *prom.SqlConnect {
	timezone := goapi.AppConfig.GetString("timezone")
	var sqlc *prom.SqlConnect = nil
	var err error
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		sqlc, err = henge.NewSqliteConnection(dir, dbname, timezone, "sqlite3", 10000, nil)
	case "pg", "pgsql", "postgres", "postgresql":
		url := goapi.AppConfig.GetString("gvabe.db.pgsql.url")
		sqlc, err = henge.NewPgsqlConnection(url, timezone, "pgx", 10000, nil)
	}
	if err != nil {
		panic(err)
	}
	return sqlc
}

func _createMongoConnect(dbtype string) *prom.MongoConnect {
	var mc *prom.MongoConnect = nil
	var err error
	switch dbtype {
	case "mongo", "mongodb":
		db := goapi.AppConfig.GetString("gvabe.db.mongodb.db")
		url := goapi.AppConfig.GetString("gvabe.db.mongodb.url")
		mc, err = prom.NewMongoConnect(url, db, 10000)
	}
	if err != nil {
		panic(err)
	}
	return mc
}

func _createUserDaoSql(sqlc *prom.SqlConnect) user.UserDao {
	return user.NewUserDaoSql(sqlc, user.TableUser, true)
}
func _createUserDaoDynamodb(adc *prom.AwsDynamodbConnect) user.UserDao {
	return user.NewUserDaoDynamodb(adc, user.TableUser)
}
func _createUserDaoMongo(mc *prom.MongoConnect) user.UserDao {
	url := mc.GetUrl()
	return user.NewUserDaoMongo(mc, user.TableUser, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogPostDaoSql(sqlc *prom.SqlConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoSql(sqlc, blog.TableBlogPost, true)
}
func _createBlogPostDaoDynamodb(adc *prom.AwsDynamodbConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoDynamodb(adc, blog.TableBlogPost)
}
func _createBlogPostDaoMongo(mc *prom.MongoConnect) blog.BlogPostDao {
	url := mc.GetUrl()
	return blog.NewBlogPostDaoMongo(mc, blog.TableBlogPost, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogCommentDaoSql(sqlc *prom.SqlConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoSql(sqlc, blog.TableBlogComment, true)
}
func _createBlogCommentDaoDynamodb(adc *prom.AwsDynamodbConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoDynamodb(adc, blog.TableBlogComment)
}
func _createBlogCommentDaoMongo(mc *prom.MongoConnect) blog.BlogCommentDao {
	url := mc.GetUrl()
	return blog.NewBlogCommentDaoMongo(mc, blog.TableBlogComment, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogVoteDaoSql(sqlc *prom.SqlConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoSql(sqlc, blog.TableBlogVote, true)
}
func _createBlogVoteDaoDynamodb(adc *prom.AwsDynamodbConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoDynamodb(adc, blog.TableBlogVote)
}
func _createBlogVoteDaoMongo(mc *prom.MongoConnect) blog.BlogVoteDao {
	url := mc.GetUrl()
	return blog.NewBlogVoteDaoMongo(mc, blog.TableBlogVote, strings.Index(url, "replicaSet=") >= 0)
}

func _createSqlTables(sqlc *prom.SqlConnect, dbtype string) {
	switch dbtype {
	case "sqlite":
		if err := henge.InitSqliteTable(sqlc, user.TableUser, map[string]string{user.UserColMaskUid: "VARCHAR(32)"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", user.TableUser, dbtype, err)
		}
		if err := henge.InitSqliteTable(sqlc, blog.TableBlogPost, map[string]string{
			blog.PostColOwnerId: "VARCHAR(32)", blog.PostColIsPublic: "INT"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogPost, dbtype, err)
		}
		if err := henge.InitSqliteTable(sqlc, blog.TableBlogComment, map[string]string{
			blog.CommentColOwnerId: "VARCHAR(32)", blog.CommentColPostId: "VARCHAR(32)", blog.CommentColParentId: "VARCHAR(32)"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogComment, dbtype, err)
		}
		if err := henge.InitSqliteTable(sqlc, blog.TableBlogVote, map[string]string{
			blog.VoteColOwnerId: "VARCHAR(32)", blog.VoteColTargetId: "VARCHAR(32)", blog.VoteColValue: "INT"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogVote, dbtype, err)
		}
	case "pg", "pgsql", "postgres", "postgresql":
		if err := henge.InitPgsqlTable(sqlc, user.TableUser, map[string]string{user.UserColMaskUid: "VARCHAR(32)"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", user.TableUser, dbtype, err)
		}
		if err := henge.InitPgsqlTable(sqlc, blog.TableBlogPost, map[string]string{
			blog.PostColOwnerId: "VARCHAR(32)", blog.PostColIsPublic: "INT"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogPost, dbtype, err)
		}
		if err := henge.InitPgsqlTable(sqlc, blog.TableBlogComment, map[string]string{
			blog.CommentColOwnerId: "VARCHAR(32)", blog.CommentColPostId: "VARCHAR(32)", blog.CommentColParentId: "VARCHAR(32)"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogComment, dbtype, err)
		}
		if err := henge.InitPgsqlTable(sqlc, blog.TableBlogVote, map[string]string{
			blog.VoteColOwnerId: "VARCHAR(32)", blog.VoteColTargetId: "VARCHAR(32)", blog.VoteColValue: "INT"}); err != nil {
			log.Printf("[WARN] creating table %s (%s): %s\n", blog.TableBlogVote, dbtype, err)
		}
	}

	// user
	if err := henge.CreateIndexSql(sqlc, user.TableUser, true, []string{user.UserColMaskUid}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", user.TableUser, user.UserColMaskUid, dbtype, err)
	}

	// blog post
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogPost, false, []string{blog.PostColOwnerId}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostColOwnerId, dbtype, err)
	}
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogPost, false, []string{blog.PostColIsPublic}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostColIsPublic, dbtype, err)
	}

	// blog comment
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogComment, false, []string{blog.CommentColOwnerId}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentColOwnerId, dbtype, err)
	}
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogComment, false, []string{blog.CommentColPostId, blog.CommentColParentId}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentColPostId+":"+blog.CommentColParentId, dbtype, err)
	}

	// blog vote
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogVote, true, []string{blog.VoteColOwnerId, blog.VoteColTargetId}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteColOwnerId+":"+blog.VoteColTargetId, dbtype, err)
	}
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogVote, false, []string{blog.VoteColTargetId, blog.VoteColValue}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteColTargetId+":"+blog.VoteColValue, dbtype, err)
	}
}

func _dynamodbWaitforGSI(adc *prom.AwsDynamodbConnect, table, gsi string, timeout time.Duration) error {
	t := time.Now()
	for status, err := adc.GetGlobalSecondaryIndexStatus(nil, table, gsi); ; {
		if err != nil {
			return err
		}
		if strings.ToUpper(status) == "ACTIVE" {
			return nil
		}
		if time.Now().Sub(t).Milliseconds() > timeout.Milliseconds() {
			return errors.New("")
		}
	}
}

func _createDynamodbTables(adc *prom.AwsDynamodbConnect) {
	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, user.TableUser, spec); err != nil {
		log.Printf("[WARN] creating tableName %s (%s): %s\n", user.TableUser, "DynamoDB", err)
	}
	spec = &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, blog.TableBlogPost, spec); err != nil {
		log.Printf("[WARN] creating tableName %s (%s): %s\n", blog.TableBlogPost, "DynamoDB", err)
	}
	spec = &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, blog.TableBlogComment, spec); err != nil {
		log.Printf("[WARN] creating tableName %s (%s): %s\n", blog.TableBlogComment, "DynamoDB", err)
	}
	spec = &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, blog.TableBlogVote, spec); err != nil {
		log.Printf("[WARN] creating tableName %s (%s): %s\n", blog.TableBlogVote, "DynamoDB", err)
	}

	var tableName, gsiName, colName string

	// user
	tableName, colName, gsiName = user.TableUser, user.UserFieldMaskId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsAttrTypeString}},
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsKeyTypePartition}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}

	// blog post
	tableName, colName, gsiName = blog.TableBlogPost, blog.PostFieldOwnerId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsAttrTypeString}},
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsKeyTypePartition}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}
	tableName, colName, gsiName = blog.TableBlogPost, blog.PostFieldIsPublic, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsAttrTypeNumber}},
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsKeyTypePartition}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}

	// blog comment
	tableName, colName, gsiName = blog.TableBlogComment, blog.CommentFieldOwnerId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsAttrTypeString}},
		[]prom.AwsDynamodbNameAndType{{Name: colName, Type: prom.AwsKeyTypePartition}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}
	tableName, colName, gsiName = blog.TableBlogComment, blog.CommentFieldPostId+"_"+blog.CommentFieldParentId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: blog.CommentFieldPostId, Type: prom.AwsAttrTypeString}, {Name: blog.CommentFieldParentId, Type: prom.AwsAttrTypeString}},
		[]prom.AwsDynamodbNameAndType{{Name: blog.CommentFieldPostId, Type: prom.AwsKeyTypePartition}, {Name: blog.CommentFieldParentId, Type: prom.AwsKeyTypeSort}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}

	// blog vote
	tableName, colName, gsiName = blog.TableBlogVote, blog.VoteFieldOwnerId+"_"+blog.VoteFieldTargetId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: blog.VoteFieldOwnerId, Type: prom.AwsAttrTypeString}, {Name: blog.VoteFieldTargetId, Type: prom.AwsAttrTypeString}},
		[]prom.AwsDynamodbNameAndType{{Name: blog.VoteFieldOwnerId, Type: prom.AwsKeyTypePartition}, {Name: blog.VoteFieldTargetId, Type: prom.AwsKeyTypeSort}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}
	tableName, colName, gsiName = blog.TableBlogVote, blog.VoteFieldTargetId+"_"+blog.VoteFieldValue, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]prom.AwsDynamodbNameAndType{{Name: blog.VoteFieldTargetId, Type: prom.AwsAttrTypeString}, {Name: blog.VoteFieldValue, Type: prom.AwsAttrTypeNumber}},
		[]prom.AwsDynamodbNameAndType{{Name: blog.VoteFieldTargetId, Type: prom.AwsKeyTypePartition}, {Name: blog.VoteFieldValue, Type: prom.AwsKeyTypeSort}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}
}

func _createMongoCollections(mc *prom.MongoConnect) {
	if err := henge.InitMongoCollection(mc, user.TableUser); err != nil {
		log.Printf("[WARN] creating collection %s (%s): %s\n", user.TableUser, "MongoDB", err)
	}
	if err := henge.InitMongoCollection(mc, blog.TableBlogPost); err != nil {
		log.Printf("[WARN] creating collection %s (%s): %s\n", blog.TableBlogPost, "MongoDB", err)
	}
	if err := henge.InitMongoCollection(mc, blog.TableBlogComment); err != nil {
		log.Printf("[WARN] creating collection %s (%s): %s\n", blog.TableBlogComment, "MongoDB", err)
	}
	if err := henge.InitMongoCollection(mc, blog.TableBlogVote); err != nil {
		log.Printf("[WARN] creating collection %s (%s): %s\n", blog.TableBlogVote, "MongoDB", err)
	}

	unique := true
	nonUnique := false
	var idxName string

	// user
	idxName = "idx_" + user.UserColMaskUid
	if _, err := mc.CreateCollectionIndexes(user.TableUser, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{user.UserColMaskUid, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &unique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", user.TableUser, user.UserColMaskUid, "MongoDB", err)
	}

	// blog post
	idxName = "idx_" + blog.PostColOwnerId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogPost, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.PostColOwnerId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostColOwnerId, "MongoDB", err)
	}
	idxName = "idx_" + blog.PostColIsPublic
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogPost, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.PostColIsPublic, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostColIsPublic, "MongoDB", err)
	}

	// blog comment
	idxName = "idx_" + blog.CommentColOwnerId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogComment, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.CommentColOwnerId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentColOwnerId, "MongoDB", err)
	}
	idxName = "idx_" + blog.CommentColPostId + "_" + blog.CommentColParentId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogComment, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.CommentColPostId, 1},
			{blog.CommentColParentId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentColPostId+":"+blog.CommentColParentId, "MongoDB", err)
	}

	// blog vote
	idxName = "idx_" + blog.VoteColOwnerId + "_" + blog.VoteColTargetId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogVote, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.VoteColOwnerId, 1},
			{blog.VoteColTargetId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &unique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteColOwnerId+":"+blog.VoteColTargetId, "MongoDB", err)
	}
	idxName = "idx_" + blog.VoteColTargetId + "_" + blog.VoteColValue
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogVote, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.VoteColTargetId, 1},
			{blog.VoteColValue, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteColTargetId+":"+blog.VoteColValue, "MongoDB", err)
	}
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))

	// create DB connect instance
	sqlc := _createSqlConnect(dbtype)
	mc := _createMongoConnect(dbtype)
	adc := _createDynamodbConnect(dbtype)
	if sqlc == nil && mc == nil && adc == nil {
		panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
	}

	if sqlc != nil {
		// create database tables
		_createSqlTables(sqlc, dbtype)

		// create DAO instances
		userDaov2 = _createUserDaoSql(sqlc)
		blogPostDaov2 = _createBlogPostDaoSql(sqlc)
		blogCommentDaov2 = _createBlogCommentDaoSql(sqlc)
		blogVoteDaov2 = _createBlogVoteDaoSql(sqlc)
	}
	if adc != nil {
		// create AWS DynamoDB tables
		_createDynamodbTables(adc)

		// create DAO instances
		userDaov2 = _createUserDaoDynamodb(adc)
		blogPostDaov2 = _createBlogPostDaoDynamodb(adc)
		blogCommentDaov2 = _createBlogCommentDaoDynamodb(adc)
		blogVoteDaov2 = _createBlogVoteDaoDynamodb(adc)
	}
	if mc != nil {
		// create MongoDB collections
		_createMongoCollections(mc)

		// create DAO instances
		userDaov2 = _createUserDaoMongo(mc)
		blogPostDaov2 = _createBlogPostDaoMongo(mc)
		blogCommentDaov2 = _createBlogCommentDaoMongo(mc)
		blogVoteDaov2 = _createBlogVoteDaoMongo(mc)
	}

	_initUsers()
	_initBlog()
}

func _initUsers() {
	adminUserId := goapi.AppConfig.GetString("gvabe.init.admin_user_id")
	adminUserPwd := goapi.AppConfig.GetString("gvabe.init.admin_user_pwd")
	adminUserName := goapi.AppConfig.GetString("gvabe.init.admin_user_name")
	if adminUserId == "" {
		log.Printf("[WARN] Admin user-id not found at config [gvabe.init.admin_user_id], will not create admin account")
		return
	}
	if adminUserPwd == "" {
		log.Printf("[INFO] Admin password not found at config [gvabe.init.admin_user_pwd], use default value")
		adminUserPwd = "s3cr3t"
	}
	if adminUserName == "" {
		log.Printf("[INFO] Admin display-name not found at config [gvabe.init.admin_user_name], use default value")
		adminUserName = adminUserId
	}
	adminUser, err := userDaov2.Get(adminUserId)
	if err != nil {
		panic(fmt.Sprintf("error while getting user [%s]: %e", adminUserId, err))
	}
	if adminUser == nil {
		log.Printf("[INFO] Admin user [%s] not found, creating one...", adminUserId)
		adminUser = user.NewUser(goapi.AppVersionNumber, adminUserId, utils.UniqueId())
		adminUser.SetPassword(encryptPassword(adminUserId, adminUserPwd)).SetDisplayName(adminUserName).SetAdmin(true)
		result, err := userDaov2.Create(adminUser)
		if err != nil {
			panic(fmt.Sprintf("error while creating user [%s]: %e", adminUserId, err))
		}
		if !result {
			log.Printf("[ERROR] Cannot create user [%s]", adminUserId)
		}
	}
}

func _initBlog() {
	adminUserId := goapi.AppConfig.GetString("gvabe.init.admin_user_id")
	adminUser, err := userDaov2.Get(adminUserId)
	if err != nil {
		panic(fmt.Sprintf("error while getting user [%s]: %e", adminUserId, err))
	}

	postId := "1"
	introBlogPost, err := blogPostDaov2.Get(postId)
	if err != nil {
		panic(fmt.Sprintf("error while getting blog post [%s]: %e", postId, err))
	}
	if introBlogPost == nil {
		log.Printf("[INFO] Introduction blog post [%s] not found, creating one...", postId)
		appName := goapi.AppConfig.GetString("app.name")
		title := "Welcome to " + appName + " v" + goapi.AppVersion
		content := `This is the introduction blog post. It will quickly introduce highlighted features.

**Manage your blog**

You can create, edit or delete your blog posts by accessing **_My Blog_** link on the menu.
Furthermore, you can quickly create a new blog post from **_Create Blog Post_** link.

Blog content supports <a href="https://en.wikipedia.org/wiki/Markdown" target="_blank">Markdown</a> syntax.

**Share your blog posts and interact with others**

_Public_ posts are visible to all users for _commenting_ (coming soon) and _voting_.
`
		introBlogPost = blog.NewBlogPost(goapi.AppVersionNumber, adminUser, true, title, content)
		introBlogPost.SetId(postId)
		result, err := blogPostDaov2.Create(introBlogPost)
		if err != nil {
			panic(fmt.Sprintf("error while creating blog post [%s]: %e", postId, err))
		}
		if !result {
			log.Printf("[ERROR] Cannot create blog post [%s]", postId)
		}
	}
}
