package config

import (
	"encoding/json"
	"github.com/spf13/viper"
	log "github.com/sirupsen/logrus"
)

type GClusterConfig interface {
	ConfigString() string
}

type LogConfiguration interface {
	GetLogLevelConfig() string
}

type ServiceRegistryConfig struct {
	Address string
}

type ServiceRegistryConfiguration interface {
	GetServiceRegistryConfig() *ServiceRegistryConfig
}

type OpenTracingConfig struct {
	Address string
}

type OpenTracingConfiguration interface {
	GetOpenTracingConfig() *OpenTracingConfig
}

type DatabaseConfig struct {
	Address     string
	LogMode     bool
	MaxIdle     int
	MaxConns    int
	MaxLifetime int
}

type DatabaseConfiguration interface {
	GetDataBaseConfiguration() *DatabaseConfig
}

type ServerConfig struct {
	Address string
	Port    int
}

type ServerConfiguration interface {
	GetServerConfig() *ServerConfig
}

type GClusterConfigLoader struct {
	Name     string
	FilePath string
	Config   GClusterConfig
}

func (loader *GClusterConfigLoader) Load() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("GClusterConfigLoader ReadConfig failed, err: %v", err)
		return err
	}

	if err := viper.Unmarshal(loader.Config); err != nil {
		log.Errorf("GClusterConfigLoader ReadConfig failed %v", err)
		return err
	}

	configData, _ := json.MarshalIndent(viper.AllSettings(), "", "")
	log.Printf("GClusterConfigLoader ReadRemoteConfig success, %s", configData)

	return nil
}
