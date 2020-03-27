package main

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container"
	"hexa_micro/pkg/shortenservice/container/servicecontainer"
	rpcClient "hexa_micro/pkg/shortenservice/interface/rpcClient"
	"hexa_micro/pkg/shortenservice/usecase"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/pkg/errors"
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

	rpcHandler := rpcClient.NewRPCHandler(shortenURLUseCase.(usecase.IShortenUseCase))
	err = rpc.Register(rpcHandler)
	if err != nil {
		log.Fatal("error registering API ", err)
	}

	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		log.Fatal("Listener error ", err)
	}

	log.Printf("serving rpc on port %d", 4040)
	http.Serve(listener, nil)
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
