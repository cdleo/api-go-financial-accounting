package config

import (
	"github.com/cdleo/gonfig"
)

type ClientConfig struct {
	Endpoint string
	Timeout  int
	Interval int
	APIKey   string
}

type DBServers struct {
	Host string
	Port int
}

type DBConfig struct {
	DriverName string
	Servers    []DBServers
	User       string
	Password   string //Deberia estar cifrado
	DBName     string
	Opts       string
}

// APIConfig struct
type APIConfig struct {
	ServerPort          int
	ExchangeRateService ClientConfig
	DB                  DBConfig
}

// GetAPIConfig obtiene configuracion API
func GetAPIConfig(configFilePath string) (*APIConfig, error) {

	var conf APIConfig
	if err := gonfig.GetConf(configFilePath, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}
