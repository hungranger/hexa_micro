package inmemory

import (
	"fmt"
	"hexa_micro/config"
	"hexa_micro/model"
	"hexa_micro/repository"
	"log"

	"github.com/pkg/errors"
)

type inmemoryRepository struct {
	db map[string]interface{}
}

func NewInmemoryRepository() repository.IRedirectRepository {
	log.Println("repository.NewInmemoryRepo: Connect Inmemory Successfully")
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
	// for _, v := range r.db {
	// 	log.Printf("%s: %s", v.(*model.Redirect).Code, v.(*model.Redirect).URL)
	// }
	return nil
}
