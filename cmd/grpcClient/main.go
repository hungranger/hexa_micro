package main

import (
	"context"
	"fmt"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"log"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure())
	if err != nil {
		log.Panic(err)
	}

	client := protobuf.NewShortenerServiceClient(conn)

	var result []protobuf.Redirect

	git := protobuf.Redirect{Url: "https://github.com"}
	vnexpress := protobuf.Redirect{Url: "https://vnexpress.net"}
	goolge := protobuf.Redirect{Url: "https://google.com"}

	item, err := client.Store(context.Background(), &git)
	if err != nil {
		log.Fatal(err)
	}
	result = append(result, *item)
	item, _ = client.Store(context.Background(), &vnexpress)
	result = append(result, *item)
	item, _ = client.Store(context.Background(), &goolge)
	result = append(result, *item)

	for _, item := range result {
		redirect, _ := client.Find(context.Background(), &item)
		fmt.Printf("%s => %s\n", redirect.GetUrl(), redirect.GetCode())
	}
}
