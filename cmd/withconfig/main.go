package main

import (
	"fmt"
	h "hexa_micro/api"
	"hexa_micro/config"
	"hexa_micro/container"
	"hexa_micro/container/servicecontainer"
	"hexa_micro/usecase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pkg/errors"
)

const (
	DEV_CONFIG  = "../../config/appConfigDev.yaml"
	PROD_CONFIG = "../../config/appConfigProd.yaml"
)

func main() {
	configPath := DEV_CONFIG
	container, err := loadConfig(configPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	shortenURLUseCase, err := container.BuildUseCase(config.SHORTEN_URL)
	if err != nil {
		log.Fatal(err)
		return
	}
	handler := h.NewHandler(shortenURLUseCase.(usecase.IShortenUseCase))

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Printf("Listening on port %s\n", httpPort())
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
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
