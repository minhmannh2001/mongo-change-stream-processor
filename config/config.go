package config

import (
	"errors"
	"log"
	"os"

	"github.com/google/wire"
	"github.com/spf13/viper"
)

var Set = wire.NewSet(NewConfig)

// Configuration
type Configuration struct {
	Logger    Logger
	MongoDB   MongoDB
	AppConfig AppConfig
}

// AppConfig struct
type AppConfig struct {
	Name       string
	AppVersion string
	Mode       string
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// MongoDB config
type MongoDB struct {
	MongoURI        string
	MongoUser       string
	MongoPassword   string
	ConnectTimeout  int
	MaxConnIdleTime int
	MinPoolSize     uint64
	MaxPoolSize     uint64
}

// Get config path for local or docker
func getDefaultConfig() string {
	return "./config/config"
}

var DefaultConfig = Configuration{}

// Load config file from given path
func NewConfig() (*Configuration, error) {
	path := os.Getenv("cfgPath")
	if path == "" {
		path = getDefaultConfig()
	}

	v := viper.New()

	v.SetConfigName(path)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	err := v.Unmarshal(&DefaultConfig)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &DefaultConfig, nil
}
