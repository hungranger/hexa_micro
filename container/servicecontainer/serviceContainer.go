package servicecontainer

import (
	"hexa_micro/config"

	"github.com/pkg/errors"
)

type ServiceContainer struct {
	FactoryMap map[string]interface{}
	AppConfig  *config.AppConfig
}

func (sc *ServiceContainer) InitApp(filename string) error {
	var err error
	config, err := loadConfig(filename)
	if err != nil {
		return errors.Wrap(err, "")
	}
	sc.AppConfig = config
	return nil
}

func loadConfig(filename string) (*config.AppConfig, error) {
	ac, err := config.ReadConfig(filename)
	if err != nil {
		return nil, errors.Wrap(err, "servicecontainer.loadconfig")
	}
	return ac, nil
}

func (sc *ServiceContainer) BuildUseCase(code string) (interface{}, error) {
	return nil, nil
}

func (sc *ServiceContainer) Get(code string) (interface{}, bool) {
	return nil, false
}

func (sc *ServiceContainer) Put(code string, value interface{}) {
	return
}
