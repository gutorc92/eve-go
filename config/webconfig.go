package config

import (
	"github.com/gutorc92/api-farm/metrics"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	uri      = "uri"
	database = "database"
	logLevel = "log-level"
	files    = "files"
)

// WebConfig defines the parametric information of a data-controller server instance
type WebConfig struct {
	Metrics  *metrics.Metrics
	Database string
	Uri      string
	File     string
}

func AddFlags(flags *pflag.FlagSet) {
	flags.StringP(uri, "u", "", "Mongo db uri")
	flags.StringP(database, "d", "", "Mongo database")
	flags.StringP(logLevel, "l", "info", "[optional] The loggin level for this service")
	flags.StringP(files, "f", "", "The files to generate api service")
}

// Init initializes the web config with properties retrieved from Viper.
func (c *WebConfig) Init(v *viper.Viper) *WebConfig {
	c.Metrics = metrics.New()
	if v.GetString(database) != "" {
		c.Database = v.GetString(database)
	}
	if v.GetString(uri) != "" {
		c.Uri = v.GetString(uri)
	}
	return c
}
