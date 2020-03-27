package loggerfactory

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/loggerfactory/zap"

	"github.com/pkg/errors"
)

// receiver for zap factory
type ZapFactory struct{}

// build zap logger
func (mf *ZapFactory) Build(lc *config.LogConfig) error {
	err := zap.RegisterLog(*lc)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
