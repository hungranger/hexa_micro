package middleware

import (
	"context"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	FIND_TIMEOUT = 200
)

type TimeoutCallFindMiddleware struct {
	Next CallFinder
}

func (to TimeoutCallFindMiddleware) CallFind(ctx context.Context, client protobuf.ShortenerServiceClient, msg *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	var cancelFunc context.CancelFunc
	var ch = make(chan bool)
	var err error
	var value *protobuf.RedirectFindResponse
	ctx, cancelFunc = context.WithTimeout(ctx, FIND_TIMEOUT*time.Millisecond)
	go func() {
		value, err = to.Next.CallFind(ctx, client, msg)
		ch <- true
	}()

	select {
	case <-ctx.Done():
		cancelFunc()
		err = status.Error(codes.DeadlineExceeded, "Client Timeout")
	case <-ch:
		// logger.Log.Info("CallFind finished normally")
	}

	return value, err
}
