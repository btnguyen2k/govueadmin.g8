package goapi

import (
	"encoding/json"
	"fmt"
	hocon "github.com/go-akka/configuration"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	pb "main/grpc"
	"main/src/itineris"
	"main/src/utils"
	"net"
	"os"
	"regexp"
	"time"
)

const (
	defaultConfigFile = "./config/application.conf"
)

var (
	AppConfig *hocon.Config
	ApiRouter *itineris.ApiRouter
)

/*
Start bootstraps the application.
*/
func Start(bootstrappers ...IBootstrapper) {
	var err error

	// load application configurations
	AppConfig = initAppConfig()
	httpHeaderAppId = AppConfig.GetString("api.http.header_app_id")
	httpHeaderAccessToken = AppConfig.GetString("api.http.header_access_token")

	// setup api-router
	ApiRouter = itineris.NewApiRouter()

	// initialize "Location"
	utils.Location, err = time.LoadLocation(AppConfig.GetString("timezone"))
	if err != nil {
		panic(err)
	}

	// bootstrapping
	if bootstrappers != nil {
		for _, b := range bootstrappers {
			log.Println("Bootstrapping", b)
			err := b.Bootstrap()
			if err != nil {
				log.Println(err)
			}
		}
	}

	// initialize and start gRPC server
	initGrpcServer()

	// initialize and start echo server
	initEchoServer()
}

func initAppConfig() *hocon.Config {
	configFile := os.Getenv("APP_CONFIG")
	if configFile == "" {
		log.Printf("No environment APP_CONFIG found, fallback to [%s]", defaultConfigFile)
		configFile = defaultConfigFile
	}
	return loadAppConfig(configFile)
}

/*
@since template-v0.4.r2
*/
func initGrpcServer() {
	listenPort := AppConfig.GetInt32("api.grpc.listen_port", 0)
	if listenPort <= 0 {
		log.Println("No valid [api.grpc.listen_port] configured, gRPC API gateway is disabled.")
		return
	}
	listenAddr := AppConfig.GetString("api.grpc.listen_addr", "127.0.0.1")
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", listenAddr, listenPort))
	if err != nil {
		log.Printf("Failed to listen gRPC: %v", err)
		return
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPApiServiceServer(grpcServer, newGrpcGateway())
	log.Printf("Starting [%s] gRPC server on [%s:%d]...\n", AppConfig.GetString("app.name")+" v"+AppConfig.GetString("app.version"), listenAddr, listenPort)
	go grpcServer.Serve(lis)
}

func initEchoServer() {
	listenPort := AppConfig.GetInt32("api.http.listen_port", 0)
	if listenPort <= 0 {
		log.Println("No valid [api.http.listen_port] configured, REST API gateway is disabled.")
		return
	}
	listenAddr := AppConfig.GetString("api.http.listen_addr", "127.0.0.1")
	e := echo.New()
	requestTimeout := AppConfig.GetTimeDuration("api.request_timeout", time.Duration(0))
	if requestTimeout > 0 {
		e.Server.ReadTimeout = requestTimeout
	}
	bodyLimit := AppConfig.GetByteSize("api.max_request_size")
	if bodyLimit != nil && bodyLimit.Int64() > 0 {
		e.Use(middleware.BodyLimit(bodyLimit.String()))
	}
	allowOgirinsStr := AppConfig.GetString("api.http.allow_origins", "*")
	if allowOgirins := regexp.MustCompile("[,;\\s]+").Split(allowOgirinsStr, -1); len(allowOgirins) > 0 {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: allowOgirins,
		}))
	}

	// register API http endpoints
	hasEndpoints := false
	confV := AppConfig.GetValue("api.http.endpoints")
	if confV != nil && confV.IsObject() {
		for uri, uriO := range confV.GetObject().Items() {
			if uriO.IsObject() && !uriO.IsEmpty() {
				hasEndpoints = true
				e.Any(uri, apiHttpHandler)
				for httpMethod, apiName := range uriO.GetObject().Items() {
					registerHttpHandler(uri, httpMethod, apiName.GetString())
				}
			}
		}
	}
	js, _ := json.Marshal(httpRoutingMap)
	log.Println("API http endpoints: " + string(js))
	if !hasEndpoints {
		log.Println("No valid HTTP endpoints defined at key [api.http.endpoints].")
	}
	log.Printf("Starting [%s] RESTful server on [%s:%d]...\n", AppConfig.GetString("app.name")+" v"+AppConfig.GetString("app.version"), listenAddr, listenPort)
	go e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", listenAddr, listenPort)))
}
