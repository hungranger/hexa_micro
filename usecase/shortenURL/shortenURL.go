package shortenURL

import (
	"hexa_micro/config"
	"hexa_micro/model"
	"hexa_micro/repository"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

type ShortenURLUseCase struct {
	redirectRepo repository.IRedirectRepository
}

func NewShortenURLUseCase(redirectRepo repository.IRedirectRepository) *ShortenURLUseCase {
	return &ShortenURLUseCase{redirectRepo}
}

func (r *ShortenURLUseCase) Find(code string) (*model.Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *ShortenURLUseCase) Store(redirect *model.Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(config.ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreateAt = time.Now().UTC().Unix()
	return r.redirectRepo.Store(redirect)
}
