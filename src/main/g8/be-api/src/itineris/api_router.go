package itineris

import "sync/atomic"

/*
ApiRouter is responsible to routing API call to handler.
*/
type ApiRouter struct {
	concurrency int64
	apiFilter   IApiFilter
	handlersMap map[string]IApiHandler
}

/*
NewApiRouter creates a new ApiRouter instance.
*/
func NewApiRouter() *ApiRouter {
	return &ApiRouter{concurrency: 0, handlersMap: map[string]IApiHandler{}}
}

/*
GetConcurrency returns current number of concurrent API calls.
*/
func (router *ApiRouter) GetConcurrency() int64 {
	return router.concurrency
}

/*
GetApiFilter returns the associated api-filter.
*/
func (router *ApiRouter) GetApiFilter() IApiFilter {
	return router.apiFilter
}

/*
SetApiFilter associated an api-filter to the router.
*/
func (router *ApiRouter) SetApiFilter(apiFilter IApiFilter) *ApiRouter {
	router.apiFilter = apiFilter
	return router
}

/*
GetHandler returns an api-handler by name.
*/
func (router *ApiRouter) GetHandler(apiName string) IApiHandler {
	f, ok := router.handlersMap[apiName]
	if ok {
		return f
	}
	return nil
}

/*
SetHandler maps an handler to api name.
*/
func (router *ApiRouter) SetHandler(apiName string, handler IApiHandler) *ApiRouter {
	if handler == nil {
		return router.RemoveHandler(apiName)
	}
	router.handlersMap[apiName] = handler
	return router
}

/*
RemoveHandler removes an api-handler.
*/
func (router *ApiRouter) RemoveHandler(apiName string) *ApiRouter {
	delete(router.handlersMap, apiName)
	return router
}

/*
GetAllHandlers returns all api-handlers as a map.
*/
func (router *ApiRouter) GetAllHandlers() map[string]IApiHandler {
	return router.handlersMap
}

/*
CallApi performs an API call.
*/
func (router *ApiRouter) CallApi(ctx *ApiContext, auth *ApiAuth, params *ApiParams) *ApiResult {
	atomic.AddInt64(&router.concurrency, 1)
	defer atomic.AddInt64(&router.concurrency, -1)
	apiName := ctx.GetApiName()
	handler := router.GetHandler(apiName)
	var apiResult *ApiResult
	if handler == nil {
		apiResult = NewApiResult(StatusNotImplemented).SetMessage("No handler for API [" + apiName + "].")
	} else if router.apiFilter != nil {
		apiResult = router.apiFilter.Call(handler, ctx, auth, params)
	} else {
		apiResult = handler(ctx, auth, params)
	}
	return apiResult
}
