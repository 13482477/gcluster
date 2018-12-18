package config

import (
	"mcloud/public.v2/config"
	"encoding/json"
)

type EcsConfig struct {
	LogLevel              string
	ServiceRegistryConfig *config.ServiceRegistryConfig
	OpenTracingConfig     *config.OpenTracingConfig
	DatabaseConfig        *config.DatabaseConfig
	ServerConfig          *config.ServerConfig
}

func (cfg *EcsConfig) GetLogLevelConfig() string {
	return cfg.LogLevel
}

func (cfg *EcsConfig) GetServiceRegistryConfig() *config.ServiceRegistryConfig {
	return cfg.ServiceRegistryConfig
}

func (cfg *EcsConfig) GetOpenTracingConfig() *config.OpenTracingConfig {
	return cfg.OpenTracingConfig
}

func (cfg *EcsConfig) GetDataBaseConfiguration() *config.DatabaseConfig {
	return cfg.DatabaseConfig
}

func (cfg *EcsConfig) GetServerConfig() *config.ServerConfig {
	return cfg.ServerConfig
}

func (cfg *EcsConfig) ConfigString() string {
	bytes, _ := json.Marshal(cfg)
	return string(bytes)
}
