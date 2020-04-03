package service

import (
	"context"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
)

type ShortenURLClient struct {
}

func (sc ShortenURLClient) callStore(ctx context.Context, client protobuf.ShortenerServiceClient, url string) *protobuf.RedirectStoreResponse {
	msg := protobuf.RedirectStoreRequest{Url: url}
	item, err := client.Store(ctx, &msg)
	if err != nil {
		logger.Log.Fatalf("%+v", err)
	}
	return item
}

func (sc ShortenURLClient) CallFind(ctx context.Context, client protobuf.ShortenerServiceClient, msg *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	return client.Find(ctx, msg)
}
