package middleware

import (
	"context"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	FIND_TIMEOUT = 300
)

type TimeoutFindMiddleware struct {
	Next protobuf.ShortenerServiceServer
}

func (to *TimeoutFindMiddleware) Find(ctx context.Context, request *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	var cancelFunc context.CancelFunc
	var ch = make(chan bool)
	var err error
	var value *protobuf.RedirectFindResponse
	ctx, cancelFunc = context.WithTimeout(ctx, FIND_TIMEOUT*time.Millisecond)

	go func() {
		value, err = to.Next.Find(ctx, request)
		ch <- true
	}()

	select {
	case <-ctx.Done():
		cancelFunc()
		err = status.Error(codes.DeadlineExceeded, "Server Timeout")
	case <-ch:
		// logger.Log.Info("CallFind finished normally")
	}

	return value, err
}

func (to *TimeoutFindMiddleware) Store(ctx context.Context, request *protobuf.RedirectStoreRequest) (*protobuf.RedirectStoreResponse, error) {
	return to.Next.Store(ctx, request)
}
