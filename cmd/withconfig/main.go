package main

import (
	"hexa_micro/config"
	"hexa_micro/container"
	"hexa_micro/container/servicecontainer"
	"log"

	"github.com/pkg/errors"
)

const (
	DEV_CONFIG  = "../../config/appConfigDev.yaml"
	PROD_CONFIG = "../../config/appConfigProd.yaml"
)

func main() {
	configPath := DEV_CONFIG
	_, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
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
