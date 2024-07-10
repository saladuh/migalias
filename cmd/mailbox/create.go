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
	"errors"
	"fmt"

	"git.sr.ht/~salad/migagoapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:   "create [local part] [name]",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		// var mailOutput strings.Builder
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domain := viper.GetStringSlice("domains")[0]
		localPart := args[0]
		name := args[1]
		client, err := migagoapi.NewClient(&userEmail, &userToken, nil, &domain, nil)
		cobra.CheckErr(err)
		if pass, err := cmd.Flags().GetString("password"); err == nil && pass != "" {
			newMailbox, err := client.CreateMailboxWithPassword(context.Background(), name, localPart, pass, false)
			cobra.CheckErr(err)
			fmt.Printf("The account %s@%s has been created sucessfully.\nMigadu Returned:\n", localPart, domain)
			out, err := json.MarshalIndent(newMailbox, "", "\t")
			cobra.CheckErr(err)
			fmt.Println(string(out))
		} else if err != nil {
			cobra.CheckErr(errors.New("Password flag non existent"))
		}
	},
}

func init() {
	createCmd.Flags().StringP("password", "p", "", "Password to be used for mailbox")
	createCmd.Flags().StringP("invite", "i", "", "Email to send mailbox invitation to (will also be set as the recovery email)")
	createCmd.MarkFlagsMutuallyExclusive("password", "invite")
	createCmd.MarkFlagsOneRequired("password", "invite")
}
