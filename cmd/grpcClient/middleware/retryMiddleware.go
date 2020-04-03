package middleware

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	RETRY_COUNT    = 3
	RETRY_INTERVAL = 2000
)

type RetryCallFindMiddleware struct {
	Next CallFinder
}

func (rcf *RetryCallFindMiddleware) CallFind(ctx context.Context, client protobuf.ShortenerServiceClient, msg *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	var err error
	var value *protobuf.RedirectFindResponse
	for i := 0; i < RETRY_COUNT; i++ {
		value, err = rcf.Next.CallFind(ctx, client, msg)
		if err == nil {
			return value, nil
		} else {
			e, _ := status.FromError(err)
			switch e.Code() {
			case codes.DeadlineExceeded:
				fallthrough
			case codes.ResourceExhausted:
				logger.Log.Infof("Retry number %v|error=%+v", i+1, err)
				time.Sleep(time.Duration(RETRY_INTERVAL) * time.Millisecond)
			default:
				return nil, errors.Wrap(err, "")
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("Service's unavailable, retry: %d times, interval: %dms", RETRY_COUNT, RETRY_INTERVAL))
}
