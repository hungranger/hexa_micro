package config

import (
	"log"

	"github.com/pkg/errors"
)

// database code. Need to map to the database code (DataStoreConfig) in the configuration yaml file.
const (
	MONGODB    string = "mongodb"
	REDISDB    string = "redis"
	INMEMORYDB string = "inmemory"
)

func validateConfig(appConfig AppConfig) error {
	err := validateDataStore(appConfig)
	if err != nil {
		return errors.Wrap(err, "validateDatastoreConfig")
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
	log.Print(imc)
	if INMEMORYDB != key {
		errMsg := INMEMORYDB + mgcMsg + key
		return errors.New(errMsg)
	}

	return nil
}
