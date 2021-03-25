package gvabe

import (
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/henge"
	"github.com/btnguyen2k/prom"

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

func _createUserDao(sqlc *prom.SqlConnect) user.UserDao {
	return user.NewUserDaoSql(sqlc, user.TableUser)
}

func _createBlogPostDao(sqlc *prom.SqlConnect) blog.BlogPostDao {
	return blog.NewBlogPostDaoSql(sqlc, blog.TableBlogPost)
}

func _createBlogCommentDao(sqlc *prom.SqlConnect) blog.BlogCommentDao {
	return blog.NewBlogCommentDaoSql(sqlc, blog.TableBlogComment)
}

func _createBlogVoteDao(sqlc *prom.SqlConnect) blog.BlogVoteDao {
	return blog.NewBlogVoteDaoSql(sqlc, blog.TableBlogVote)
}

func _createSqlTables(sqlc *prom.SqlConnect, dbtype string) {
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, user.TableUser, map[string]string{user.UserCol_MaskUid: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, blog.TableBlogPost, map[string]string{
			blog.PostCol_OwnerId: "VARCHAR(32)", blog.PostCol_IsPublic: "INT"})
		henge.InitSqliteTable(sqlc, blog.TableBlogComment, map[string]string{
			blog.CommentCol_OwnerId: "VARCHAR(32)", blog.CommentCol_PostId: "VARCHAR(32)", blog.CommentCol_ParentId: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, blog.TableBlogVote, map[string]string{
			blog.VoteCol_OwnerId: "VARCHAR(32)", blog.VoteCol_TargetId: "VARCHAR(32)", blog.VoteCol_Value: "INT"})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, user.TableUser, map[string]string{user.UserCol_MaskUid: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, blog.TableBlogPost, map[string]string{
			blog.PostCol_OwnerId: "VARCHAR(32)", blog.PostCol_IsPublic: "INT"})
		henge.InitPgsqlTable(sqlc, blog.TableBlogComment, map[string]string{
			blog.CommentCol_OwnerId: "VARCHAR(32)", blog.CommentCol_PostId: "VARCHAR(32)", blog.CommentCol_ParentId: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, blog.TableBlogVote, map[string]string{
			blog.VoteCol_OwnerId: "VARCHAR(32)", blog.VoteCol_TargetId: "VARCHAR(32)", blog.VoteCol_Value: "INT"})
	}

	// user
	henge.CreateIndexSql(sqlc, user.TableUser, true, []string{user.UserCol_MaskUid})

	// blog post
	henge.CreateIndexSql(sqlc, blog.TableBlogPost, false, []string{blog.PostCol_OwnerId})
	henge.CreateIndexSql(sqlc, blog.TableBlogPost, false, []string{blog.PostCol_IsPublic})

	// blog comment
	henge.CreateIndexSql(sqlc, blog.TableBlogComment, false, []string{blog.CommentCol_OwnerId})
	henge.CreateIndexSql(sqlc, blog.TableBlogComment, false, []string{blog.CommentCol_PostId, blog.CommentCol_ParentId})

	// blog vote
	henge.CreateIndexSql(sqlc, blog.TableBlogVote, true, []string{blog.VoteCol_OwnerId, blog.VoteCol_TargetId})
	henge.CreateIndexSql(sqlc, blog.TableBlogVote, false, []string{blog.VoteCol_TargetId, blog.VoteCol_Value})
}

func initDaos() {
	// create SqlConnect instance
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype) // only SQL-based datastore is supported
	if sqlc == nil {
		panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
	}

	// create database tables (assuming SQL-based datastore)
	_createSqlTables(sqlc, dbtype)

	// create DAO instances
	userDaov2 = _createUserDao(sqlc)
	blogPostDaov2 = _createBlogPostDao(sqlc)
	blogCommentDaov2 = _createBlogCommentDao(sqlc)
	blogVoteDaov2 = _createBlogVoteDao(sqlc)

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
