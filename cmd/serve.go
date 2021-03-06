package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/gutorc92/api-farm/config"
	"github.com/gutorc92/api-farm/handlers"

)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the HTTP REST APIs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		wc := new(config.WebConfig).Init(viper.GetViper())
		server := new(handlers.Server).InitFromWebConfig(wc)

		err := server.Serve()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	config.AddFlags(serveCmd.Flags())

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
