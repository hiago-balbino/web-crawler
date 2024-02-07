package cmd

import (
	"github.com/hiago-balbino/web-crawler/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "A main CLI to crawler",
	Long:  "CLI used to crawl the HTML page and return the slice of links given a depth",
}

func Execute() error {
	cobra.OnInitialize(config.InitConfigurations)
	rootCmd.AddCommand(apiCmd)

	return rootCmd.Execute()
}
