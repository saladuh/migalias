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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

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
		var wg sync.WaitGroup
		var mailOutput strings.Builder
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domains := viper.GetStringSlice("domains")
		verbosity, err := cmd.Flags().GetCount("verbosity")
		cobra.CheckErr(err)
		boxes := make([]utils.Wrapped[[]migagoapi.Mailbox], len(domains))
		wg.Add(len(domains))
		outVerbosity := ""
		if len(args) > 0 {
			outVerbosity = args[0]
		}
		verbosity = utils.ProcessVerboseArgs(outVerbosity, verbosity, lstMaxVerbosity)

		for i, domain := range domains {
			go func() {
				client, err := migagoapi.NewClient(&userEmail, &userToken, nil, &domain, nil)
				cobra.CheckErr(err)
				domainBoxes, err := client.GetMailboxes(context.Background())
				boxes[i] = utils.Wrapped[[]migagoapi.Mailbox]{Value: *domainBoxes, Err: err}
				wg.Done()
			}()
		}

		wg.Wait()

		for i, domain := range domains {
			mailOutput.WriteString(fmt.Sprintf("\nDomain: %s\n", domain))
			if boxes[i].IsErr() {
				mailOutput.WriteString(boxes[i].Err.Error())
			} else {
				listMailboxes(&mailOutput, boxes[i].Get(), verbosity)
			}
		}

		fmt.Println(mailOutput.String())

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

func listMailboxes(output *strings.Builder, mailboxes []migagoapi.Mailbox, verbosity int) {
	switch verbosity {
	case 1:
		utils.ListAddressesWithIdentities(output, mailboxes, "\n\t", "\n\t", "\n")
	case 2:
		for _, box := range mailboxes {
			out, err := json.MarshalIndent(box, "", "\t")
			cobra.CheckErr(err)
			output.Write(out)
		}
	default:
		utils.ListAddresses(output, mailboxes, "\n\t", "\n\t", "\n")
	}
}
