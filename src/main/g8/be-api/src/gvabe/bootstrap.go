/*
Package gvabe provides backend API for GoVueAdmin Frontend.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.1.0
*/
package gvabe

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/prom"

	"main/src/goapi"
	"main/src/gvabe/bo"
	"main/src/gvabe/bo/group"
	"main/src/gvabe/bo/user"
	blogv2 "main/src/gvabe/bov2/blog"
	userv2 "main/src/gvabe/bov2/user"
	"main/src/itineris"
	"main/src/utils"
)

var (
	userDaov2        userv2.UserDao
	blogPostDaov2    blogv2.BlogPostDao
	blogCommentDaov2 blogv2.BlogCommentDao
	blogVoteDaov2    blogv2.BlogVoteDao

	groupDao group.GroupDao
	userDao  user.UserDao
)

// MyBootstrapper implements goapi.IBootstrapper
type MyBootstrapper struct {
	name string
}

var Bootstrapper = &MyBootstrapper{name: "gvabe"}

/*
Bootstrap implements goapi.IBootstrapper.Bootstrap

Bootstrapper usually does the following:
- register api-handlers with the global ApiRouter
- other initializing work (e.g. creating DAO, initializing database, etc)
*/
func (b *MyBootstrapper) Bootstrap() error {
	if os.Getenv("DEBUG") != "" {
		DEBUG = true
	}
	go startUpdateSystemInfo()

	initRsaKeys()
	initExter()
	initDaos()
	initApiHandlers(goapi.ApiRouter)
	initApiFilters(goapi.ApiRouter)
	return nil
}

// available since template-v0.2.0
func initExter() {
	if exterAppId = goapi.AppConfig.GetString("gvabe.exter.app_id"); exterAppId == "" {
		log.Printf("[WARN] No Exter app-id configured at [gvabe.exter.app_id], Exter login is disabled.")
	} else if exterBaseUrl = goapi.AppConfig.GetString("gvabe.exter.base_url"); exterBaseUrl == "" {
		log.Printf("[WARN] No Exter base-url configured at [gvabe.exter.base_url], default value will be used.")
		exterBaseUrl = "https://exteross.gpvcloud.com"
	}
	exterBaseUrl = strings.TrimSuffix(exterBaseUrl, "/") // trim trailing slashes
	if exterAppId != "" {
		exterClient = NewExterClient(exterAppId, exterBaseUrl)
	}
	log.Printf("[INFO] Exter app-id: %s / Base Url: %s", exterAppId, exterBaseUrl)

	go goFetchExterInfo(60)
}

// available since template-v0.2.0
func initRsaKeys() {
	rsaPrivKeyFile := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_file")
	if rsaPrivKeyFile == "" {
		log.Println("[WARN] No RSA private key file configured at [gvabe.keys.rsa_privkey_file], generating one...")
		privKey, err := genRsaKey(2048)
		if err != nil {
			panic(err)
		}
		rsaPrivKey = privKey
	} else {
		log.Println(fmt.Sprintf("[INFO] Loading RSA private key from [%s]...", rsaPrivKeyFile))
		content, err := ioutil.ReadFile(rsaPrivKeyFile)
		if err != nil {
			panic(err)
		}
		block, _ := pem.Decode(content)
		if block == nil {
			panic(fmt.Sprintf("cannot decode PEM from file [%s]", rsaPrivKeyFile))
		}
		var der []byte
		passphrase := goapi.AppConfig.GetString("gvabe.keys.rsa_privkey_passphrase")
		if passphrase != "" {
			log.Println("[INFO] RSA private key is pass-phrase protected")
			if decrypted, err := x509.DecryptPEMBlock(block, []byte(passphrase)); err != nil {
				panic(err)
			} else {
				der = decrypted
			}
		} else {
			der = block.Bytes
		}
		if block.Type == "RSA PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS1PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey
			}
		} else if block.Type == "PRIVATE KEY" {
			if privKey, err := x509.ParsePKCS8PrivateKey(der); err != nil {
				panic(err)
			} else {
				rsaPrivKey = privKey.(*rsa.PrivateKey)
			}
		}
	}

	rsaPubKey = &rsaPrivKey.PublicKey

	if DEBUG {
		if DEBUG {
			log.Printf("[DEBUG] Exter public key: {Size: %d / Exponent: %d / Modulus: %x}",
				rsaPubKey.Size()*8, rsaPubKey.E, rsaPubKey.N)

			pubBlockPKCS1 := pem.Block{
				Type:    "RSA PUBLIC KEY",
				Headers: nil,
				Bytes:   x509.MarshalPKCS1PublicKey(rsaPubKey),
			}
			rsaPubKeyPemPKCS1 := pem.EncodeToMemory(&pubBlockPKCS1)
			log.Printf("[DEBUG] Exter public key (PKCS1): %s", string(rsaPubKeyPemPKCS1))

			pubPKIX, _ := x509.MarshalPKIXPublicKey(rsaPubKey)
			pubBlockPKIX := pem.Block{
				Type:    "PUBLIC KEY",
				Headers: nil,
				Bytes:   pubPKIX,
			}
			rsaPubKeyPemPKIX := pem.EncodeToMemory(&pubBlockPKIX)
			log.Printf("[DEBUG] Exter public key (PKIX): %s", string(rsaPubKeyPemPKIX))
		}
	}
}

func createSqlConnect() *prom.SqlConnect {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		dir := goapi.AppConfig.GetString("gvabe.db.sqlite.directory")
		dbname := goapi.AppConfig.GetString("gvabe.db.sqlite.dbname")
		return bo.NewSqliteConnection(dir, dbname)
	case "pg", "pgsql", "postgres", "postgresql":
		url := goapi.AppConfig.GetString("gvabe.db.pgsql.url")
		return bo.NewPgsqlConnection(url, goapi.AppConfig.GetString("timezone"))
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func createGroupDao(sqlc *prom.SqlConnect) group.GroupDao {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		return group.NewGroupDaoSqlite(sqlc, bo.TableGroup)
	case "pg", "pgsql", "postgres", "postgresql":
		return group.NewGroupDaoPgsql(sqlc, bo.TableGroup)
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

func createUserDao(sqlc *prom.SqlConnect) user.UserDao {
	dbtype := strings.ToLower(goapi.AppConfig.GetString("gvabe.db.type"))
	switch dbtype {
	case "sqlite":
		return user.NewUserDaoSqlite(sqlc, bo.TableUser)
	case "pg", "pgsql", "postgres", "postgresql":
		return user.NewUserDaoPgsql(sqlc, bo.TableUser)
	}
	panic(fmt.Sprintf("unknown databbase type: %s", dbtype))
}

/*----------------------------------------------------------------------*/

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

/*----------------------------------------------------------------------*/

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

		// FIXME this is for demo purpose only!
		if group.Id == systemGroupId {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot edit system group [%s]", group.Id))
		}

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
		return itineris.NewApiResult(itineris.StatusOk)
	}
}

// API handler "deleteGroup"
func apiDeleteGroup(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
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

		// FIXME this is for demo purpose only!
		if group.Id == systemGroupId {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot delete system group [%s]", group.Id))
		}

		if ok, err := groupDao.Delete(group); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		} else if !ok {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("Group [%s] has not been deleted", id.(string)))
		}
		return itineris.NewApiResult(itineris.StatusOk)
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

/*----------------------------------------------------------------------*/

// API handler "userList"
func apiUserList(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	userList, err := userDao.GetAll()
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	}
	data := make([]map[string]interface{}, 0)
	for _, u := range userList {
		data = append(data, map[string]interface{}{
			"username": u.GetUsername(), "name": u.GetName(), "gid": u.GetGroupId(),
		})
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// API handler "getUser"
func apiGetUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	username, _ := params.GetParamAsType("username", reddo.TypeString)
	if username == nil || strings.TrimSpace(username.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	}
	if user, err := userDao.Get(username.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if user == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	} else {
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{
			"username": user.GetUsername(), "name": user.GetName(), "gid": user.GetGroupId(),
		})
	}
}

// API handler "updateUser"
func apiUpdateUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	username, _ := params.GetParamAsType("username", reddo.TypeString)
	if username == nil || strings.TrimSpace(username.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	}
	if user, err := userDao.Get(username.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if user == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	} else {
		// TODO check current user's permission

		// FIXME this is for demo purpose only!
		if user.GetUsername() == systemAdminUsername {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot edit system admin user [%s]", user.GetUsername()))
		}

		password, _ := params.GetParamAsType("password", reddo.TypeString)
		var newPassword, newPassword2 interface{}
		if password != nil && strings.TrimSpace(password.(string)) != "" {
			password = strings.TrimSpace(password.(string))
			if encryptPassword(user.GetUsername(), password.(string)) != user.GetPassword() {
				return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Current password does not match")
			}

			newPassword, _ = params.GetParamAsType("new_password", reddo.TypeString)
			if newPassword == nil || strings.TrimSpace(newPassword.(string)) == "" {
				return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [new_password]")
			}
			newPassword = strings.TrimSpace(newPassword.(string))
			newPassword2, _ = params.GetParamAsType("new_password2", reddo.TypeString)
			if newPassword2 == nil || strings.TrimSpace(newPassword2.(string)) != newPassword {
				return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("New password does not match confirmed one")
			}
		}

		name, _ := params.GetParamAsType("name", reddo.TypeString)
		if name == nil || strings.TrimSpace(name.(string)) == "" {
			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
		}
		name = strings.TrimSpace(name.(string))

		groupId, _ := params.GetParamAsType("group_id", reddo.TypeString)
		if groupId == nil || strings.TrimSpace(groupId.(string)) == "" {
			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [group_id]")
		}
		groupId = strings.TrimSpace(strings.ToLower(groupId.(string)))
		if group, err := groupDao.Get(groupId.(string)); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		} else if group == nil {
			return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Group [%s] does not exist", groupId))
		}

		user.SetName(strings.TrimSpace(name.(string))).
			SetGroupId(groupId.(string))
		if password != nil && strings.TrimSpace(password.(string)) != "" {
			user.SetPassword(encryptPassword(user.GetUsername(), newPassword.(string)))
		}

		if ok, err := userDao.Update(user); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		} else if !ok {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been updated", username.(string)))
		}
		return itineris.NewApiResult(itineris.StatusOk)
	}
}

// API handler "deleteUser"
func apiDeleteUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	username, _ := params.GetParamAsType("username", reddo.TypeString)
	if username == nil || strings.TrimSpace(username.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	}
	if user, err := userDao.Get(username.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if user == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(fmt.Sprintf("User [%s] not found", username))
	} else {
		// TODO check current user's permission

		// FIXME this is for demo purpose only!
		if user.GetUsername() == systemAdminUsername {
			return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(fmt.Sprintf("Cannot delete system admin user [%s]", user.GetUsername()))
		}

		if ok, err := userDao.Delete(user); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
		} else if !ok {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been deleted", username.(string)))
		}
		return itineris.NewApiResult(itineris.StatusOk)
	}
}

// API handler "createUser"
func apiCreateUser(_ *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	username, _ := params.GetParamAsType("username", reddo.TypeString)
	if username == nil || strings.TrimSpace(username.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [username]")
	}
	username = strings.TrimSpace(strings.ToLower(username.(string)))
	if !regexp.MustCompile("^[0-9a-z_]+$").Match([]byte(username.(string))) {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Invalid value for parameter [username]")
	}

	password, _ := params.GetParamAsType("password", reddo.TypeString)
	if password == nil || strings.TrimSpace(password.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [password]")
	}
	password = strings.TrimSpace(password.(string))
	password2, _ := params.GetParamAsType("password2", reddo.TypeString)
	if password2 == nil || strings.TrimSpace(password2.(string)) != password {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Password does not match confirmed one")
	}

	name, _ := params.GetParamAsType("name", reddo.TypeString)
	if name == nil || strings.TrimSpace(name.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [name]")
	}
	name = strings.TrimSpace(name.(string))

	groupId, _ := params.GetParamAsType("group_id", reddo.TypeString)
	if groupId == nil || strings.TrimSpace(groupId.(string)) == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("Missing or empty parameter [group_id]")
	}
	groupId = strings.TrimSpace(strings.ToLower(groupId.(string)))
	if group, err := groupDao.Get(groupId.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if group == nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("Group [%s] does not exist", groupId))
	}

	if user, err := userDao.Get(username.(string)); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if user != nil {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(fmt.Sprintf("User [%s] already existed", username))
	}
	user := user.NewUserBo(username.(string), "").
		SetPassword(encryptPassword(username.(string), password.(string))).
		SetName(name.(string)).
		SetGroupId(groupId.(string)).
		SetAesKey(utils.RandomString(16))
	if ok, err := userDao.Create(user); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(err.Error())
	} else if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(fmt.Sprintf("User [%s] has not been created", username))
	}
	return itineris.NewApiResult(itineris.StatusOk)
}
