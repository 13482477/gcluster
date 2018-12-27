package config

import (
	"gcluster/essential/config"
	"encoding/json"
)

type UserConfig struct {
	LogLevel              string
	ServiceRegistryConfig *config.ServiceRegistryConfig
	OpenTracingConfig     *config.OpenTracingConfig
	DatabaseConfig        *config.DatabaseConfig
	ServerConfig          *config.ServerConfig
}

func (cfg *UserConfig) GetLogLevelConfig() string {
	return cfg.LogLevel
}

func (cfg *UserConfig) GetServiceRegistryConfig() *config.ServiceRegistryConfig {
	return cfg.ServiceRegistryConfig
}

func (cfg *UserConfig) GetOpenTracingConfig() *config.OpenTracingConfig {
	return cfg.OpenTracingConfig
}

func (cfg *UserConfig) GetDataBaseConfig() *config.DatabaseConfig {
	return cfg.DatabaseConfig
}

func (cfg *UserConfig) GetServerConfig() *config.ServerConfig {
	return cfg.ServerConfig
}

func (cfg *UserConfig) ConfigString() string {
	bytes, _ := json.Marshal(cfg)
	return string(bytes)
}
