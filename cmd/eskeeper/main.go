package main

import (
	"context"
	"fmt"
	"os"

	"github.com/po3rin/eskeeper"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	TERMINAL = "terminal"
	PIPE     = "pipe"
)

var rootCmd = &cobra.Command{
	Use:   "eskeeper",
	Short: "eskeeper synchronizes index and alias with configuration files while ensuring idempotency.",
	Run: func(cmd *cobra.Command, args []string) {
		k, err := eskeeper.New(
			viper.GetStringSlice("es_urls"),
			eskeeper.UserName(viper.GetString("es_user")),
			eskeeper.Pass(viper.GetString("es_pass")),
			eskeeper.Verbose(viper.GetBool("verbose")),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Println("Currently does not support interactive mode")
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

var validate = &cobra.Command{
	Use:   "validate",
	Short: "Validates config",
	Run: func(cmd *cobra.Command, args []string) {
		k, err := eskeeper.New(
			[]string{},
			eskeeper.Verbose(true),
		)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if terminal.IsTerminal(int(os.Stdin.Fd())) {
			fmt.Println("Currently does not support interactive mode")
			os.Exit(1)
		}

		ctx := context.Background()
		err = k.Validate(ctx, os.Stdin)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("pass")
	},
}

func init() {
	rootCmd.AddCommand(validate)
	viper.SetEnvPrefix("eskeeper")
	viper.AutomaticEnv()

	pflag.StringP("es_user", "u", "", "Elasticsearch user name")
	pflag.StringP("es_pass", "p", "", "Elasticsearch password")
	pflag.StringSliceP("es_urls", "e", []string{"http://localhost:9200"}, "Elasticserch endpoint URLs (comma delimited)")
	pflag.BoolP("verbose", "v", false, "Make the operation more talkative")

	viper.BindPFlags(pflag.CommandLine)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
