package config

import (
	"github.com/gutorc92/api-farm/metrics"
)

const (
	bbAPIURL      = "bb-api-url"
	port          = "port"
	logLevel      = "log-level"
	redisHost     = "redis-host"
	redisPassword = "redis-password"
	redisDb       = "redis-db"
	shutdownTime  = "shutdown-time"
	jsonFile      = "json-file"
	userName      = "username"
	userPassword  = "password"
)

// WebConfig defines the parametric information of a data-controller server instance
type WebConfig struct {
	teste string
	Metrics			 *metrics.Metrics
}

// Init initializes the web config with properties retrieved from Viper.
func (c *WebConfig) Init() *WebConfig {
	c.teste = "testando"
	c.Metrics = metrics.New()
	return c
}


