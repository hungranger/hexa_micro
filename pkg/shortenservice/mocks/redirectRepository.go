package mocks

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/repository"

	"github.com/pkg/errors"
)

type redirectReposiotyStub struct {
	db map[string]*model.Redirect
}

func NewRedirectReposiotyStub() repository.IRedirectRepository {
	fmt.Println("Connect RedirectRepositoryStub Successfully")
	return &redirectReposiotyStub{
		map[string]*model.Redirect{
			"hXH2eyrZg": &model.Redirect{
				Code:     "hXH2eyrZg",
				URL:      "https://github.com",
				CreateAt: 1,
			},
			"_XH2ey9WR": &model.Redirect{
				Code:     "_XH2ey9WR",
				URL:      "https://vnexpress.net",
				CreateAt: 2,
			},
		},
	}
}

func (r *redirectReposiotyStub) Find(code string) (*model.Redirect, error) {
	redirect, ok := r.db[code]
	if !ok {
		return nil, errors.Wrap(config.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	return redirect, nil
}

func (r *redirectReposiotyStub) Store(redirect *model.Redirect) error {
	r.db[redirect.Code] = redirect

	return nil
}
