package repository

import "hexa_micro/pkg/shortenservice/model"

type IRedirectRepository interface {
	Find(code string) (*model.Redirect, error)
	Store(redirect *model.Redirect) error
}
