package loggerfactory

import (
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/loggerfactory/logrus"

	"github.com/pkg/errors"
)

// receiver for logrus factory
type LogrusFactory struct{}

// build logrus logger
func (mf *LogrusFactory) Build(lc *config.LogConfig) error {
	err := logrus.RegisterLog(*lc)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}
