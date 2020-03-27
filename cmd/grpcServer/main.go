package main

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/container/servicecontainer"
	grpcClient "hexa_micro/pkg/shortenservice/interface/grpcClient"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"hexa_micro/pkg/shortenservice/usecase"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DEV_CONFIG  = "../../pkg/shortenservice/config/appConfigDev.yaml"
	PROD_CONFIG = "../../pkg/shortenservice/config/appConfigProd.yaml"
)

func main() {
	configPath := DEV_CONFIG
	container, err := loadConfig(configPath)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
		return
	}

	shortenURLUseCase, _ := container.BuildUseCase(config.SHORTEN_URL)
	grpcHandler := grpcClient.NewGRPCHandler(shortenURLUseCase.(usecase.IShortenUseCase))

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	srv := grpc.NewServer()
	protobuf.RegisterShortenerServiceServer(srv, grpcHandler)
	reflection.Register(srv)

	logger.Log.Infof("serving grpc on port %d\n", 4040)
	if err := srv.Serve(listener); err != nil {
		logger.Log.Fatalf("%+v", err)
	}
}

func loadConfig(filePath string) (container.Container, error) {
	factoryMap := make(map[string]interface{})
	appConfig := config.AppConfig{}
	container := servicecontainer.ServiceContainer{factoryMap, &appConfig}
	err := container.InitApp(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "")
	}
	return &container, nil
}
