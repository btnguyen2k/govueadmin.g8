package itineris

/*
IApiAuthenticator is interface used to authenticate API call.
*/
type IApiAuthenticator interface {
	Authenticate(*ApiContext, *ApiAuth) bool
}
