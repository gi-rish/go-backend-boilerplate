package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/NewStreetTechnologies/go-backend-boilerplate/config"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/controllers"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/db"
	_logger "github.com/NewStreetTechnologies/go-backend-boilerplate/logger"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/repository"
	"github.com/NewStreetTechnologies/go-backend-boilerplate/routes"

	localization_pb "github.com/NewStreetTechnologies/go-grpc-proto/localization-service"

	"github.com/gin-gonic/gin"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	RELEASE = "release"
)

var (
	logger = _logger.Logger
	// cache  = _cache.Cache
)

func main() {
	port := fmt.Sprintf(":%s", config.GetConfig("service.port"))
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(fmt.Sprintf("failed to listen: %v", err))
	}
	multiplexer := cmux.New(l)
	tls := config.GetConfigInBool("service.tls")
	r := gin.New()

	tranRepo := &repository.TranslationRepo{
		DB: db.GetDB(),
	}
	routerConfig := &routes.RouterConfig{
		R:                     r,
		TranslationRepository: tranRepo,
		Logger:                logger,
	}
	router := routerConfig.InitRouter()

	server := controllers.InitTranslationController(tranRepo, logger)

	go serveGRPC(multiplexer, tls, server, port)

	go serveHTTP(multiplexer, tls, router)

	defer l.Close()
	if err := multiplexer.Serve(); !strings.Contains(err.Error(),
		"use of closed network connection") {
		logger.Fatal(err)
	}
}

func serveGRPC(m cmux.CMux, tls bool, server *controllers.TranslationController, port string) {
	l := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	grpcServer := grpc.NewServer()
	if tls {
		tlsCredentials, err := loadTLSCredentials()
		if err != nil {
			panic(err)
		}
		grpcServer = grpc.NewServer(grpc.Creds(tlsCredentials))
	}
	localization_pb.RegisterLocalizationServiceServer(grpcServer, server)
	logger.Infof("Starting gRPC server, running on port %s...", port)
	if err := grpcServer.Serve(l); err != nil {
		panic(fmt.Sprintf("Failed to start gRPC server: %+v !!!!!!", err))
	}
}

func serveHTTP(m cmux.CMux, tls bool, router *gin.Engine) {
	if tls {
		logger.Infof("Starting https server, running on port %s...", config.GetConfig("service.port"))
		if err := http.ServeTLS(m.Match(cmux.TLS()), router, config.GetConfig("service.cert"), config.GetConfig("service.key")); err != nil {
			panic(fmt.Sprintf("Failed to start https server: %+v !!!!!!", err))
		}
	} else {
		logger.Infof("Starting http server, running on port %s...", config.GetConfig("service.port"))
		if err := http.Serve(m.Match(cmux.HTTP1Fast()), router); err != nil {
			panic(fmt.Sprintf("Failed to start http server: %+v !!!!!!", err))
		}
	}
}

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	serverCert, err := tls.LoadX509KeyPair(config.GetConfig("service.cert"), config.GetConfig("service.key"))
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}
	return credentials.NewTLS(config), nil
}
