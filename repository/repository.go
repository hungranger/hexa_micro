package repository

import "hexa_micro/model"

type IRedirectRepository interface {
	Find(code string) (*model.Redirect, error)
	Store(redirect *model.Redirect) error
}
