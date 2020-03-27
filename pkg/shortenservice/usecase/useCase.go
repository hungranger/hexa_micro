// It is the entry point for the application's business logic. It is a top level package for a Micro-service application.
// This top level package only defines interface, the concrete implementations are defined in sub-package of it.
// It only depends on model package. No other package should dependent on it except cmd.
package usecase

import "hexa_micro/pkg/shortenservice/model"

type IShortenUseCase interface {
	Find(code string) (*model.Redirect, error)
	Store(redirect *model.Redirect) error
}
