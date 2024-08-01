package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage .kgrc settings",
		Long:  `View, set, or unset configuration options in the .kgrc file.`,
	}

	cmd.AddCommand(
		newConfigSetCmd(),
		newConfigGetCmd(),
		newConfigListCmd(),
		newConfigUnsetCmd(),
	)

	return cmd
}

func newConfigSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Set a configuration option",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setConfig(args[0], args[1])
		},
	}
}

func newConfigGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get [key]",
		Short: "Get a configuration option",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getConfig(args[0])
		},
	}
}

func newConfigListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configuration options",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return listConfig()
		},
	}
}

func newConfigUnsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "unset [key]",
		Short: "Unset a configuration option",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return unsetConfig(args[0])
		},
	}
}

func setConfig(key, value string) error {
	if err := validateConfigKey(key); err != nil {
		return err
	}

	viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Set %s = %s\n", key, value)
	return nil
}

func getConfig(key string) error {
	if err := validateConfigKey(key); err != nil {
		return err
	}

	value := viper.Get(key)
	if value == nil {
		return fmt.Errorf("key '%s' not found in config", key)
	}

	fmt.Printf("%s = %v\n", key, value)
	return nil
}

func listConfig() error {
	settings := viper.AllSettings()
	for key, value := range settings {
		fmt.Printf("%s = %v\n", key, value)
	}
	return nil
}

func unsetConfig(key string) error {
	if err := validateConfigKey(key); err != nil {
		return err
	}

	viper.Set(key, nil)
	if err := viper.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	fmt.Printf("Unset %s\n", key)
	return nil
}

func validateConfigKey(key string) error {
	validKeys := []string{
		"knowledge_graph_dir",
		"backup_dir",
		"max_backups",
		"editor",
		"default_tags",
		"date_format",
	}

	key = strings.ToLower(key)
	for _, validKey := range validKeys {
		if key == validKey {
			return nil
		}
	}

	return fmt.Errorf("invalid config key: %s", key)
}