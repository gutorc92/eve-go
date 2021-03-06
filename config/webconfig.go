package config

import (
	"github.com/gutorc92/api-farm/metrics"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	database      = "database"
	logLevel      = "log-level"
)

// WebConfig defines the parametric information of a data-controller server instance
type WebConfig struct {
	Metrics			 *metrics.Metrics
	Database	string
}

func AddFlags(flags *pflag.FlagSet) {
	flags.StringP(database, "d", "", "Mongo db uri")
	flags.StringP(logLevel, "l", "info", "[optional] The loggin level for this service")
}
// Init initializes the web config with properties retrieved from Viper.
func (c *WebConfig) Init(v *viper.Viper) *WebConfig {
	c.Metrics = metrics.New()
	return c
}




