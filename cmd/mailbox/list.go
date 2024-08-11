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
	"os"
	"strings"
	"sync"

	"git.sr.ht/~salad/migagoapi"
	"git.sr.ht/~salad/migalias/internal/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const lstMaxVerbosity = 2
const maxThreads = 3

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list [min|extra|max]",
		Short: "list all mailboxes of domain(s)",
		Long: `A longer description that spans multiple lines and likely contains examples
	and usage of using your command. For example:

	migalias mailbox list`,
		Args:      cobra.MaximumNArgs(1),
		ValidArgs: []string{"min", "minimal", "extra", "max", "maximum"},
		Run:       listRun,
	}

	return cmd
}

func listRun(_ *cobra.Command, args []string) {
	fmt.Println("list called")
	var wg sync.WaitGroup
	var mailOutput strings.Builder
	var outputLevel int
	var outVerbosity string
	userEmail := viper.GetString("user_email")
	userToken := viper.GetString("user_token")
	domains := viper.GetStringSlice("domains")
	boxes := make([]*utils.Wrapped[[]migagoapi.Mailbox], len(domains))
	wg.Add(len(domains))
	if len(args) > 0 {
		outVerbosity = args[0]
	}
	fmt.Println(outVerbosity)
	outputLevel, err := utils.ProcessOutputLevel(outVerbosity, lstMaxVerbosity)
	if err != nil {
		a := []string{"min", "extra", "max"}
		fmt.Fprintf(os.Stderr, "mailbox list error: %s\n", err)
		fmt.Fprintf(os.Stderr, "You must pass one of: %+q\n", a)
		os.Exit(1)
	}

	maxRoutines := make(chan int, maxThreads)
	for i, domain := range domains {
		maxRoutines <- 0
		go func() {
			defer func() {
				wg.Done()
				<-maxRoutines
			}()
			client, err := migagoapi.NewClient(userEmail, userToken, "", domain, nil)
			cobra.CheckErr(err)
			boxes[i] = utils.WrapUp(client.GetMailboxes(context.Background()))
		}()
	}

	wg.Wait()

	for i, domain := range domains {
		mailOutput.WriteString(fmt.Sprintf("\nDomain: %s\n", domain))
		if boxes[i].IsErr() {
			mailOutput.WriteString(boxes[i].Err.Error())
		} else {
			listMailboxes(&mailOutput, boxes[i].Get(), outputLevel)
		}
	}

	fmt.Println(mailOutput.String())
}

func init() {
}

func listMailboxes(output *strings.Builder, mailboxes []migagoapi.Mailbox, verbosity int) {
	switch verbosity {
	case 1:
		utils.ListWithFunc(output, mailboxes, "\n\t", "\n\t", "\n", func(m *migagoapi.Mailbox) string {
			var outString strings.Builder
			outString.WriteString(m.GetAddress())
			utils.ListWithFunc(&outString, m.Identities, "\n\t\t", "\n\t\t", "", (*migagoapi.Identity).GetAddress)
			return outString.String()
		})
		// utils.ListAddressesWithIdentities(output, mailboxes, "\n\t", "\n\t", "\n")
	case 2:
		for _, box := range mailboxes {
			out, err := json.MarshalIndent(box, "", "\t")
			cobra.CheckErr(err)
			output.Write(out)
		}
	default:
		utils.ListWithFunc(output, mailboxes, "\n\t", "\n\t", "\n", (*migagoapi.Mailbox).GetAddress)
	}
}
