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
package mailbox

import (
	"context"
	"fmt"
	"strings"

	"git.sr.ht/~salad/migagoapi"
	"git.sr.ht/~salad/migalias/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const lstMaxVerbosity = 2

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all mailboxes of domain(s)",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

migalias mailbox list`,
	Args:      cobra.MaximumNArgs(1),
	ValidArgs: []string{"min", "minimal", "extra", "max", "maximum"},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domains := viper.GetStringSlice("domains")
		verbosity, err := cmd.Flags().GetCount("verbosity")
		cobra.CheckErr(err)
		verbosity = processListArgs(args, verbosity)
		// client, err := migagoapi.NewClient(&userEmail, &userToken,
		// 	nil, &viper.GetStringSlice("domains")[0], nil)
		// cobra.CheckErr(err)
		// fmt.Println(client.GetMailboxes(context.Background()))
		for _, domain := range domains {
			fmt.Printf("\nDomain: %s\n", domain)
			client, err := migagoapi.NewClient(&userEmail, &userToken, nil, &domain, nil)
			cobra.CheckErr(err)
			domainMailboxes, err := client.GetMailboxes(context.Background())
			cobra.CheckErr(err)
			listMailboxes(*domainMailboxes, verbosity)
		}
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listMailboxes(mailboxes []migagoapi.Mailbox, verbosity int) {
	var output strings.Builder
	switch verbosity {
	case 1:
		utils.ListAddressesWithIdentities(&output, mailboxes, "\n\t", "\n\t", "\n")
	default:
		utils.ListAddresses(&output, mailboxes, "\n\t", "\n\t", "\n")
	}
	fmt.Print(output.String())
}

func processListArgs(args []string, verbosity int) int {
	fmt.Println(args)
	var outputVerbosity int
	if len(args) == 0 {
		outputVerbosity = 0
	} else {
		switch verboseLevel := args[0]; verboseLevel {
		case "min", "minimal":
			outputVerbosity = 0
		case "extra":
			outputVerbosity = 1
		case "max", "maximum":
			outputVerbosity = 2
		default:
			panic("What the frick\n")
		}
	}

	outputVerbosity = max(outputVerbosity, min(verbosity, lstMaxVerbosity))
	return outputVerbosity
}
