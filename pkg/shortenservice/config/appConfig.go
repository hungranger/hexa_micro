package config

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var (
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	ErrRedirectInvalid  = errors.New("Redirect Invalid")
)

type AppConfig struct {
	MongoDBConfig  DataStoreConfig `yaml:"mongoDBConfig"`
	RedisConfig    DataStoreConfig `yaml:"redisConfig"`
	InMemoryConfig DataStoreConfig `yaml:"inMemoryConfig"`
	ZapConfig      LogConfig       `yaml:"zapConfig"`
	LorusConfig    LogConfig       `yaml:"logrusConfig"`
	Log            LogConfig       `yaml:"logConfig"`
	UseCase        UseCaseConfig   `yaml:"useCaseConfig"`
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

type UseCaseConfig struct {
	ShortenURL ShortenURLConfig `yaml:"shortenURL"`
}

type ShortenURLConfig struct {
	Code               string     `yaml:"code"`
	RedirectRepoConfig RepoConfig `yaml:"redirectRepoConfig"`
}

// RepoConfig represents handlers for data store. It can be a database or a gRPC connection
type RepoConfig struct {
	Code            string          `yaml:"code"`
	DataStoreConfig DataStoreConfig `yaml:"dataStoreConfig"`
}

// LogConfig represents logger handler
// Logger has many parameters can be set or changed. Currently, only three are listed here. Can add more into it to
// fits your needs.
type LogConfig struct {
	// log library name
	Code string `yaml:"code"`
	// log level
	Level string `yaml:"level"`
	// show caller in log message
	EnableCaller bool `yaml:"enableCaller"`
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

	return &ac, nil
}
