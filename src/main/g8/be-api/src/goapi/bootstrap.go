package goapi

/*
IBootstrapper defines an interface for application to hook bootstrapping routines.

Bootstrapper has access to:
- Application configurations via global variable goapi.AppConfig
- itineris.ApiRouter instance via global variable goapi.ApiRouter
*/
type IBootstrapper interface {
	Bootstrap() error
}
