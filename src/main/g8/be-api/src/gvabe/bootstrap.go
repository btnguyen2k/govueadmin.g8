/*
Package gvabe provides backend API for GoVueAdmin Frontend.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package gvabe

import (
	"fmt"
	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/prom"
	"log"
	"main/src/goapi"
	"main/src/gvabe/bo"
	"main/src/gvabe/bo/group"
	"main/src/gvabe/bo/user"
	"main/src/itineris"
	"main/src/utils"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
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
	apiFilter = &GVAFEAuthenticationFilter{BaseApiFilter: &itineris.BaseApiFilter{ApiRouter: apiRouter, NextFilter: apiFilter}}
	// apiFilter = itineris.NewLoggingFilter(goapi.ApiRouter, apiFilter, itineris.NewWriterRequestLogger(os.Stdout, appName, appVersion))

	apiRouter.SetApiFilter(apiFilter)
}

/*
GVAFEAuthenticationFilter performs authentication check before calling API and issues new access token if existing one is about to expire.

	- AppId must be "$shortname$_fe"
	- AccessToken must be valid (allocated and active)
*/
type GVAFEAuthenticationFilter struct {
	*itineris.BaseApiFilter
}

/*
Call implements IApiFilter.Call

This function first authenticates API call. If successful and login session is about to expire,
this function renews the access token by generating a new token and returning it in result's extra field.
The returned access token is in the following format: "username:login_token:expiry" (without quotes).
*/
func (f *GVAFEAuthenticationFilter) Call(handler itineris.IApiHandler, ctx *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	authed, username, loginToken := f.authenticate(ctx, auth)
	if !authed {
		return itineris.ResultNoPermission
	}
	if f.NextFilter != nil {
		return f.NextFilter.Call(handler, ctx, auth, params)
	}
	result := handler(ctx, auth, params)
	if username != "" && loginToken != "" {
		if loginData, err := decodeLoginToken(username, loginToken); err == nil && loginData != nil {
			if expiry, err := reddo.ToInt(loginData[loginAttrExpiry]); err == nil && expiry-loginSessionNearExpiry < time.Now().Unix() {
				if user, err := userDao.Get(username); err == nil && user != nil {
					if loginToken, err := genLoginToken(user); err == nil {
						expiry := time.Now().Unix() + loginSessionTtl
						result.AddExtraInfo(apiResultExtraAccessToken, username+":"+loginToken+":"+strconv.FormatInt(expiry, 10))
					}
				}
			}
		}
	}
	return result
}

/*
authenticate authenticates an API call.

This function expects auth.access_token in the following format: "username:login_token" (without quotes).

Upon successful authentication, this function returns (true, username, login_token), where username and login_token were generated by apiLogin.
*/
func (f *GVAFEAuthenticationFilter) authenticate(ctx *itineris.ApiContext, auth *itineris.ApiAuth) (bool, string, string) {
	if "$shortname$_fe" != auth.GetAppId() {
		return false, "", ""
	}
	if ctx.GetApiName() != "info" && ctx.GetApiName() != "login" {
		tokens := strings.SplitN(auth.GetAccessToken(), ":", 2)
		if len(tokens) != 2 {
			return false, "", ""
		}
		status, err := verifyLoginToken(tokens[0], tokens[1])
		return err == nil && status == sessionStatusOk, tokens[0], tokens[1]
	}
	return true, "", ""
}

/*----------------------------------------------------------------------*/

/*
Setup API handlers: application register its api-handlers by calling router.SetHandler(apiName, apiHandlerFunc)

    - api-handler function must has the following signature: func (itineris.ApiContext, itineris.ApiAuth, itineris.ApiParams) *itineris.ApiResult
*/
func initApiHandlers(router *itineris.ApiRouter) {
	router.SetHandler("info", apiInfo)
	router.SetHandler("login", apiLogin)
	router.SetHandler("checkLoginToken", apiCheckLoginToken)
	router.SetHandler("systemInfo", apiSystemInfo)

	router.SetHandler("groupList", apiGroupList)
	router.SetHandler("getGroup", apiGetGroup)
	router.SetHandler("createGroup", apiCreateGroup)
	router.SetHandler("updateGroup", apiUpdateGroup)

	router.SetHandler("userList", apiUserList)
}

// API handler "info"
func apiInfo(_ *itineris.ApiContext, auth *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
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
apiLogin handles API call "login".
Upon login successfully, this API returns the following map:

	{
		"uid": username of logged-in user,
		"token": login token, used for latter authentication (e.g. apiCheckLoginToken),
		"expiry": session's expiry, in UNIX timestamp (seconds),
	}
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
	encPwd := encryptPassword(user.GetUsername(), password.(string))
	if encPwd != user.GetPassword() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage("login failed")
	}
	if token, err := genLoginToken(user); err == nil {
		t := time.Now()
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{
			"uid":    user.GetUsername(),
			"token":  token,
			"expiry": t.Unix() + loginSessionTtl,
		}).AddExtraInfo(apiResultExtraAccessToken, user.GetUsername()+":"+token)
	} else {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
}

/*
apiLogin handles API call "checkLoginToken".
This API expects an input map:

	{
		"uid": username of logged-in user (returned by apiLogin),
		"token": login token (returned by apiLogin),
	}

Upon successful, this API return itineris.StatusOk
*/
func apiCheckLoginToken(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	token, _ := params.GetParamAsType("token", reddo.TypeString)
	if token == nil || token == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty token")
	}
	username, _ := params.GetParamAsType("uid", reddo.TypeString)
	if username == nil || username == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty username")
	}
	if status, err := verifyLoginToken(username.(string), token.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(err.Error())
	} else {
		switch status {
		case sessionStatusOk:
			return itineris.NewApiResult(itineris.StatusOk)
		case sessionStatusInvalid, sessionStatusUserNotFound, sessionStatusExpired:
			return itineris.NewApiResult(itineris.StatusNoPermission)
		default:
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage("Unknown error, status = " + strconv.Itoa(status))
		}
	}
}

// API handler "systemInfo"
func apiSystemInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	data := lastSystemInfo()
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// API handler "groupList"
func apiGroupList(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	groupList, err := groupDao.GetAll()
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
	data := make([]map[string]interface{}, 0)
	for _, g := range groupList {
		data = append(data, map[string]interface{}{"id": g.Id, "name": g.Name})
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// API handler "getGroup"
func apiGetGroup(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("Group [%s] not found", id))
	}
	if group, err := groupDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if group == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("Group [%s] not found", id))
	} else {
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": group.Id, "name": group.Name})
	}
}

// API handler "updateGroup"
func apiUpdateGroup(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("Group [%s] not found", id))
	}
	if group, err := groupDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if group == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("Group [%s] not found", id))
	} else {
		// TODO check current user's permission

		name, _ := params.GetParamAsType("name", reddo.TypeString)
		if name == nil || strings.TrimSpace(name.(string)) == "" {
			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
		}
		group.Name = strings.TrimSpace(name.(string))
		if ok, err := groupDao.Update(group); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		} else if !ok {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Group [%s] has not been updated", id.(string)))
		}
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"id": group.Id, "name": group.Name})
	}
}

// API handler "createGroup"
func apiCreateGroup(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	id, _ := params.GetParamAsType("id", reddo.TypeString)
	if id == nil || strings.TrimSpace(id.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [id]")
	}
	id = strings.TrimSpace(strings.ToLower(id.(string)))
	if !regexp.MustCompile("^[0-9a-z_]+$").Match([]byte(id.(string))) {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Invalid value for parameter [id]")
	}

	name, _ := params.GetParamAsType("name", reddo.TypeString)
	if name == nil || strings.TrimSpace(name.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
	}
	name = strings.TrimSpace(name.(string))

	if group, err := groupDao.Get(id.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if group != nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Group [%s] already existed", id))
	}
	group := &group.Group{
		Id:   strings.TrimSpace(strings.ToLower(id.(string))),
		Name: strings.TrimSpace(name.(string)),
	}
	if ok, err := groupDao.Create(group); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Group [%s] has not been created", id))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(group)
}

// API handler "userList"
func apiUserList(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	data := lastSystemInfo()
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}
