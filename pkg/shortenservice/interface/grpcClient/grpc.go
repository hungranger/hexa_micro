package grpcClient

import (
	"context"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/usecase"
	"time"
)

// type RPCServer interface {
// 	Find(string, *shortener.Redirect) error
// 	Store(shortener.Redirect, *shortener.Redirect) error
// }

type GRPCHandler struct {
	shortenUseCase usecase.IShortenUseCase
}

func NewGRPCHandler(shortenUseCase usecase.IShortenUseCase) protobuf.ShortenerServiceServer {
	return &GRPCHandler{shortenUseCase}
}

func (h *GRPCHandler) Find(ctx context.Context, request *protobuf.RedirectFindRequest) (*protobuf.RedirectFindResponse, error) {
	// test client|server timeout middleware
	time.Sleep(2 * time.Second)
	redirect, err := h.shortenUseCase.Find(request.GetCode())
	if err != nil {
		return nil, err
	}
	return &protobuf.RedirectFindResponse{
		Code: redirect.Code,
		Url:  redirect.URL,
	}, nil
}

func (h *GRPCHandler) Store(ctx context.Context, request *protobuf.RedirectStoreRequest) (*protobuf.RedirectStoreResponse, error) {
	item := &model.Redirect{URL: request.GetUrl()}
	err := h.shortenUseCase.Store(item)
	if err != nil {
		return nil, err
	}

	return &protobuf.RedirectStoreResponse{
		Code: item.Code,
		Url:  item.URL,
	}, nil
}
