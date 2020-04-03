package middleware

import (
	"context"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	SERVICE_THROTTLE = 5 // no more than 5 request at the same time.
)

var tm throttleMutex

type ThrottleMiddleware struct {
	Next protobuf.ShortenerServiceServer
}

type throttleMutex struct {
	mu       sync.RWMutex
	throttle int
}

func (t *throttleMutex) getThrottle() int {
	t.mu.RLock()
	// logger.Log.Infof("Get throttle=%v", t.throttle)
	defer t.mu.RUnlock()
	return t.throttle
}

func (t *throttleMutex) changeThrottle(delta int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.throttle += delta
}

func (tg *ThrottleMiddleware) Find(ctx context.Context, request *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	throttle := tm.getThrottle()
	if throttle >= SERVICE_THROTTLE {
		logger.Log.Infof("Get throttle=%v reached", throttle)
		return nil, status.Error(codes.ResourceExhausted, "Rate limit has reached")
	} else {
		tm.changeThrottle(1)
		resp, err := tg.Next.Find(ctx, request)
		tm.changeThrottle(-1)
		return resp, err
	}
}

func (tg *ThrottleMiddleware) Store(ctx context.Context, request *protobuf.RedirectStoreRequest) (*protobuf.RedirectStoreResponse, error) {
	return tg.Next.Store(ctx, request)
}
