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

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list mailboxLocalPart [min|extra|max]",
		Short: "A brief description of your command",
		Long:  ``,
		Args:  cobra.RangeArgs(1, 2),
		Run:   listRun,
	}

	return cmd
}

func init() {
}

func listRun(cmd *cobra.Command, args []string) {
	fmt.Println("list called")
	var wg sync.WaitGroup
	var identityOutput strings.Builder
	var outVerbosity string
	var outputLevel int
	localPart := args[0]
	userEmail := viper.GetString("user_email")
	userToken := viper.GetString("user_token")
	domains := viper.GetStringSlice("domains")
	boxes := make([]*utils.Wrapped[[]migagoapi.Identity], len(domains))
	wg.Add(len(domains))
	if len(args) > 1 {
		outVerbosity = args[1]
	}
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
			boxes[i] = utils.WrapUp(client.GetIdentities(context.Background(), localPart))
		}()
	}

	wg.Wait()

	for i, domain := range domains {
		identityOutput.WriteString(fmt.Sprintf("\nDomain: %s\n", domain))
		if boxes[i].IsErr() {
			identityOutput.WriteString(boxes[i].Err.Error())
		} else {
			listIdentities(&identityOutput, boxes[i].Get(), outputLevel)
		}
	}

	fmt.Println(identityOutput.String())
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
		utils.ListWithFunc(output, identities, "\n\t", "\n\t", "\n", (*migagoapi.Identity).GetAddress)
	}
}
