package middleware

import (
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"

	"golang.org/x/net/context"
)

type CallFinder interface {
	CallFind(ctx context.Context, client protobuf.ShortenerServiceClient, msg *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error)
}

func BuildGetMiddleware(cc CallFinder) CallFinder {
	cbcg := CircuitBreakerCallFind{cc}
	tcg := TimeoutCallFindMiddleware{cbcg}
	rcg := RetryCallFindMiddleware{tcg}
	return &rcg
}
