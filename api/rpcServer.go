package api

import (
	shortener "hexa_micro/shotener"
	"log"

	"github.com/pkg/errors"
)

// type RPCServer interface {
// 	Find(string, *shortener.Redirect) error
// 	Store(shortener.Redirect, *shortener.Redirect) error
// }

type SimpleRPC struct {
	redirectService shortener.RedirectService
}

func NewSimpleRPC(redirectService shortener.RedirectService) *SimpleRPC {
	return &SimpleRPC{redirectService}
}

func (s *SimpleRPC) Find(code string, reply *shortener.Redirect) error {
	redirect, err := s.redirectService.Find(code)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Find")
	}

	*reply = *redirect
	return nil
}

func (s *SimpleRPC) Store(item *shortener.Redirect, reply *shortener.Redirect) error {
	err := s.redirectService.Store(item)
	if err != nil {
		log.Println(err)
		return errors.Wrap(err, "api.SimpleRPC.Store")
	}

	*reply = *item
	return nil
}
