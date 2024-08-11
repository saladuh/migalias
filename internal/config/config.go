package config

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

func RegisterRootFlags(cmd *cobra.Command) {

	globalFlags := cmd.PersistentFlags()

	globalFlags.StringVar(&cfgFile, "config", "", "config file (default is $XDG_CONFIG_HOME/.migalias.yaml)")
	globalFlags.StringP("token", "t", "", "API token associated with Migadu account")
	globalFlags.StringP("useremail", "e", "example@example.com", "User Email of the Migadu account")
	globalFlags.StringSliceP("domains", "d", nil, "domains as comma seperated list")
	globalFlags.CountP("verbosity", "v", "-v, -vv, -vvv, to increase verbosity")

	viper.BindPFlag("user_token", globalFlags.Lookup("token"))
	viper.BindPFlag("user_email", globalFlags.Lookup("useremail"))
	viper.BindPFlag("domains", globalFlags.Lookup("domains"))

	cobra.OnInitialize(initConfig)
}

func RegisterMailboxCreateFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.StringP("password", "p", "", "Password to be used for mailbox")
	flags.StringP("invite", "i", "", "Email to send mailbox invitation to (will also be set as the recovery email)")
	cmd.MarkFlagsMutuallyExclusive("password", "invite")
	cmd.MarkFlagsOneRequired("password", "invite")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		configHome, err := os.UserConfigDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".migalias" (without extension).
		viper.AddConfigPath(configHome)
		viper.SetConfigType("yaml")
		viper.SetConfigName("migalias")
	}

	viper.SetEnvPrefix("migadu")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		// Fallback to oldschool home directory dotfile.
	} else if cfgFile == "" {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigName(".migalias")
		if err := viper.ReadInConfig(); err == nil {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
