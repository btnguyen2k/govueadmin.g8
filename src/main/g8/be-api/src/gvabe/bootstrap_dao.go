package gvabe

import (
	"fmt"
	"log"
	"strings"

	"github.com/btnguyen2k/prom"

	"main/src/goapi"
	"main/src/gvabe/bo"
	"main/src/gvabe/bo/group"
	"main/src/gvabe/bo/user"
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

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)
	if sqlc == nil {
		panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
	}
	switch dbtype {
	case "sqlite":
		henge.InitSqliteTable(sqlc, userv2.TableUser, map[string]string{userv2.UserCol_MaskUid: "VARCHAR(32)"})
		henge.CreateIndex(sqlc, userv2.TableUser, true, []string{userv2.UserCol_MaskUid})
	case "pg", "pgsql", "postgres", "postgresql":
		henge.InitPgsqlTable(sqlc, userv2.TableUser, map[string]string{userv2.UserCol_MaskUid: "VARCHAR(32)"})
		henge.CreateIndex(sqlc, userv2.TableUser, true, []string{userv2.UserCol_MaskUid})
	}
	userDaov2 = _createUserDao(sqlc)
	_initUsers()

	switch dbtype {
	case "sqlite":
		group.InitSqliteTableGroup(sqlc, bo.TableGroup)
		user.InitSqliteTableUser(sqlc, bo.TableUser)
	case "pg", "pgsql", "postgres", "postgresql":
		group.InitPgsqlTableGroup(sqlc, bo.TableGroup)
		user.InitPgsqlTableUser(sqlc, bo.TableUser)
	}

	groupDao = createGroupDao(sqlc)
	systemGroup, err := groupDao.Get(systemGroupId)
	if err != nil {
		panic("error while getting group [" + systemGroupId + "]: " + err.Error())
	}
	if systemGroup == nil {
		log.Printf("System group [%s] not found, creating one...", systemGroupId)
		result, err := groupDao.Create(&group.Group{Id: systemGroupId, Name: "System User Group"})
		if err != nil {
			panic("error while creating group [" + systemGroupId + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create group [%s]", systemGroupId)
		}
	}

	userDao = createUserDao(sqlc)
	systemAdminUser, err := userDao.Get(systemAdminUsername)
	if err != nil {
		panic("error while getting user [" + systemAdminUsername + "]: " + err.Error())
	}
	if systemAdminUser == nil {
		pwd := "s3cr3t"
		log.Printf("System admin user [%s] not found, creating one with password [%s]...", systemAdminUsername, pwd)
		systemAdminUser = user.NewUserBo(systemAdminUsername, "").
			SetPassword(encryptPassword(systemAdminUsername, pwd)).
			SetName(systemAdminName).
			SetGroupId(systemGroupId).
			SetAesKey(utils.RandomString(16))
		result, err := userDao.Create(systemAdminUser)
		if err != nil {
			panic("error while creating user [" + systemAdminUsername + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", systemAdminUsername)
		}
	}
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
