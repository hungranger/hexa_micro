package mocks

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/model"
)

type ShortenUseCaseFake struct {
}

func (s *ShortenUseCaseFake) Find(code string) (*model.Redirect, error) {
	fmt.Printf("ShortenUseCaseFake: Find(%v)\n", code)

	if code == "" {
		return nil, config.ErrRedirectInvalid
	} else if code == "abc" {
		return nil, config.ErrRedirectNotFound
	} else {
		return &model.Redirect{
			Code:     code,
			URL:      "http://vnexpress.net",
			CreateAt: 123456,
		}, nil
	}
}

func (s *ShortenUseCaseFake) Store(redirect *model.Redirect) error {
	fmt.Println("ShortenUseCaseFake: Store(redirect) func")
	fmt.Printf("value passed in: %v\n", redirect)

	if redirect.URL == "" {
		return config.ErrRedirectNotFound
	} else if redirect.URL == "http/vnexpress.net" {
		return config.ErrRedirectInvalid
	} else {
		redirect.Code = "XH2ey9WR"
		redirect.CreateAt = 123456
		return nil
	}
}
