package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func execute() error {
	rootCmd := &cobra.Command{
		Use:   "kg",
		Short: "Knowledge Graph Manager",
		Long:  `A CLI tool to manage a knowledge graph using markdown files with YAML frontmatter.`,
	}

	rootCmd.PersistentFlags().String("config", "", "config file (default is $HOME/.kgrc)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))

	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(
		newConnectCmd(),
		newListCmd(),
		newSearchCmd(),
		newAddCmd(),
		newEditCmd(),
		newVisualizeCmd(),
		newStatsCmd(),
		newExportCmd(),
		newImportCmd(),
		newBackupCmd(),
		newConfigCmd(),
		newFrontmatterCmd(),
	)

	return rootCmd.Execute()
}

func initConfig() {
	if cfgFile := viper.GetString("config"); cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath("..")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".kgrc")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
