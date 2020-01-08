/*
Package gvabe provides backend API for GoVueAdmin Frontend.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package gvabe

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/prom"
	"log"
	"main/src/goapi"
	"main/src/gvabe/bo"
	"main/src/gvabe/bo/group"
	"main/src/gvabe/bo/user"
	"main/src/itineris"
	"runtime"
	"strconv"
	"strings"
)

var (
	groupDao group.GroupDao
	userDao  user.UserDao
)

type MyBootstrapper struct {
	name string
}

var Bootstrapper = &MyBootstrapper{name: "gvabe"}

/*
Bootstrap implements goapi.IBootstrapper.Bootstrap

Bootstrapper usually does:
- register api-handlers with the global ApiRouter
- other initializing work (e.g. creating DAO, initializing database, etc)
*/
func (b *MyBootstrapper) Bootstrap() error {
	go startUpdateSystemInfo()

	initDaos()
	initApiHandlers(goapi.ApiRouter)
	initApiFilters(goapi.ApiRouter)
	return nil
}

func createSqlConnect() *prom.SqlConnect {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		return bo.NewSqliteConnection(dir, dbname)
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func createGroupDao(sqlc *prom.SqlConnect) group.GroupDao {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		return group.NewGroupDaoSqlite(sqlc, bo.TableGroup)
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func createUserDao(sqlc *prom.SqlConnect) user.UserDao {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		return user.NewUserDaoSqlite(sqlc, bo.TableUser)
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func initDaos() {
	sqlc := createSqlConnect()
	group.InitSqliteTableGroup(sqlc, bo.TableGroup)
	user.InitSqliteTableUser(sqlc, bo.TableUser)

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
		systemAdminUser = &user.User{
			Username: systemAdminUsername,
			Password: encryptPassword(systemAdminUsername, pwd),
			Name:     systemAdminName,
			GroupId:  systemGroupId,
		}
		result, err := userDao.Create(systemAdminUser)
		if err != nil {
			panic("error while creating user [" + systemAdminUsername + "]: " + err.Error())
		}
		if !result {
			log.Printf("Cannot create user [%s]", systemAdminUsername)
		}
	}
}

/*
Setup API filters: application register its api-handlers by calling router.SetHandler(apiName, apiHandlerFunc)

    - api-handler function must has the following signature: func (itineris.ApiContext, itineris.ApiAuth, itineris.ApiParams) *itineris.ApiResult
*/
func initApiFilters(apiRouter *itineris.ApiRouter) {
	var apiFilter itineris.IApiFilter = nil
	// appName := goapi.AppConfig.GetString("app.name")
	// appVersion := goapi.AppConfig.GetString("app.version")

	// filters are LIFO:
	// - request goes through the last filter to the first one
	// - response goes through the first filter to the last one
	// suggested order of filters:
	// - Request logger should be the last one to capture full request/response

	// apiFilter = itineris.NewAddPerfInfoFilter(goapi.ApiRouter, apiFilter)
	// apiFilter = itineris.NewLoggingFilter(goapi.ApiRouter, apiFilter, itineris.NewWriterPerfLogger(os.Stderr, appName, appVersion))
	apiFilter = itineris.NewAuthenticationFilter(goapi.ApiRouter, apiFilter, &GVAFEApiAuthenticator{})
	// apiFilter = itineris.NewLoggingFilter(goapi.ApiRouter, apiFilter, itineris.NewWriterRequestLogger(os.Stdout, appName, appVersion))

	apiRouter.SetApiFilter(apiFilter)
}

/*
GVAFEApiAuthenticator is an "IApiAuthenticator" which checks:

	- AppId must be "$shortname$_fe"
	- AccessToken must be valid (allocated and active)s
*/
type GVAFEApiAuthenticator struct {
}

/*
Authenticate implements IApiAuthenticator.Authenticate.
*/
func (a *GVAFEApiAuthenticator) Authenticate(ctx *itineris.ApiContext, auth *itineris.ApiAuth) bool {
	if "$shortname$_fe" != auth.GetAppId() {
		return false
	}
	if ctx.GetApiName() != "login" {
		// TODO verify token
	}
	return true
}

/*----------------------------------------------------------------------*/

/*
Setup API handlers: application register its api-handlers by calling router.SetHandler(apiName, apiHandlerFunc)

    - api-handler function must has the following signature: func (itineris.ApiContext, itineris.ApiAuth, itineris.ApiParams) *itineris.ApiResult
*/
func initApiHandlers(router *itineris.ApiRouter) {
	router.SetHandler("info", apiInfo)
	router.SetHandler("login", apiLogin)
	router.SetHandler("systemInfo", apiSystemInfo)
}

/*
API handler "info"
*/
func apiInfo(ctx *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	appInfo := map[string]interface{}{
		"name":        goapi.AppConfig.GetString("app.name"),
		"shortname":   goapi.AppConfig.GetString("app.shortname"),
		"version":     goapi.AppConfig.GetString("app.version"),
		"description": goapi.AppConfig.GetString("app.desc"),
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	result := map[string]interface{}{
		"app": appInfo,
		"memory": map[string]interface{}{
			"alloc":     m.Alloc,
			"alloc_str": strconv.FormatFloat(float64(m.Alloc)/1024.0/1024.0, 'f', 1, 64) + " MiB",
			"sys":       m.Sys,
			"sys_str":   strconv.FormatFloat(float64(m.Sys)/1024.0/1024.0, 'f', 1, 64) + " MiB",
			"gc":        m.NumGC,
		},
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
}

/*
API handler "login"
*/
func apiLogin(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	username, _ := params.GetParamAsType("username", reddo.TypeString)
	if username == nil || username == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty username")
	}
	password, _ := params.GetParamAsType("password", reddo.TypeString)
	if password == nil || password == "" {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("login failed")
	}
	user, err := userDao.Get(username.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
	if user == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("login failed")
	}
	encPwd := encryptPassword(user.Username, password.(string))
	if encPwd != user.Password {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("login failed")
	}
	js, _ := json.Marshal(map[string]interface{}{"username": user.Username, "group_id": user.GroupId})
	token := base64.StdEncoding.EncodeToString(js)
	return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"token": token})
}

/*
API handler "systemInfo"
*/
func apiSystemInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	data := lastSystemInfo()
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}
