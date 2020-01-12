/*
Package itineris creates a framework that help building API server over HTTP and gRPC.

    HTTP/gRPC are API communication protocols only. Business logic is handled by one single code repository for all communication protocols.
    Data is interchanged (request/response) in JSON format; can be gzipped to reduce space/transmit time consumption.

@author Thanh Nguyen <btnguyen2k@gmail.com>
@since template-v0.4.r1
*/
package itineris

import (
	"encoding/json"
	"github.com/btnguyen2k/consu/reddo"
	"main/src/utils"
	"reflect"
	"time"
)

/*----------------------------------------------------------------------*/

type IApiHandler func(*ApiContext, *ApiAuth, *ApiParams) *ApiResult

/*----------------------------------------------------------------------*/

const (
	ctxId        = "id"
	ctxTimestamp = "time"
	ctxApiName   = "api_name"
	ctxGateway   = "gateway"
)

/**
ApiContext encapsulates the context information of an API call.
*/
type ApiContext struct {
	contextData map[string]interface{}
}

/*
NewApiContext creates a new ApiContext instance.
*/
func NewApiContext() *ApiContext {
	ctx := &ApiContext{contextData: map[string]interface{}{
		ctxId:        utils.UniqueId(),
		ctxTimestamp: time.Now(),
	}}
	return ctx
}

/*
GetId returns the unique id associated with this API context.
*/
func (ctx *ApiContext) GetId() string {
	v, err := ctx.GetContextValueAsType(ctxId, reddo.TypeString)
	if err != nil {
		panic(err)
	}
	return v.(string)
}

/*
SetId associates a unique id with this API context.
*/
func (ctx *ApiContext) SetId(id string) *ApiContext {
	return ctx.SetContextValue(ctxId, id)
}

/*
GetApiName returns the API name associated with this context.
*/
func (ctx *ApiContext) GetApiName() string {
	v, err := ctx.GetContextValueAsType(ctxApiName, reddo.TypeString)
	if err != nil {
		panic(err)
	}
	return v.(string)
}

/*
SetApiName associates the API name with this context.
*/
func (ctx *ApiContext) SetApiName(apiName string) *ApiContext {
	return ctx.SetContextValue(ctxApiName, apiName)
}

/*
GetGateway returns gateway name associated with this API context.
*/
func (ctx *ApiContext) GetGateway() string {
	v, err := ctx.GetContextValueAsType(ctxGateway, reddo.TypeString)
	if err != nil {
		panic(err)
	}
	return v.(string)
}

/*
SetApiName associates a gateway name this API context.
*/
func (ctx *ApiContext) SetGateway(gateway string) *ApiContext {
	return ctx.SetContextValue(ctxGateway, gateway)
}

/*
GetTimestamp returns the timestamp associated with this API context.
*/
func (ctx *ApiContext) GetTimestamp() time.Time {
	v, err := ctx.GetContextValueAsType(ctxTimestamp, reddo.TypeTime)
	if err != nil {
		panic(err)
	}
	return v.(time.Time)
}

/*
SetTimestamp associates a timestamp with this API context.
*/
func (ctx *ApiContext) SetTimestamp(timestamp time.Time) *ApiContext {
	return ctx.SetContextValue(ctxTimestamp, timestamp)
}

/*
SetContextValue sets a context value.

    If value is nil, the associated context field is removed
*/
func (ctx *ApiContext) SetContextValue(field string, value interface{}) *ApiContext {
	if value == nil {
		return ctx.RemoveContextValue(field)
	}
	ctx.contextData[field] = value
	return ctx
}

/*
RemoveContextValue removes an associated context field value.
*/
func (ctx *ApiContext) RemoveContextValue(field string) *ApiContext {
	delete(ctx.contextData, field)
	return ctx
}

/*
GetContextValue returns a context value.
*/
func (ctx *ApiContext) GetContextValue(field string) interface{} {
	v, ok := ctx.contextData[field]
	if ok {
		return v
	}
	return nil
}

/*
GetContextValueAsType returns a context value casted to a specified type.
*/
func (ctx *ApiContext) GetContextValueAsType(field string, typ reflect.Type) (interface{}, error) {
	return reddo.Convert(ctx.GetContextValue(field), typ)
}

/*
GetAllContextValues returns all context values as a map.
*/
func (ctx *ApiContext) GetAllContextValues() map[string]interface{} {
	return ctx.contextData
}

/*
ToJsonString serializes the ApiContext to JSON string.
*/
func (ctx *ApiContext) ToJsonString() string {
	js, _ := json.Marshal(ctx.contextData)
	return string(js)
}

/*----------------------------------------------------------------------*/

/**
ApiAuth encapsulates authentication information of an API call.
*/
type ApiAuth struct {
	appId       string
	accessToken string
}

/*
NewApiAuth creates a new ApiAuth instance.
*/
func NewApiAuth(appId, accessToken string) *ApiAuth {
	return &ApiAuth{appId: appId, accessToken: accessToken}
}

/*
GetAppId returns the app-id associated with the ApiAuth instance.
*/
func (auth *ApiAuth) GetAppId() string {
	return auth.appId
}

/*
GetAccessToken returns the access-token associated with the ApiAuth instance.
*/
func (auth *ApiAuth) GetAccessToken() string {
	return auth.accessToken
}

/*----------------------------------------------------------------------*/

/**
ApiParams encapsulates parameters to be passed to the API.
*/
type ApiParams struct {
	params map[string]interface{}
}

/*
NewApiParams creates a new ApiParams instance.
*/
func NewApiParams() *ApiParams {
	return &ApiParams{params: map[string]interface{}{}}
}

/*
SetParam sets a parameter value.

    If value is nil, the associated param is removed
*/
func (prm *ApiParams) SetParam(key string, value interface{}) *ApiParams {
	if value == nil {
		return prm.RemoveParam(key)
	}
	prm.params[key] = value
	return prm
}

/*
RemoveParam removes a parameter value.
*/
func (prm *ApiParams) RemoveParam(key string) *ApiParams {
	delete(prm.params, key)
	return prm
}

/*
GetParam returns a parameter value.
*/
func (prm *ApiParams) GetParam(key string) interface{} {
	v, ok := prm.params[key]
	if ok {
		return v
	}
	return nil
}

/*
GetParamAsType returns a parameter value casted to a specified type.
*/
func (prm *ApiParams) GetParamAsType(key string, typ reflect.Type) (interface{}, error) {
	return reddo.Convert(prm.GetParam(key), typ)
}

/*
GetAllParams returns all parameters as a map.
*/
func (prm *ApiParams) GetAllParams() map[string]interface{} {
	return prm.params
}

/*----------------------------------------------------------------------*/

const (
	StatusOk             = 200
	StatusErrorClient    = 400
	StatusNoPermission   = 403
	StatusNotFound       = 404
	StatusDeprecated     = 410
	StatusErrorServer    = 500
	StatusNotImplemented = 501
)

var (
	ResultNoPermission   = NewApiResult(StatusNoPermission).SetMessage("Client do not has permission to call this API")
	ResultNotImplemented = NewApiResult(StatusNotImplemented).SetMessage("API not found")
	ResultNotFound       = NewApiResult(StatusNotFound).SetMessage("Item not found")
)

/*
ApiResult encapsulates result from an API call.
*/
type ApiResult struct {
	Status    int                    `json:"status"`
	Message   string                 `json:"message"`
	Data      interface{}            `json:"data"`
	DebugInfo interface{}            `json:"debug"`
	Extras    map[string]interface{} `json:"extras"`
}

/*
NewApiResult creates a new ApiResult instance.
*/
func NewApiResult(status int) *ApiResult {
	return &ApiResult{Status: status, Extras: make(map[string]interface{})}
}

/*
GetStatus returns result Status.
*/
func (rst *ApiResult) GetStatus() int {
	return rst.Status
}

/*
GetMessage returns result Message.
*/
func (rst *ApiResult) GetMessage() string {
	return rst.Message
}

/*
SetMessage sets the value of result Message.
*/
func (rst *ApiResult) SetMessage(message string) *ApiResult {
	rst.Message = message
	return rst
}

/*
GetData returns result Data.
*/
func (rst *ApiResult) GetData() interface{} {
	return rst.Data
}

/*
SetData sets the value of result Data.
*/
func (rst *ApiResult) SetData(data interface{}) *ApiResult {
	rst.Data = data
	return rst
}

/*
GetExtras returns result extra info.
*/
func (rst *ApiResult) GetExtras() map[string]interface{} {
	return rst.Extras
}

/*
SetExtras sets the value of result extra info.
*/
func (rst *ApiResult) SetExtras(extras map[string]interface{}) *ApiResult {
	rst.Extras = extras
	return rst
}

/*
AddExtraInfo adds an extra info to the result.
*/
func (rst *ApiResult) AddExtraInfo(key string, value interface{}) *ApiResult {
	if rst.Extras == nil {
		rst.Extras = make(map[string]interface{})
	}
	rst.Extras[key] = value
	return rst
}

/*
RemoveExtraInfo removes an extra info from the result.
*/
func (rst *ApiResult) RemoveExtraInfo(key string) *ApiResult {
	if rst.Extras != nil {
		delete(rst.Extras, key)
	}
	return rst
}

/*
GetDebugInfo returns result debug info.
*/
func (rst *ApiResult) GetDebugInfo() interface{} {
	return rst.DebugInfo
}

/*
SetDebugInfo sets the value of result debug Data.
*/
func (rst *ApiResult) SetDebugInfo(debugInfo interface{}) *ApiResult {
	rst.DebugInfo = debugInfo
	return rst
}

/*
Clone replicates the ApiResult instance.
*/
func (rst *ApiResult) Clone() *ApiResult {
	return &ApiResult{
		Status:    rst.Status,
		Message:   rst.Message,
		Data:      rst.Data,
		DebugInfo: rst.DebugInfo,
		Extras:    rst.Extras,
	}
}

/*
ToJsonString serializes the ApiResult to JSON string.
*/
func (rst *ApiResult) ToJsonString() string {
	js, _ := json.Marshal(*rst)
	return string(js)
}

/*
ToMap exports the ApiResult data to a map.
*/
func (rst *ApiResult) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"status": rst.Status,
	}
	if rst.Message != "" {
		m["message"] = rst.Message
	}
	if rst.Data != nil {
		m["data"] = rst.Data
	}
	if rst.DebugInfo != nil {
		m["debug"] = rst.DebugInfo
	}
	if rst.Extras != nil {
		m["extras"] = rst.Extras
	}
	return m
}
