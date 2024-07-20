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
package identity

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

const maxThreads = 3
const lstMaxVerbosity = 1

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		var wg sync.WaitGroup
		var identityOutput strings.Builder
		var outVerbosity string
		localPart := args[0]
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domains := viper.GetStringSlice("domains")
		verbosity, err := cmd.Flags().GetCount("verbosity")
		cobra.CheckErr(err)
		boxes := make([]*utils.Wrapped[[]migagoapi.Identity], len(domains))
		wg.Add(len(domains))
		if len(args) > 1 {
			outVerbosity = args[1]
		}
		defer func() {
			if r := recover(); r != nil {
				a := []string{"min", "extra", "max"}
				fmt.Fprintf(os.Stderr, "mailbox list error: %s\n", r)
				fmt.Fprintf(os.Stderr, "You must pass one of: %+q\n", a)
				os.Exit(1)
			}
		}()
		verbosity = utils.ProcessVerboseArgs(outVerbosity, verbosity, lstMaxVerbosity)

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
				boxes[i] = utils.WrapUp(client.GetIdentities(context.Background(), localPart))
			}()
		}

		wg.Wait()

		for i, domain := range domains {
			identityOutput.WriteString(fmt.Sprintf("\nDomain: %s\n", domain))
			if boxes[i].IsErr() {
				identityOutput.WriteString(boxes[i].Err.Error())
			} else {
				listIdentities(&identityOutput, boxes[i].Get(), verbosity)
			}
		}

		fmt.Println(identityOutput.String())

	},
}

func init() {
}

func listIdentities(output *strings.Builder, identities []migagoapi.Identity, verbosity int) {
	switch verbosity {
	case 1:
		for _, box := range identities {
			out, err := json.MarshalIndent(box, "", "\t")
			cobra.CheckErr(err)
			output.Write(out)
		}
	default:
		utils.ListWithFunc(output, identities, (*migagoapi.Identity).GetAddress, "\n\t", "\n\t", "\n")
	}
}
