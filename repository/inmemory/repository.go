package inmemory

import (
	"fmt"
	shortener "hexa_micro/shotener"
	"log"

	"github.com/pkg/errors"
)

type inmemoryRepository struct {
	db map[string]interface{}
}

func NewInmemoryRepository() shortener.RedirectRepository {
	log.Println("repository.NewInmemoryRepo: Connect Inmemory Successfully")
	return &inmemoryRepository{
		map[string]interface{}{},
	}
}

func (r *inmemoryRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect: %s", code)
}

func (r *inmemoryRepository) Find(code string) (*shortener.Redirect, error) {
	key := r.generateKey(code)
	redirect, ok := r.db[key]
	if !ok {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	return redirect.(*shortener.Redirect), nil
}

func (r *inmemoryRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	r.db[key] = redirect
	// for _, v := range r.db {
	// 	log.Printf("%s: %s", v.(*shortener.Redirect).Code, v.(*shortener.Redirect).URL)
	// }
	return nil
}
