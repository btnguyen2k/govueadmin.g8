package gvabe

import (
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/prom"

	"main/src/goapi"
	blogv2 "main/src/gvabe/bov2/blog"
	userv2 "main/src/gvabe/bov2/user"
	"main/src/henge"
	"main/src/utils"
)

func _createSqlConnect(dbtype string) *prom.SqlConnect {
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		return henge.NewSqliteConnection(dir, dbname)
	case "pg", "pgsql", "postgres", "postgresql":
		url := goapi.AppConfig.GetString("gvabe.db.pgsql.url")
		return henge.NewPgsqlConnection(url, goapi.AppConfig.GetString("timezone"))
	}
	return nil
}

func _createUserDao(sqlc *prom.SqlConnect) userv2.UserDao {
	return userv2.NewUserDaoSql(sqlc, userv2.TableUser)
}

func _createBlogPostDao(sqlc *prom.SqlConnect) blogv2.BlogPostDao {
	return blogv2.NewBlogPostDaoSql(sqlc, blogv2.TableBlogPost)
}

func _createBlogCommentDao(sqlc *prom.SqlConnect) blogv2.BlogCommentDao {
	return blogv2.NewBlogCommentDaoSql(sqlc, blogv2.TableBlogComment)
}

func _createBlogVoteDao(sqlc *prom.SqlConnect) blogv2.BlogVoteDao {
	return blogv2.NewBlogVoteDaoSql(sqlc, blogv2.TableBlogVote)
}

func _createSqlTables(sqlc *prom.SqlConnect, dbtype string) {
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, userv2.TableUser, map[string]string{userv2.UserCol_MaskUid: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, blogv2.TableBlogPost, map[string]string{
			blogv2.PostCol_OwnerId: "VARCHAR(32)", blogv2.PostCol_IsPublic: "INT"})
		henge.InitSqliteTable(sqlc, blogv2.TableBlogComment, map[string]string{
			blogv2.CommentCol_OwnerId: "VARCHAR(32)", blogv2.CommentCol_PostId: "VARCHAR(32)", blogv2.CommentCol_ParentId: "VARCHAR(32)"})
		henge.InitSqliteTable(sqlc, blogv2.TableBlogVote, map[string]string{
			blogv2.VoteCol_OwnerId: "VARCHAR(32)", blogv2.VoteCol_TargetId: "VARCHAR(32)", blogv2.VoteCol_Value: "INT"})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, userv2.TableUser, map[string]string{userv2.UserCol_MaskUid: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, blogv2.TableBlogPost, map[string]string{
			blogv2.PostCol_OwnerId: "VARCHAR(32)", blogv2.PostCol_IsPublic: "INT"})
		henge.InitPgsqlTable(sqlc, blogv2.TableBlogComment, map[string]string{
			blogv2.CommentCol_OwnerId: "VARCHAR(32)", blogv2.CommentCol_PostId: "VARCHAR(32)", blogv2.CommentCol_ParentId: "VARCHAR(32)"})
		henge.InitPgsqlTable(sqlc, blogv2.TableBlogVote, map[string]string{
			blogv2.VoteCol_OwnerId: "VARCHAR(32)", blogv2.VoteCol_TargetId: "VARCHAR(32)", blogv2.VoteCol_Value: "INT"})
	}

	// user
	henge.CreateIndex(sqlc, userv2.TableUser, true, []string{userv2.UserCol_MaskUid})

	// blog post
	henge.CreateIndex(sqlc, blogv2.TableBlogPost, false, []string{blogv2.PostCol_OwnerId})
	henge.CreateIndex(sqlc, blogv2.TableBlogPost, false, []string{blogv2.PostCol_IsPublic})

	// blog comment
	henge.CreateIndex(sqlc, blogv2.TableBlogComment, false, []string{blogv2.CommentCol_OwnerId})
	henge.CreateIndex(sqlc, blogv2.TableBlogComment, false, []string{blogv2.CommentCol_PostId, blogv2.CommentCol_ParentId})

	// blog vote
	henge.CreateIndex(sqlc, blogv2.TableBlogVote, true, []string{blogv2.VoteCol_OwnerId, blogv2.VoteCol_TargetId})
	henge.CreateIndex(sqlc, blogv2.TableBlogVote, false, []string{blogv2.VoteCol_TargetId, blogv2.VoteCol_Value})
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
		adminUser = userv2.NewUser(goapi.AppVersionNumber, adminUserId, utils.UniqueId())
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
		introBlogPost = blogv2.NewBlogPost(goapi.AppVersionNumber, adminUser, true, title, content)
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
