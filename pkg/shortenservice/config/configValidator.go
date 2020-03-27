package config

import (
	"github.com/pkg/errors"
)

// database code. Need to map to the database code (DataStoreConfig) in the configuration yaml file.
const (
	MONGODB    string = "mongodb"
	REDISDB    string = "redis"
	INMEMORYDB string = "inmemory"
)

// use case code. Need to map to the use case code (UseCaseConfig) in the configuration yaml file.
// Client app use those to retrieve use case from the container
const (
	SHORTEN_URL = "shortenURL"
)

const (
	REDIRECT_REPO = "redirectRepo"
)

func validateConfig(appConfig AppConfig) error {
	err := validateDataStore(appConfig)
	if err != nil {
		return errors.Wrap(err, "validateDatastoreConfig")
	}

	useCase := appConfig.UseCase
	err = validateUseCase(useCase)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func validateDataStore(appConfig AppConfig) error {
	mgc := appConfig.MongoDBConfig
	key := mgc.Code
	mgcMsg := " in validateDataStore doesn't match key = "
	if MONGODB != key {
		errMsg := MONGODB + mgcMsg + key
		return errors.New(errMsg)
	}

	rc := appConfig.RedisConfig
	key = rc.Code
	if REDISDB != key {
		errMsg := REDISDB + mgcMsg + key
		return errors.New(errMsg)
	}

	imc := appConfig.InMemoryConfig
	key = imc.Code
	if INMEMORYDB != key {
		errMsg := INMEMORYDB + mgcMsg + key
		return errors.New(errMsg)
	}

	return nil
}

func validateUseCase(useCase UseCaseConfig) error {
	err := validateShortenURL(useCase)
	if err != nil {
		return errors.Wrap(err, "")
	}
	return nil
}

func validateShortenURL(usecase UseCaseConfig) error {
	sc := usecase.ShortenURL
	key := sc.Code
	scMsg := " in validateShortenURL doesn't match key = "
	if SHORTEN_URL != key {
		errMsg := SHORTEN_URL + scMsg + key
		return errors.New(errMsg)
	}

	key = sc.RedirectRepoConfig.Code
	if REDIRECT_REPO != key {
		errMsg := REDIRECT_REPO + scMsg + key
		return errors.New(errMsg)
	}
	return nil
}
