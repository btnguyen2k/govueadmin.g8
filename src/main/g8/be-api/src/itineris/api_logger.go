package itineris

import (
	"encoding/json"
	"fmt"
	"io"
)

/*
IApiLogger is interface used to log API call.
*/
type IApiLogger interface {
	/*
	   PreApiCall is called just before API is actually invoked.
	*/
	PreApiCall(concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams)

	/*
	   PostApiCall is called right after API is invoked.
	*/
	PostApiCall(durationNano, concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams, result *ApiResult)
}

/*----------------------------------------------------------------------*/

/*
BasePerfLogger is base struct to implement API performance logger.
*/
type BasePerfLogger struct {
	FieldAppName        string
	FieldAppVersion     string
	FieldId             string
	FieldStage          string
	FieldApiName        string
	FieldGateway        string
	FieldTimestampStart string
	FieldDuration       string
	FieldConcurrency    string
}

/*
WriterPerfLogger writes API performance logs to an io.Writer in JSON format.
*/
type WriterPerfLogger struct {
	*BasePerfLogger
	writer              io.Writer
	appName, appVersion string
}

/*
NewWriterPerfLogger creates a new WriterPerfLogger instance.
*/
func NewWriterPerfLogger(writer io.Writer, appName, appVersion string) *WriterPerfLogger {
	return &WriterPerfLogger{BasePerfLogger: &BasePerfLogger{
		FieldAppName:        "app_name",
		FieldAppVersion:     "app_version",
		FieldId:             "id",
		FieldStage:          "s",
		FieldApiName:        "api",
		FieldGateway:        "gw",
		FieldTimestampStart: "t",
		FieldDuration:       "d",
		FieldConcurrency:    "c",
	}, writer: writer, appName: appName, appVersion: appVersion}
}

/*
PreApiCall implements IApiLogger.PreApiCall
*/
func (logger *WriterPerfLogger) PreApiCall(concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams) {
	data := map[string]interface{}{
		logger.FieldAppName:        logger.appName,
		logger.FieldAppVersion:     logger.appVersion,
		logger.FieldId:             ctx.GetId(),
		logger.FieldStage:          "START",
		logger.FieldApiName:        ctx.GetApiName(),
		logger.FieldGateway:        ctx.GetGateway(),
		logger.FieldTimestampStart: ctx.GetTimestamp().UnixNano() / 1000000, // convert to milliseconds
		logger.FieldConcurrency:    concurrency,
	}
	logger.writeLog(data)
}

/*
PostApiCall implements IApiLogger.PostApiCall
*/
func (logger *WriterPerfLogger) PostApiCall(durationNano, concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams, result *ApiResult) {
	data := map[string]interface{}{
		logger.FieldAppName:        logger.appName,
		logger.FieldAppVersion:     logger.appVersion,
		logger.FieldId:             ctx.GetId(),
		logger.FieldStage:          "END",
		logger.FieldApiName:        ctx.GetApiName(),
		logger.FieldGateway:        ctx.GetGateway(),
		logger.FieldTimestampStart: ctx.GetTimestamp().UnixNano() / 1000000, // convert to milliseconds
		logger.FieldDuration:       durationNano,
		logger.FieldConcurrency:    concurrency,
	}
	logger.writeLog(data)
}

func (logger *WriterPerfLogger) writeLog(data map[string]interface{}) {
	if logger.writer != nil {
		js, _ := json.Marshal(data)
		fmt.Fprintln(logger.writer, string(js))
	}
}

/*----------------------------------------------------------------------*/

/*
BaseRequestLogger is base struct to implement API request/response logger.
*/
type BaseRequestLogger struct {
	FieldAppName        string
	FieldAppVersion     string
	FieldId             string
	FieldApiName        string
	FieldTimestampStart string
	FieldStage          string
	FieldGateway        string
	FieldDuration       string
	FieldConcurrency    string
	FieldContext        string
	FieldAuth           string
	FieldParams         string
	FieldResult         string
}

/*
WriterRequestLogger writes API request/response logs to an io.Writer in JSON format.
*/
type WriterRequestLogger struct {
	*BaseRequestLogger
	writer              io.Writer
	appName, appVersion string
}

/*
NewWriterRequestLogger creates a new WriterRequestLogger instance.
*/
func NewWriterRequestLogger(writer io.Writer, appName, appVersion string) *WriterRequestLogger {
	return &WriterRequestLogger{BaseRequestLogger: &BaseRequestLogger{
		FieldAppName:        "api_name",
		FieldAppVersion:     "api_version",
		FieldId:             "id",
		FieldApiName:        "api",
		FieldTimestampStart: "t",
		FieldStage:          "s",
		FieldGateway:        "gw",
		FieldDuration:       "d",
		FieldConcurrency:    "c",
		FieldContext:        "context",
		FieldAuth:           "auth",
		FieldParams:         "params",
		FieldResult:         "result",
	}, writer: writer, appName: appName, appVersion: appVersion}
}

func (logger *WriterRequestLogger) buildAuthMap(auth *ApiAuth) map[string]interface{} {
	accessToken := auth.GetAccessToken()
	if len(accessToken) <= 3 {
		accessToken = "***"
	} else {
		accessToken = "***" + accessToken[0:len(accessToken)-3]
	}
	return map[string]interface{}{
		"app_id":      auth.GetAppId(),
		"accessToken": accessToken,
	}
}

/*
PreApiCall implements IApiLogger.PreApiCall
*/
func (logger *WriterRequestLogger) PreApiCall(concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams) {
	data := map[string]interface{}{
		logger.FieldAppName:        logger.appName,
		logger.FieldAppVersion:     logger.appVersion,
		logger.FieldId:             ctx.GetId(),
		logger.FieldStage:          "START",
		logger.FieldApiName:        ctx.GetApiName(),
		logger.FieldGateway:        ctx.GetGateway(),
		logger.FieldTimestampStart: ctx.GetTimestamp().UnixNano() / 1000000, // convert to milliseconds
		logger.FieldConcurrency:    concurrency,
		logger.FieldContext:        ctx.GetAllContextValues(),
		logger.FieldAuth:           logger.buildAuthMap(auth),
		logger.FieldParams:         params.GetAllParams(),
	}
	logger.writeLog(data)
}

/*
PostApiCall implements IApiLogger.PostApiCall
*/
func (logger *WriterRequestLogger) PostApiCall(durationNano, concurrency int64, ctx *ApiContext, auth *ApiAuth, params *ApiParams, result *ApiResult) {
	data := map[string]interface{}{
		logger.FieldAppName:        logger.appName,
		logger.FieldAppVersion:     logger.appVersion,
		logger.FieldId:             ctx.GetId(),
		logger.FieldStage:          "END",
		logger.FieldApiName:        ctx.GetApiName(),
		logger.FieldGateway:        ctx.GetGateway(),
		logger.FieldTimestampStart: ctx.GetTimestamp().UnixNano() / 1000000, // convert to milliseconds
		logger.FieldDuration:       durationNano,
		logger.FieldConcurrency:    concurrency,
		logger.FieldContext:        ctx.GetAllContextValues(),
		logger.FieldAuth:           logger.buildAuthMap(auth),
		logger.FieldParams:         params.GetAllParams(),
		logger.FieldResult:         result.ToMap(),
	}
	logger.writeLog(data)
}

func (logger *WriterRequestLogger) writeLog(data map[string]interface{}) {
	if logger.writer != nil {
		js, _ := json.Marshal(data)
		fmt.Fprintln(logger.writer, string(js))
	}
}
