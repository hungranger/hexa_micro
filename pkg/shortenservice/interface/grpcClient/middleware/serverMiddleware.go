package middleware

import (
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
)

func BuildMiddleware(service protobuf.ShortenerServiceServer) protobuf.ShortenerServiceServer {
	to := TimeoutFindMiddleware{service}
	tm := ThrottleMiddleware{&to}
	return &tm
}
