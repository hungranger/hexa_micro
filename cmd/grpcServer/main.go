package main

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/container/servicecontainer"
	grpcClient "hexa_micro/pkg/shortenservice/interface/grpcClient"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/middleware"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"hexa_micro/pkg/shortenservice/usecase"
	"io"
	"net"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	DEV_CONFIG                      = "../../pkg/shortenservice/config/appConfigDev.yaml"
	PROD_CONFIG                     = "../../pkg/shortenservice/config/appConfigProd.yaml"
	HOST_URL                 string = "localhost:4040"
	SERVICE_NAME_SHORTEN_URL string = "ShortenURL Service"
)

func newTracer(service string) (opentracing.Tracer, io.Closer) {
	cfg := &jaegerCfg.Configuration{
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			LogSpans: true,
		},
		Disabled: true,
	}
	tracer, closer, err := cfg.New(service, jaegerCfg.Logger(jaeger.StdLogger))
	if err != nil {
		logger.Log.Fatalf("Cannot init Jaeger: %+v", err)
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer
}

func main() {
	configPath := DEV_CONFIG
	container, err := loadConfig(configPath)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
		return
	}

	shortenURLUseCase, _ := container.BuildUseCase(config.SHORTEN_URL)
	grpcHandler := grpcClient.NewGRPCHandler(shortenURLUseCase.(usecase.IShortenUseCase))
	throttleHandler := middleware.BuildMiddleware(grpcHandler)

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	tracer, closer := newTracer(SERVICE_NAME_SHORTEN_URL)
	defer closer.Close()

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer, otgrpc.LogPayloads()),
		),
	}

	srv := grpc.NewServer(opts...)
	protobuf.RegisterShortenerServiceServer(srv, throttleHandler)
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
