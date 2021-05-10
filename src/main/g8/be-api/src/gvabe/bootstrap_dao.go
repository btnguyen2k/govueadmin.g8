package gvabe

import (
	"fmt"
	"log"
	"strings"

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
	return user.NewUserDaoSql(sqlc, user.TableUser)
}
func _createUserDaoMongo(mc *prom.MongoConnect) user.UserDao {
	return user.NewUserDaoMongo(mc, user.TableUser)
}

func _createBlogPostDaoSql(sqlc *prom.SqlConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoSql(sqlc, blog.TableBlogPost)
}
func _createBlogPostDaoMongo(mc *prom.MongoConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoMongo(mc, blog.TableBlogPost)
}

func _createBlogCommentDaoSql(sqlc *prom.SqlConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoSql(sqlc, blog.TableBlogComment)
}
func _createBlogCommentDaoMongo(mc *prom.MongoConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoMongo(mc, blog.TableBlogComment)
}

func _createBlogVoteDaoSql(sqlc *prom.SqlConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoSql(sqlc, blog.TableBlogVote)
}
func _createBlogVoteDaoMongo(mc *prom.MongoConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoMongo(mc, blog.TableBlogVote)
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

	// create SqlConnect instance
	sqlc := _createSqlConnect(dbtype)
	mc := _createMongoConnect(dbtype)
	if sqlc == nil && mc == nil {
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
