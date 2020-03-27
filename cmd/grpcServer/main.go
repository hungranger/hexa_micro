package main

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container"
	"hexa_micro/pkg/shortenservice/container/servicecontainer"
	grpcClient "hexa_micro/pkg/shortenservice/interface/grpcClient"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"hexa_micro/pkg/shortenservice/usecase"
	"log"
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
		log.Fatal(err)
		return
	}

	shortenURLUseCase, _ := container.BuildUseCase(config.SHORTEN_URL)
	grpcHandler := grpcClient.NewGRPCHandler(shortenURLUseCase.(usecase.IShortenUseCase))

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Panicln(err)
	}

	srv := grpc.NewServer()
	protobuf.RegisterShortenerServiceServer(srv, grpcHandler)
	reflection.Register(srv)

	log.Printf("serving grpc on port %d\n", 4040)
	if err := srv.Serve(listener); err != nil {
		log.Panic(err)
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
