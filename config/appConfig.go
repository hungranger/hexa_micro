package config

import (
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	MongoDBConfig  DataStoreConfig `yaml:"mongoDBConfig"`
	RedisConfig    DataStoreConfig `yaml:"redisConfig"`
	InMemoryConfig DataStoreConfig `yaml:"inMemoryConfig"`
}

type DataStoreConfig struct {
	Code string `yaml:"code"`
	// Only database has a driver name, for grpc it is "tcp" ( network) for server
	DriverName string `yaml:"driverName"`
	// For database, this is datasource name; for grpc, it is target url
	UrlAddress string `yaml:"urlAddress"`
	// Only some databases need this password
	Password string `yaml:"password"`
	// Only some databases need this database name
	DbName string `yaml:"dbName"`
	// Only some databases need this timeout
	Timeout int `yaml:"timeout"`
}

// ReadConfig reads the file of the filename (in the same folder) and put it into the AppConfig
func ReadConfig(filename string) (*AppConfig, error) {
	var ac AppConfig
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read config file error")
	}

	err = yaml.Unmarshal(file, &ac)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal")
	}

	err = validateConfig(ac)
	if err != nil {
		return nil, errors.Wrap(err, "validate config")
	}

	log.Print("appConfig:", ac)
	return &ac, nil
}
