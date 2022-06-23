package gvabe

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/btnguyen2k/consu/reddo"
	"github.com/btnguyen2k/goyai"

	"github.com/btnguyen2k/henge"

	"main/src/goapi"
	"main/src/gvabe/bov2/blog"
	"main/src/gvabe/bov2/user"
	"main/src/itineris"
)

// Setup API handlers: application register its api-handlers by calling router.SetHandler(apiName, apiHandlerFunc)
//   - api-handler function must have the following signature:
//     func (itineris.ApiContext, itineris.ApiAuth, itineris.ApiParams) *itineris.ApiResult
func initApiHandlers(router *itineris.ApiRouter) {
	router.SetHandler("info", apiInfo)
	router.SetHandler("login", apiLogin)
	router.SetHandler("verifyLoginToken", apiVerifyLoginToken)
	router.SetHandler("systemInfo", apiSystemInfo)

	router.SetHandler("myFeed", apiMyFeed)
	router.SetHandler("myBlog", apiMyBlog)
	router.SetHandler("createBlogPost", apiCreateBlogPost)
	router.SetHandler("getBlogPost", apiGetBlogPost)
	router.SetHandler("updateBlogPost", apiUpdateBlogPost)
	router.SetHandler("deleteBlogPost", apiDeleteBlogPost)

	router.SetHandler("getUserVoteForPost", apiGetUserVoteForPost)
	router.SetHandler("voteForPost", apiVoteForPost)
}

/*------------------------------ shared variables and functions ------------------------------*/

var (
	// those APIs will not need authentication.
	// "false" means client, however, needs to sends app-id along with the API call
	// "true" means the API is free for public call
	publicApis = map[string]bool{
		"login":            false,
		"info":             true,
		"getApp":           false,
		"verifyLoginToken": true,
		"loginChannelList": true,
	}
)

// available since template-v0.2.0
func _extractParam(params *itineris.ApiParams, paramName string, typ reflect.Type, defValue interface{}, regexp *regexp.Regexp) interface{} {
	v, _ := params.GetParamAsType(paramName, typ)
	if v == nil {
		v = defValue
	}
	if v != nil {
		if _, ok := v.(string); ok {
			v = strings.TrimSpace(v.(string))
			if regexp != nil && !regexp.Match([]byte(v.(string))) {
				return nil
			}
		}
	}
	return v
}

// available since template-v0.2.0
func _currentUserFromContext(ctx *itineris.ApiContext) (*SessionClaims, *user.User, error) {
	sessClaims, ok := ctx.GetContextValue(ctxFieldSession).(*SessionClaims)
	if !ok || sessClaims == nil {
		return nil, nil, nil
	}
	user, err := userDaov2.Get(sessClaims.UserId)
	return sessClaims, user, err
}

/*------------------------------ APIs ------------------------------*/

// API handler "info"
func apiInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	var publicPEM []byte
	if pubDER, err := x509.MarshalPKIXPublicKey(rsaPubKey); err == nil {
		pubBlock := pem.Block{
			Type:    "PUBLIC KEY",
			Headers: nil,
			Bytes:   pubDER,
		}
		publicPEM = pem.EncodeToMemory(&pubBlock)
	} else {
		publicPEM = []byte(err.Error())
	}

	// var m runtime.MemStats
	result := map[string]interface{}{
		"app": map[string]interface{}{
			"name":        goapi.AppConfig.GetString("app.name"),
			"shortname":   goapi.AppConfig.GetString("app.shortname"),
			"version":     goapi.AppConfig.GetString("app.version"),
			"description": goapi.AppConfig.GetString("app.desc"),
		},
		"exter": map[string]interface{}{
			"app_id":   exterAppId,
			"base_url": exterBaseUrl,
		},
		"rsa_public_key": string(publicPEM),
		/*!!! demo purpose only! exposing memory usage is generally not a good idea !!!*/
		// "memory": map[string]interface{}{
		// 	"alloc":     m.Alloc,
		// 	"alloc_str": strconv.FormatFloat(float64(m.Alloc)/1024.0/1024.0, 'f', 1, 64) + " MiB",
		// 	"sys":       m.Sys,
		// 	"sys_str":   strconv.FormatFloat(float64(m.Sys)/1024.0/1024.0, 'f', 1, 64) + " MiB",
		// 	"gc":        m.NumGC,
		// },
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(result)
}

// API handler "systemInfo"
func apiSystemInfo(_ *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	data := lastSystemInfo()
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

func _doLoginExter(ctx *itineris.ApiContext, params *itineris.ApiParams) *itineris.ApiResult {
	token := _extractParam(params, "token", reddo.TypeString, "", nil)
	if token == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage("empty token")
	}
	if DEBUG && exterRsaPubKey != nil {
		exterToken, err := parseExterJwt(token.(string))
		if err != nil {
			log.Printf("[DEBUG] Error parsing submitted JWT: %e", err)
		} else {
			log.Printf("[DEBUG] Submitted JWT: {Id: %s / Type: %s / AppId: %s / UserId: %s / UserName: %s}",
				exterToken.Id, exterToken.Type, exterToken.AppId, exterToken.UserId, exterToken.UserName)
		}
	}

	if exterClient == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_not_enabled",
				&goyai.LocalizeConfig{DefaultMessage: "Exter login is not enabled"}),
		)
	}
	resp, err := exterClient.VerifyLoginToken(token.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_token_validation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if resp.Status != 200 {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"code": resp.Status, "message": resp.Message}}),
		)
	}
	if exterRsaPubKey == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: "Exter login failed, please retry", PluralCount: 0}),
		)
	}
	exterJwt := resp.GetString("data")
	exterToken, err := parseExterJwt(exterJwt)
	if DEBUG {
		if err != nil {
			log.Printf("[DEBUG] Error parsing returned JWT: %e", err)
		} else {
			log.Printf("[DEBUG] Submitted JWT: {Id: %s / Type: %s / AppId: %s / UserId: %s / UserName: %s}",
				exterToken.Id, exterToken.Type, exterToken.AppId, exterToken.UserId, exterToken.UserName)
		}
	}
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_token_validation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if exterToken.Type != "login" {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_exter_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: "Exter login failed, please retry", PluralCount: 0}),
		)
	}
	user, err := createUserFromExterToken(exterToken)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_user_creation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": resp.Message}}),
		)
	}
	if user == nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_user_creation_failed",
				&goyai.LocalizeConfig{DefaultMessage: "User account initialization failed, please retry", PluralCount: 0}),
		)
	}
	claims, err := genLoginClaims(ctx.GetId(), &Session{
		ClientRef:   ctx.GetId(),
		Channel:     loginChannelExter,
		UserId:      user.GetId(),
		DisplayName: user.GetDisplayName(),
		CreatedAt:   time.Now(),
		ExpiredAt:   time.Unix(exterToken.ExpiresAt, 0),
		Data:        []byte(exterJwt),
	})
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_jwt_generation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	jwt, err := genJws(claims)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_jwt_generation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(jwt)
}

func _doLoginForm(ctx *itineris.ApiContext, params *itineris.ApiParams) *itineris.ApiResult {
	username := _extractParam(params, "username", reddo.TypeString, "", nil)
	resultLoginFailed := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_login_failed",
			&goyai.LocalizeConfig{DefaultMessage: "Authentication failed", PluralCount: 0}),
	)
	if username == "" {
		return resultLoginFailed
	}
	password := _extractParam(params, "password", reddo.TypeString, "", nil)
	if password == "" {
		return resultLoginFailed
	}
	user, err := userDaov2.Get(username.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if user == nil {
		return resultLoginFailed
	}
	if encryptPassword(user.GetId(), password.(string)) != user.GetPassword() {
		return resultLoginFailed
	}
	now := time.Now()
	claims, err := genLoginClaims(ctx.GetId(), &Session{
		ClientRef:   ctx.GetId(),
		Channel:     loginChannelForm,
		UserId:      user.GetId(),
		DisplayName: user.GetDisplayName(),
		CreatedAt:   now,
		ExpiredAt:   now.Add(3600 * time.Second),
		Data:        nil,
	})
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_jwt_generation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	jwt, err := genJws(claims)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_jwt_generation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(jwt)
}

// apiLogin handles API call "login".
//   - Upon login successfully, this API returns the login token as JWT.
func apiLogin(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	mode := _extractParam(params, "mode", reddo.TypeString, "form", nil)
	switch strings.ToLower(mode.(string)) {
	case "exter":
		return _doLoginExter(ctx, params)
	default:
		return _doLoginForm(ctx, params)
	}
}

// apiVerifyLoginToken handles API call "verifyLoginToken".
//   - Upon successful, this API returns the login-token.
func apiVerifyLoginToken(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	// firstly extract JWT token from request and convert it into claims
	token := _extractParam(params, "token", reddo.TypeString, "", nil)
	if token == "" {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: "Authentication failed", PluralCount: 0}),
		)
	}
	claims, err := parseLoginToken(token.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if claims.isExpired() {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_login_failed",
				&goyai.LocalizeConfig{DefaultMessage: errorExpiredJwt.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": errorExpiredJwt.Error()}}),
		)
	}

	// lastly return the login-token encoded as JWT
	jwt, err := genJws(claims)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_jwt_generation_failed",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(),
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(jwt)
}

var funcPostToMapTransform = func(m map[string]interface{}) map[string]interface{} {
	u, _ := userDaov2.Get(m[blog.PostFieldOwnerId].(string))
	// transform input map
	result := map[string]interface{}{
		"id":             m[henge.FieldId],
		"t_created":      m[henge.FieldTimeCreated],
		"is_public":      m[blog.PostFieldIsPublic],
		"owner_id":       m[blog.PostFieldOwnerId],
		"title":          m[blog.PostAttrTitle],
		"content":        m[blog.PostAttrContent],
		"num_comments":   m[blog.PostAttrNumComments],
		"num_votes_up":   m[blog.PostAttrNumVotesUp],
		"num_votes_down": m[blog.PostAttrNumVotesDown],
	}
	if t, ok := result["t_created"].(time.Time); ok {
		result["t_created"] = t.In(time.UTC)
	}
	if u != nil {
		result["owner"] = u.ToMap(func(m map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{
				"id":           m[henge.FieldId],
				"mid":          m[user.UserFieldMaskId],
				"is_admin":     m[user.UserAttrIsAdmin],
				"display_name": m[user.UserAttrDisplayName],
			}
		})
	}
	return result
}

// apiMyFeed handles API call "myFeed"
//
// @available since template-v0.2.0
func apiMyFeed(ctx *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if user == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
				&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
		)
	}
	blogPostList, err := blogPostDaov2.GetUserFeedAll(user)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	data := make([]map[string]interface{}, 0)
	for _, p := range blogPostList {
		data = append(data, p.ToMap(funcPostToMapTransform))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// apiMyBlog handles API call "myBlog"
//
// @available since template-v0.2.0
func apiMyBlog(ctx *itineris.ApiContext, _ *itineris.ApiAuth, _ *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if user == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
				&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
		)
	}
	blogPostList, err := blogPostDaov2.GetUserPostsAll(user)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	data := make([]map[string]interface{}, 0)
	for _, p := range blogPostList {
		data = append(data, p.ToMap(funcPostToMapTransform))
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(data)
}

// apiCreateBlogPost handles API call "createBlogPost"
//
// @available since template-v0.2.0
func apiCreateBlogPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if user == nil {
		return itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
				&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
		)
	}
	isPublic := _extractParam(params, "is_public", reddo.TypeBool, false, nil)
	title := _extractParam(params, "title", reddo.TypeString, "", nil)
	if title == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_empty_blog_title",
				&goyai.LocalizeConfig{DefaultMessage: "Blog title is empty, please provide one"}),
		)
	}
	content := _extractParam(params, "content", reddo.TypeString, "", nil)
	if content == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_empty_blog_content",
				&goyai.LocalizeConfig{DefaultMessage: "Blog content is empty, please provide one"}),
		)
	}
	blogPost := blog.NewBlogPost(goapi.AppVersionNumber, user, isPublic.(bool), title.(string), content.(string))
	ok, err := blogPostDaov2.Create(blogPost)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: "General server error occurred", PluralCount: 0}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk)
}

// apiGetBlogPost handles API call "getBlogPost"
//
// @available since template-v0.2.0
func apiGetBlogPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	resultNoPermission := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
			&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
	)
	if user == nil {
		return resultNoPermission
	}
	id := _extractParam(params, "id", reddo.TypeString, "", nil)
	blogPost, err := blogPostDaov2.Get(id.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if blogPost == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_blog_not_exist",
				&goyai.LocalizeConfig{DefaultMessage: "Blog post not found",
					TemplateData: map[string]interface{}{"id": id}}),
		)
	}
	if !blogPost.IsPublic() && blogPost.GetOwnerId() != user.GetId() {
		return resultNoPermission
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(blogPost.ToMap(funcPostToMapTransform))
}

// apiUpdateBlogPost handles API call "updateBlogPost"
//
// @available since template-v0.2.0
func apiUpdateBlogPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	resultNoPermission := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
			&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
	)
	if user == nil {
		return resultNoPermission
	}
	id := _extractParam(params, "id", reddo.TypeString, "", nil)
	blogPost, err := blogPostDaov2.Get(id.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if blogPost == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_blog_not_exist",
				&goyai.LocalizeConfig{DefaultMessage: "Blog post not found",
					TemplateData: map[string]interface{}{"id": id}}),
		)
	}
	if blogPost.GetOwnerId() != user.GetId() {
		return resultNoPermission
	}
	isPublic := _extractParam(params, "is_public", reddo.TypeBool, false, nil)
	title := _extractParam(params, "title", reddo.TypeString, "", nil)
	if title == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_empty_blog_title",
				&goyai.LocalizeConfig{DefaultMessage: "Blog title is empty, please provide one"}),
		)
	}
	content := _extractParam(params, "content", reddo.TypeString, "", nil)
	if content == "" {
		return itineris.NewApiResult(itineris.StatusErrorClient).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_empty_blog_content",
				&goyai.LocalizeConfig{DefaultMessage: "Blog content is empty, please provide one"}),
		)
	}
	blogPost.SetPublic(isPublic.(bool)).SetTitle(title.(string)).SetContent(content.(string))
	ok, err := blogPostDaov2.Update(blogPost)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: "General server error occurred", PluralCount: 0}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk)
}

// apiDeleteBlogPost handles API call "deleteBlogPost"
//
// @available since template-v0.2.0
func apiDeleteBlogPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	resultNoPermission := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
			&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
	)
	if user == nil {
		return resultNoPermission
	}
	id := _extractParam(params, "id", reddo.TypeString, "", nil)
	blogPost, err := blogPostDaov2.Get(id.(string))
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if blogPost == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_blog_not_exist",
				&goyai.LocalizeConfig{DefaultMessage: "Blog post not found",
					TemplateData: map[string]interface{}{"id": id}}),
		)
	}
	if blogPost.GetOwnerId() != user.GetId() {
		return resultNoPermission
	}
	ok, err := blogPostDaov2.Delete(blogPost)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if !ok {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: "General server error occurred", PluralCount: 0}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk)
}

// apiGetUserVoteForPost handles API call "getUserVoteForPost"
//
// @available since template-v0.2.0
func apiGetUserVoteForPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	resultNoPermission := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
			&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
	)
	if user == nil {
		return resultNoPermission
	}
	postId := _extractParam(params, "postId", reddo.TypeString, "", nil).(string)
	vote, err := blogVoteDaov2.GetUserVoteForTarget(user, postId)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	value := 0
	if vote != nil {
		value = vote.GetValue()
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(value)
}

// apiVoteForPost handles API call "voteForPost"
//
// @available since template-v0.2.0
func apiVoteForPost(ctx *itineris.ApiContext, _ *itineris.ApiAuth, params *itineris.ApiParams) *itineris.ApiResult {
	_, user, err := _currentUserFromContext(ctx)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	resultNoPermission := itineris.NewApiResult(itineris.StatusNoPermission).SetMessage(
		i18n.Localize(ctx.GetClientLocale(), "error_no_permission",
			&goyai.LocalizeConfig{DefaultMessage: "Not authorized"}),
	)
	if user == nil {
		return resultNoPermission
	}
	value := _extractParam(params, "vote", reddo.TypeInt, 0, nil).(int64)
	if value == 0 {
		return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{"vote": false})
	}
	postId := _extractParam(params, "postId", reddo.TypeString, "", nil).(string)
	blogPost, err := blogPostDaov2.Get(postId)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	if blogPost == nil {
		return itineris.NewApiResult(itineris.StatusNotFound).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_blog_not_exist",
				&goyai.LocalizeConfig{DefaultMessage: "Blog post not found",
					TemplateData: map[string]interface{}{"id": postId}}),
		)
	}
	if !blogPost.IsPublic() && blogPost.GetOwnerId() != user.GetId() {
		return resultNoPermission
	}
	if value > 1 {
		value = 1
	} else if value < -1 {
		value = -1
	}
	existingVote, err := blogVoteDaov2.GetUserVoteForTarget(user, postId)
	if err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	log.Printf("Existing vote: %#v\n", existingVote)
	newVote := blog.NewBlogVote(goapi.AppVersionNumber, user, blogPost.GetId(), int(value))
	if existingVote == nil {
		// new vote
		if value > 0 {
			blogPost.IncNumVotesUp(1)
		} else {
			blogPost.IncNumVotesDown(1)
		}
		if _, err := blogVoteDaov2.Create(newVote); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
				i18n.Localize(ctx.GetClientLocale(), "error_server",
					&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
						TemplateData: map[string]interface{}{"error": err.Error()}}),
			)
		}
	} else {
		newVote.SetId(existingVote.GetId())
		if existingVote.GetValue() == newVote.GetValue() {
			// cancel existing vote
			newVote.SetValue(0)
			if value > 0 {
				blogPost.IncNumVotesUp(-1)
			} else {
				blogPost.IncNumVotesDown(-1)
			}
		} else {
			// chance existing vote
			if value > 0 {
				blogPost.IncNumVotesUp(1)
				if existingVote.GetValue() != 0 {
					blogPost.IncNumVotesDown(-1)
				}
			} else {
				if existingVote.GetValue() != 0 {
					blogPost.IncNumVotesUp(-1)
				}
				blogPost.IncNumVotesDown(1)
			}
		}
		log.Printf("New vote: %#v\n", newVote)
		if _, err := blogVoteDaov2.Update(newVote); err != nil {
			return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
				i18n.Localize(ctx.GetClientLocale(), "error_server",
					&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
						TemplateData: map[string]interface{}{"error": err.Error()}}),
			)
		}
	}
	if _, err := blogPostDaov2.Update(blogPost); err != nil {
		return itineris.NewApiResult(itineris.StatusErrorServer).SetMessage(
			i18n.Localize(ctx.GetClientLocale(), "error_server",
				&goyai.LocalizeConfig{DefaultMessage: err.Error(), PluralCount: -1,
					TemplateData: map[string]interface{}{"error": err.Error()}}),
		)
	}
	return itineris.NewApiResult(itineris.StatusOk).SetData(map[string]interface{}{
		"vote": true, "value": newVote.GetValue(), "num_votes_up": blogPost.GetNumVotesUp(), "num_votes_down": blogPost.GetNumVotesDown(),
	})
}
