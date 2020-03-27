package main

import (
	"context"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}

	client := protobuf.NewShortenerServiceClient(conn)

	var result []protobuf.Redirect

	git := protobuf.Redirect{Url: "https://github.com"}
	vnexpress := protobuf.Redirect{Url: "https://vnexpress.net"}
	goolge := protobuf.Redirect{Url: "https://google.com"}

	item, err := client.Store(context.Background(), &git)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}
	result = append(result, *item)
	item, _ = client.Store(context.Background(), &vnexpress)
	result = append(result, *item)
	item, _ = client.Store(context.Background(), &goolge)
	result = append(result, *item)

	for _, item := range result {
		redirect, _ := client.Find(context.Background(), &item)
		logger.Log.Infof("%s => %s\n", redirect.GetUrl(), redirect.GetCode())
	}
}
