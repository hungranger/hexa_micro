package api

import (
	shortener "hexa_micro/shotener"
	"log"

	"github.com/pkg/errors"
)

type RPCHandler struct {
	redirectService shortener.RedirectService
}

func NewRPCHandler(redirectService shortener.RedirectService) *RPCHandler {
	return &RPCHandler{redirectService}
}

func (h *RPCHandler) Find(code string, reply *shortener.Redirect) error {
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Find")
	}

	*reply = *redirect
	return nil
}

func (h *RPCHandler) Store(item *shortener.Redirect, reply *shortener.Redirect) error {
	err := h.redirectService.Store(item)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Store")
	}

	*reply = *item
	return nil
}
