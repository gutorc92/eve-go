package cmd

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/gutorc92/eve-go/config"
	"github.com/gutorc92/eve-go/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", name.Name), nil
}

// serveCmd represents the serve command
var labdaCmd = &cobra.Command{
	Use:   "lambda",
	Short: "Starts the AWS Lambda APIs server",
	RunE: func(cmd *cobra.Command, args []string) error {
		wc := new(config.WebConfig).Init(viper.GetViper())
		fmt.Println(wc.Uri)
		server := new(handlers.Server).InitFromWebConfig(wc)
		lambda.Start(server.HandleAsLambda)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(labdaCmd)

	config.AddFlags(labdaCmd.Flags())

	err := viper.GetViper().BindPFlags(labdaCmd.Flags())
	if err != nil {
		panic(err)
	}
}
