// Package cmd provides command line processing functions
// for the authentication services.
package cmd

import (
	"fmt"
	"os"

	"github.com/dhaifley/dauth/lib"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   lib.ServiceInfo.Name,
	Short: lib.ServiceInfo.Short,
	Long:  lib.ServiceInfo.Long,
}

func init() {
	viper.SetConfigFile("dauth_config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}

	viper.SetEnvPrefix(lib.ServiceInfo.Name)
	viper.SetDefault("sql", "")
	if err := viper.BindEnv("sql"); err != nil {
		fmt.Println(err)
	}

	viper.SetDefault("cert", "")
	if err := viper.BindEnv("cert"); err != nil {
		fmt.Println(err)
	}

	viper.SetDefault("key", "")
	if err := viper.BindEnv("key"); err != nil {
		fmt.Println(err)
	}
}

// Execute starts the command processor.
func Execute() {
	fmt.Println(lib.ServiceInfo.Short)
	fmt.Println("Version:", lib.ServiceInfo.Version)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
