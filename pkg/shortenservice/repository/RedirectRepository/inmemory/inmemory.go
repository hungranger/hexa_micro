package inmemory

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/repository"

	"github.com/pkg/errors"
)

type inmemoryRepository struct {
	db map[string]interface{}
}

func NewInmemoryRepository() repository.IRedirectRepository {
	logger.Log.Info("Connect Inmemory Successfully")
	return &inmemoryRepository{
		map[string]interface{}{},
	}
}

func (r *inmemoryRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect: %s", code)
}

func (r *inmemoryRepository) Find(code string) (*model.Redirect, error) {
	key := r.generateKey(code)
	redirect, ok := r.db[key]
	if !ok {
		return nil, errors.Wrap(config.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	return redirect.(*model.Redirect), nil
}

func (r *inmemoryRepository) Store(redirect *model.Redirect) error {
	key := r.generateKey(redirect.Code)
	r.db[key] = redirect

	return nil
}

func (r *inmemoryRepository) FindAll() ([]model.Redirect, error) {
	var result []model.Redirect
	for _, v := range r.db {
		result = append(result, model.Redirect{
			Code:     v.(model.Redirect).Code,
			URL:      v.(*model.Redirect).URL,
			CreateAt: v.(*model.Redirect).CreateAt,
		})
	}

	return result, nil
}
