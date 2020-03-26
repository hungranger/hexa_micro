package servicecontainer

import (
	"hexa_micro/config"
	"hexa_micro/repository"
	"hexa_micro/repository/RedirectRepository/inmemory"
	mongo "hexa_micro/repository/RedirectRepository/mongodb"
	"hexa_micro/repository/RedirectRepository/redis"
	"hexa_micro/usecase/shortenURL"

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
	switch code {
	case config.SHORTEN_URL:
		shortenURLCfg := sc.AppConfig.UseCase.ShortenURL
		redirectRepo, err := sc.buildRepo(&shortenURLCfg.RedirectRepoConfig)
		if err != nil {
			return nil, err
		}
		return shortenURL.NewShortenURLUseCase(redirectRepo), nil
	}
	return nil, nil
}

func (sc *ServiceContainer) Get(code string) (interface{}, bool) {
	value, found := sc.FactoryMap[code]
	return value, found
}

func (sc *ServiceContainer) Put(code string, value interface{}) {
	sc.FactoryMap[code] = value
}

func (sc *ServiceContainer) buildRepo(repoCfg *config.RepoConfig) (repository.IRedirectRepository, error) {
	switch repoCfg.Code {
	case config.REDIRECT_REPO:
		var repo repository.IRedirectRepository
		var err error

		switch repoCfg.DataStoreConfig.Code {
		case config.INMEMORYDB:
			repo = inmemory.NewInmemoryRepository()
			return repo, nil
		case config.MONGODB:
			repo, err = mongo.NewMongoRepository(repoCfg.DataStoreConfig.UrlAddress, repoCfg.DataStoreConfig.DbName, repoCfg.DataStoreConfig.Timeout)
		case config.REDISDB:
			repo, err = redis.NewRedisRepository(repoCfg.DataStoreConfig.UrlAddress, repoCfg.DataStoreConfig.Password)
		}

		if err != nil {
			return nil, errors.Wrap(err, "")
		}
		return repo, nil
	}
	return nil, nil
}
