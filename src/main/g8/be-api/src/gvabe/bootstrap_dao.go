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
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func initDaos() {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	sqlc := _createSqlConnect(dbtype)

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