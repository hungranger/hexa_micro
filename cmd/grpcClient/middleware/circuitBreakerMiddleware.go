package middleware

import (
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"time"

	"github.com/pkg/errors"
	"github.com/sony/gobreaker"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var cb *gobreaker.CircuitBreaker

type CircuitBreakerCallFind struct {
	Next CallFinder
}

func init() {
	st := gobreaker.Settings{
		Name:          "CircuitBreakerCallFind",
		MaxRequests:   2,
		Timeout:       5 * time.Second,
		ReadyToTrip:   readyToTrip,
		OnStateChange: onStateChange,
	}
	cb = gobreaker.NewCircuitBreaker(st)
}

func readyToTrip(counts gobreaker.Counts) bool {
	// logger.Log.Infof("req: %d, fail: %v, success: %v", counts.Requests, counts.TotalFailures, counts.TotalSuccesses)
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.Requests >= 2 && failureRatio >= 0.6
}

func onStateChange(name string, from gobreaker.State, to gobreaker.State) {
	logger.Log.Infof("CB: %v %v=>%v", name, from, to)
}

func (cbm CircuitBreakerCallFind) CallFind(ctx context.Context, client protobuf.ShortenerServiceClient, msg *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	var err error
	var value *protobuf.RedirectFindResponse
	var serivceUp bool
	_, cbErr := cb.Execute(func() (interface{}, error) {
		value, err = cbm.Next.CallFind(ctx, client, msg)
		// need to decide which error represent to unavailable service
		if err != nil {
			e, _ := status.FromError(err)
			switch e.Code() {
			case codes.DeadlineExceeded:
				fallthrough
			case codes.ResourceExhausted:
				return nil, errors.Wrap(err, "Circuit breaker:")
			}
		}
		serivceUp = true
		return value, nil
	})

	if !serivceUp {
		cbErr = status.Error(codes.Unavailable, cbErr.Error())
		return nil, cbErr
	}
	logger.Log.Infof("value: %v, err: %v", value, err)

	return value, err
}
