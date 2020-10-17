package main

import (
	"context"
	"fmt"
	"os"

	"github.com/po3rin/eskeeper"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eskeeper",
	Short: "eskeeper synchronizes index and alias with configuration files while ensuring idempotency.",
	Run: func(cmd *cobra.Command, args []string) {
		k, err := eskeeper.New(
			[]string{"http://localhost:9200"},
			eskeeper.UserName(""),
			eskeeper.Pass(""),
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

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
