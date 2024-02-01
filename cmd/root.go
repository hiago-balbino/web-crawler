package cmd

import (
	configuration "github.com/hiago-balbino/web-crawler/configs"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "A main CLI to crawler",
	Long:  "CLI used to crawl the HTML page and return the slice of links given a depth",
}

// Execute executes the root command.
func Execute() error {
	cobra.OnInitialize(configuration.InitConfigurations)
	rootCmd.AddCommand(apiCmd)

	return rootCmd.Execute()
}
