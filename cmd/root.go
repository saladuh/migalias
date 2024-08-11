/*
Copyright © 2024 James Laverne-Cadby <james@salad.moe>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"errors"
	"log/slog"
	"os"

	"git.sr.ht/~salad/migalias/cmd/identity"
	"git.sr.ht/~salad/migalias/cmd/mailbox"
	"git.sr.ht/~salad/migalias/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const maxVerbosity = 3

// RootCmd represents the base command when called without any subcommands
func NewCommand(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migalias",
		Short: "Configure Migadu from the terminal!",
		Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,
		PersistentPreRunE: preRun,
	}

	cmd.AddCommand(
		mailbox.NewCommand(),
		identity.NewCommand(),
	)

	config.RegisterRootFlags(cmd)

	return cmd
}

func preRun(cmd *cobra.Command, _ []string) error {
	if !viper.IsSet("domains") {
		return errors.New("'domains' is not set in config or anywhere else")
	}
	doms := viper.GetViper().GetStringSlice("domains")
	if doms == nil {
		return errors.New("'domains' is not set in config or anywhere else")
	}
	if len(doms) == 0 {
		return errors.New("'domains' is not set in config or anywhere else")
	}
	verbosity, err := cmd.Flags().GetCount("verbosity")
	cobra.CheckErr(err)
	verbosity = min(verbosity, maxVerbosity)

	logVerbosity := new(slog.LevelVar)
	switch verbosity {
	default:
		logVerbosity.Set(slog.LevelError)
	case 1:
		logVerbosity.Set(slog.LevelWarn)
	case 2:
		logVerbosity.Set(slog.LevelInfo)
	case 3:
		logVerbosity.Set(slog.LevelDebug)
	}

	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logVerbosity})
	slog.SetDefault(slog.New(h))

	return nil

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	err := NewCommand(version).Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
