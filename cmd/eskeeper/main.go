package main

import (
	"context"
	"fmt"
	"os"

	"github.com/po3rin/eskeeper"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "eskeeper",
	Short: "eskeeper synchronizes index and alias with configuration files while ensuring idempotency.",
	Run: func(cmd *cobra.Command, args []string) {
		k, err := eskeeper.New(
			viper.GetStringSlice("es_urls"),
			eskeeper.UserName(viper.GetString("es_user")),
			eskeeper.Pass(viper.GetString("es_pass")),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		ctx := context.Background()
		err = k.Sync(ctx, os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	viper.SetEnvPrefix("eskeeper")
	viper.AutomaticEnv()

	pflag.StringP("es_user", "u", "", "Elasticsearch user name")
	viper.BindPFlags(pflag.CommandLine)

	pflag.StringP("es_pass", "p", "", "Elasticsearch password")
	viper.BindPFlags(pflag.CommandLine)

	pflag.StringSliceP("es_urls", "e", []string{"http://localhost:9200"}, "Elasticserch endpoint URLs (comma delimited)")
	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
