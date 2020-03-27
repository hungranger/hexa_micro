package grpcClient

import (
	"context"
	"hexa_micro/pkg/shortenservice/interface/grpcClient/protobuf"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/usecase"
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

func (h *GRPCHandler) Find(ctx context.Context, request *protobuf.Redirect) (*protobuf.Redirect, error) {
	redirect, err := h.shortenUseCase.Find(request.GetCode())
	if err != nil {
		return nil, err
	}
	return &protobuf.Redirect{
		Code:      redirect.Code,
		Url:       redirect.URL,
		CreatedAt: redirect.CreateAt,
	}, nil
}

func (h *GRPCHandler) Store(ctx context.Context, request *protobuf.Redirect) (*protobuf.Redirect, error) {
	item := &model.Redirect{URL: request.GetUrl()}
	err := h.shortenUseCase.Store(item)
	if err != nil {
		return nil, err
	}

	return &protobuf.Redirect{
		Code:      item.Code,
		Url:       item.URL,
		CreatedAt: item.CreateAt,
	}, nil
}
