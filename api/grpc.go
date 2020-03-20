package api

import (
	"context"
	"hexa_micro/serializer/protobuf"
	shortener "hexa_micro/shotener"
)

// type RPCServer interface {
// 	Find(string, *shortener.Redirect) error
// 	Store(shortener.Redirect, *shortener.Redirect) error
// }

type GRPCHandler struct {
	redirectService shortener.RedirectService
}

func NewGRPCHandler(redirectService shortener.RedirectService) protobuf.ShortenerServiceServer {
	return &GRPCHandler{redirectService}
}

func (h *GRPCHandler) Find(ctx context.Context, request *protobuf.Redirect) (*protobuf.Redirect, error) {
	redirect, err := h.redirectService.Find(request.GetCode())
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
	item := &shortener.Redirect{URL: request.GetUrl()}
	err := h.redirectService.Store(item)
	if err != nil {
		return nil, err
	}

	return &protobuf.Redirect{
		Code:      item.Code,
		Url:       item.URL,
		CreatedAt: item.CreateAt,
	}, nil
}
