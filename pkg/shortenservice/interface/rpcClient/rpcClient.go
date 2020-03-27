package rpcClient

import (
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/usecase"
	"log"

	"github.com/pkg/errors"
)

type RPCHandler struct {
	shortenURLUseCase usecase.IShortenUseCase
}

func NewRPCHandler(shortenURLUseCase usecase.IShortenUseCase) *RPCHandler {
	return &RPCHandler{shortenURLUseCase}
}

func (h *RPCHandler) Find(code string, reply *model.Redirect) error {
	redirect, err := h.shortenURLUseCase.Find(code)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Find")
	}

	*reply = *redirect
	return nil
}

func (h *RPCHandler) Store(item *model.Redirect, reply *model.Redirect) error {
	err := h.shortenURLUseCase.Store(item)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Store")
	}

	*reply = *item
	return nil
}
