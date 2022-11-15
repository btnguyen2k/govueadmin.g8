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
	promdynamodb "github.com/btnguyen2k/prom/dynamodb"
	prommongo "github.com/btnguyen2k/prom/mongo"
	promsql "github.com/btnguyen2k/prom/sql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"main/src/goapi"
	"main/src/gvabe/bov2/blog"
	"main/src/gvabe/bov2/user"
	"main/src/utils"

	_ "github.com/btnguyen2k/gocosmos"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/mattn/go-sqlite3"
)

func _createDynamodbConnect(dbtype string) *promdynamodb.AwsDynamodbConnect {
	var adc *promdynamodb.AwsDynamodbConnect = nil
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
		adc, err = promdynamodb.NewAwsDynamodbConnect(cfg, nil, nil, 10000)
	}
	if err != nil {
		panic(err)
	}
	return adc
}

func _createSqlConnect(dbtype string) *promsql.SqlConnect {
	timezone := goapi.AppConfig.GetString("timezone")
	var sqlc *promsql.SqlConnect = nil
	var err error
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		sqlc, err = henge.NewSqliteConnection(dir, dbname, timezone, "sqlite3", 10000, nil)
	case "pg", "pgsql", "postgres", "postgresql":
		url := goapi.AppConfig.GetString("gvabe.db.pgsql.url")
		sqlc, err = henge.NewPgsqlConnection(url, timezone, "pgx", 10000, nil)
	case "cosmos", "cosmosdb":
		url := goapi.AppConfig.GetString("gvabe.db.cosmosdb.url")
		sqlc, err = henge.NewCosmosdbConnection(url, timezone, "gocosmos", 10000, nil)
	}
	if err != nil {
		panic(err)
	}
	return sqlc
}

func _createMongoConnect(dbtype string) *prommongo.MongoConnect {
	var mc *prommongo.MongoConnect = nil
	var err error
	switch dbtype {
	case "mongo", "mongodb":
		db := goapi.AppConfig.GetString("gvabe.db.mongodb.db")
		url := goapi.AppConfig.GetString("gvabe.db.mongodb.url")
		mc, err = prommongo.NewMongoConnect(url, db, 10000)
	}
	if err != nil {
		panic(err)
	}
	return mc
}

func _createUserDaoSql(sqlc *promsql.SqlConnect) user.UserDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return user.NewUserDaoCosmosdb(sqlc, user.TableUser, true)
	}
	return user.NewUserDaoSql(sqlc, user.TableUser, true)
}
func _createUserDaoDynamodb(adc *promdynamodb.AwsDynamodbConnect) user.UserDao {
	return user.NewUserDaoDynamodb(adc, user.TableUser)
}
func _createUserDaoMongo(mc *prommongo.MongoConnect) user.UserDao {
	url := mc.GetUrl()
	return user.NewUserDaoMongo(mc, user.TableUser, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogPostDaoSql(sqlc *promsql.SqlConnect) blog.BlogPostDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return blog.NewBlogPostDaoCosmosdb(sqlc, blog.TableBlogPost, true)
	}
	return blog.NewBlogPostDaoSql(sqlc, blog.TableBlogPost, true)
}
func _createBlogPostDaoDynamodb(adc *promdynamodb.AwsDynamodbConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoDynamodb(adc, blog.TableBlogPost)
}
func _createBlogPostDaoMongo(mc *prommongo.MongoConnect) blog.BlogPostDao {
	url := mc.GetUrl()
	return blog.NewBlogPostDaoMongo(mc, blog.TableBlogPost, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogCommentDaoSql(sqlc *promsql.SqlConnect) blog.BlogCommentDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return blog.NewBlogCommentDaoCosmosdb(sqlc, blog.TableBlogComment, true)
	}
	return blog.NewBlogCommentDaoSql(sqlc, blog.TableBlogComment, true)
}
func _createBlogCommentDaoDynamodb(adc *promdynamodb.AwsDynamodbConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoDynamodb(adc, blog.TableBlogComment)
}
func _createBlogCommentDaoMongo(mc *prommongo.MongoConnect) blog.BlogCommentDao {
	url := mc.GetUrl()
	return blog.NewBlogCommentDaoMongo(mc, blog.TableBlogComment, strings.Index(url, "replicaSet=") >= 0)
}

func _createBlogVoteDaoSql(sqlc *promsql.SqlConnect) blog.BlogVoteDao {
	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return blog.NewBlogVoteDaoCosmosdb(sqlc, blog.TableBlogVote, true)
	}
	return blog.NewBlogVoteDaoSql(sqlc, blog.TableBlogVote, true)
}
func _createBlogVoteDaoDynamodb(adc *promdynamodb.AwsDynamodbConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoDynamodb(adc, blog.TableBlogVote)
}
func _createBlogVoteDaoMongo(mc *prommongo.MongoConnect) blog.BlogVoteDao {
	url := mc.GetUrl()
	return blog.NewBlogVoteDaoMongo(mc, blog.TableBlogVote, strings.Index(url, "replicaSet=") >= 0)
}

var _sqliteTableSchema = map[string]map[string]string{
	user.TableUser:        {user.UserColMaskUid: "VARCHAR(32)"},
	blog.TableBlogPost:    {blog.PostColOwnerId: "VARCHAR(32)", blog.PostColIsPublic: "INT"},
	blog.TableBlogComment: {blog.CommentColOwnerId: "VARCHAR(32)", blog.CommentColPostId: "VARCHAR(32)", blog.CommentColParentId: "VARCHAR(32)"},
	blog.TableBlogVote:    {blog.VoteColOwnerId: "VARCHAR(32)", blog.VoteColTargetId: "VARCHAR(32)", blog.VoteColValue: "INT"},
}

var _pgsqlTableSchema = map[string]map[string]string{
	user.TableUser:        {user.UserColMaskUid: "VARCHAR(32)"},
	blog.TableBlogPost:    {blog.PostColOwnerId: "VARCHAR(32)", blog.PostColIsPublic: "INT"},
	blog.TableBlogComment: {blog.CommentColOwnerId: "VARCHAR(32)", blog.CommentColPostId: "VARCHAR(32)", blog.CommentColParentId: "VARCHAR(32)"},
	blog.TableBlogVote:    {blog.VoteColOwnerId: "VARCHAR(32)", blog.VoteColTargetId: "VARCHAR(32)", blog.VoteColValue: "INT"},
}

var _cosmosdbTableSpec = map[string]*henge.CosmosdbCollectionSpec{
	user.TableUser:        {Pk: henge.CosmosdbColId, Uk: [][]string{{"/" + user.UserFieldMaskId}}},
	blog.TableBlogPost:    {Pk: henge.CosmosdbColId},
	blog.TableBlogComment: {Pk: henge.CosmosdbColId},
	blog.TableBlogVote:    {Pk: henge.CosmosdbColId, Uk: [][]string{{"/" + blog.VoteFieldOwnerId, "/" + blog.VoteFieldTargetId}}},
}

func _createSqlTables(sqlc *promsql.SqlConnect, dbtype string) {
	switch sqlc.GetDbFlavor() {
	case promsql.FlavorSqlite:
		for tbl, schema := range _sqliteTableSchema {
			if err := henge.InitSqliteTable(sqlc, tbl, schema); err != nil {
				log.Printf("[WARN] creating table %s (%s): %s\n", tbl, dbtype, err)
			}
		}
	case promsql.FlavorPgSql:
		for tbl, schema := range _pgsqlTableSchema {
			if err := henge.InitSqliteTable(sqlc, tbl, schema); err != nil {
				log.Printf("[WARN] creating table %s (%s): %s\n", tbl, dbtype, err)
			}
		}
	case promsql.FlavorCosmosDb:
		for tbl, spec := range _cosmosdbTableSpec {
			if err := henge.InitCosmosdbCollection(sqlc, tbl, spec); err != nil {
				log.Printf("[WARN] creating table %s (%s): %s\n", tbl, dbtype, err)
			}
		}
	}

	if sqlc.GetDbFlavor() == promsql.FlavorCosmosDb {
		return
	}

	// user
	if err := henge.CreateIndexSql(sqlc, user.TableUser, true, []string{user.UserColMaskUid}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", user.TableUser, user.UserColMaskUid, dbtype, err)
	}

	// blog post
	if err := henge.CreateIndexSql(sqlc, blog.TableBlogPost, false, []string{blog.PostColOwnerId}); err != nil {
		log.Printf("[WARN] creating table index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostColOwnerId, dbtype, err)
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

func _dynamodbWaitforGSI(adc *promdynamodb.AwsDynamodbConnect, table, gsi string, timeout time.Duration) error {
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

func _createDynamodbTables(adc *promdynamodb.AwsDynamodbConnect) {
	if err := blog.InitBlogCommentTableDynamodb(adc, blog.TableBlogComment); err != nil {
		panic(err)
	}
	if err := blog.InitBlogPostTableDynamodb(adc, blog.TableBlogPost); err != nil {
		panic(err)
	}
	if err := blog.InitBlogVoteTableDynamodb(adc, blog.TableBlogVote); err != nil {
		panic(err)
	}

	spec := &henge.DynamodbTablesSpec{MainTableRcu: 2, MainTableWcu: 1, CreateUidxTable: true, UidxTableRcu: 2, UidxTableWcu: 1}
	if err := henge.InitDynamodbTables(adc, user.TableUser, spec); err != nil {
		log.Printf("[WARN] creating tableName %s (%s): %s\n", user.TableUser, "DynamoDB", err)
	}

	var tableName, gsiName, colName string

	// user
	tableName, colName, gsiName = user.TableUser, user.UserFieldMaskId, "gsi_"+colName
	if err := adc.CreateGlobalSecondaryIndex(nil, tableName, gsiName, 2, 1,
		[]promdynamodb.AwsDynamodbNameAndType{{Name: colName, Type: promdynamodb.AwsAttrTypeString}},
		[]promdynamodb.AwsDynamodbNameAndType{{Name: colName, Type: promdynamodb.AwsKeyTypePartition}}); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	} else if err := _dynamodbWaitforGSI(adc, tableName, gsiName, 10*time.Second); err != nil {
		log.Printf("[WARN] creating GSI %s/%s (%s): %s\n", tableName, colName, "DynamoDB", err)
	}
}

func _createMongoCollections(mc *prommongo.MongoConnect) {
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
	idxName = "idx_" + user.UserFieldMaskId
	if _, err := mc.CreateCollectionIndexes(user.TableUser, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{user.UserFieldMaskId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &unique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", user.TableUser, user.UserFieldMaskId, "MongoDB", err)
	}

	// blog post
	idxName = "idx_" + blog.PostFieldOwnerId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogPost, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.PostFieldOwnerId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogPost, blog.PostFieldOwnerId, "MongoDB", err)
	}

	// blog comment
	idxName = "idx_" + blog.CommentFieldOwnerId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogComment, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.CommentFieldOwnerId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentFieldOwnerId, "MongoDB", err)
	}
	idxName = "idx_" + blog.CommentFieldPostId + "_" + blog.CommentFieldParentId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogComment, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.CommentFieldPostId, 1},
			{blog.CommentFieldParentId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogComment, blog.CommentFieldPostId+":"+blog.CommentFieldParentId, "MongoDB", err)
	}

	// blog vote
	idxName = "idx_" + blog.VoteFieldOwnerId + "_" + blog.VoteFieldTargetId
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogVote, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.VoteFieldOwnerId, 1},
			{blog.VoteFieldTargetId, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &unique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteFieldOwnerId+":"+blog.VoteFieldTargetId, "MongoDB", err)
	}
	idxName = "idx_" + blog.VoteFieldTargetId + "_" + blog.VoteFieldValue
	if _, err := mc.CreateCollectionIndexes(blog.TableBlogVote, []interface{}{mongo.IndexModel{
		Keys: bson.D{
			{blog.VoteFieldTargetId, 1},
			{blog.VoteFieldValue, 1},
		},
		Options: &options.IndexOptions{
			Name:   &idxName,
			Unique: &nonUnique,
		},
	}}); err != nil {
		log.Printf("[WARN] creating collection index %s/%s (%s): %s\n", blog.TableBlogVote, blog.VoteFieldTargetId+":"+blog.VoteFieldValue, "MongoDB", err)
	}
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	if DEBUG {
		log.Printf("[DEUBG] db-type: %s", dbtype)
	}

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
		adminUser = user.NewUser(goapi.AppVersionNumber, adminUserId, utils.UniqueId())
		adminUser.SetPassword(encryptPassword(adminUserId, adminUserPwd)).SetDisplayName(adminUserName).SetAdmin(true)
		log.Printf("[INFO] Admin user [%s] not found, creating one...(%s)", adminUserId, adminUser.GetMaskId())
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
