package main

import (
	"context"
	"hexa_micro/cmd/grpcClient/middleware"
	"hexa_micro/cmd/grpcclient/service"
	jaegerTracer "hexa_micro/cmd/grpcclient/tracer"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/container/servicecontainer"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"time"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const (
	DEV_CONFIG         = "../../pkg/shortenservice/config/appConfigDev.yaml"
	PROD_CONFIG        = "../../pkg/shortenservice/config/appConfigProd.yaml"
	HOST_URL    string = "localhost:4040"
)

func main() {
	configPath := DEV_CONFIG
	_, err := loadConfig(configPath)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
		return
	}

	tracer, closer := jaegerTracer.NewTracer()
	defer closer.Close()

	conn, err := grpc.Dial(HOST_URL,
		grpc.WithInsecure(), grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads())),
	)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	client := protobuf.NewShortenerServiceClient(conn)
	serviceClient := middleware.BuildGetMiddleware(service.ShortenURLClient{})

	// item := service.ShortenURLClient{}.CallStore(context.Background(), client, "https://github.com")

	// var result []protobuf.RedirectStoreResponse
	// result = append(result, *item)
	// item = callStore(client, "https://vnexpress.net")
	// result = append(result, *item)
	// item = callStore(client, "https://google.com")
	// result = append(result, *item)
	// for i := 0; i < 10; i++ {
	// 	req := &protobuf.RedirectFindRequest{
	// 		Code: item.GetCode(),
	// 	}
	// 	value, err := serviceClient.CallFind(context.Background(), client, req)
	// 	logger.Log.Infof("%s => http://localhost:8080/%s, err: %v", value.GetUrl(), value.GetCode(), err)
	// }

	testCircuitBreaker(serviceClient, client)

	// for _, item := range result {
	// 	req := &protobuf.RedirectFindRequest{
	// 		Code: item.GetCode(),
	// 	}
	// 	redirect := callFind(client, req)
	// 	logger.Log.Infof("%s => %s", redirect.GetUrl(), redirect.GetCode())
	// }
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

func testCircuitBreaker(serviceClient middleware.CallFinder, client protobuf.ShortenerServiceClient) {
	req := &protobuf.RedirectFindRequest{
		Code: "abc",
	}
	value, err := serviceClient.CallFind(context.Background(), client, req)
	logger.Log.Infof("value: %v, err: %v", value, err)
	time.Sleep(time.Duration(20*1000) * time.Millisecond)

	value, err = serviceClient.CallFind(context.Background(), client, req)
	logger.Log.Infof("value: %v, err: %v", value, err)
	time.Sleep(time.Duration(5*1000) * time.Millisecond)

	value, err = serviceClient.CallFind(context.Background(), client, req)
	logger.Log.Infof("value: %v, err: %v", value, err)
	time.Sleep(time.Duration(20*2000) * time.Millisecond)

	value, err = serviceClient.CallFind(context.Background(), client, req)
	logger.Log.Infof("value: %v, err: %v", value, err)
}
