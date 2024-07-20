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

	"git.sr.ht/~salad/migagoapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	password = iota + 1
	invitationEmail
)

var authType int

var createCmd = &cobra.Command{
	Use:   "create local_part name",
	Short: "",
	Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		if cmd.Flags().Changed("password") {
			authType = password
		} else if cmd.Flags().Changed("invite") {
			authType = invitationEmail
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domain := viper.GetStringSlice("domains")[0]
		localPart := args[0]
		name := args[1]
		client, err := migagoapi.NewClient(userEmail, userToken, "", domain, nil)
		cobra.CheckErr(err)

		switch authType {
		case password:
			pass, _ := cmd.Flags().GetString("password")
			newMailbox, err := client.CreateMailboxWithPassword(context.Background(), name, localPart, pass, false)
			cobra.CheckErr(err)
			fmt.Printf("The account %s@%s has been created sucessfully.\nMigadu Returned:\n", localPart, domain)
			out, err := json.MarshalIndent(newMailbox, "", "\t")
			cobra.CheckErr(err)
			fmt.Println(string(out))
		case invitationEmail:
			iEmail, _ := cmd.Flags().GetString("invite")
			newMailbox, err := client.CreateMailboxWithInvite(context.Background(), name, localPart, iEmail)
			cobra.CheckErr(err)
			fmt.Printf("The account %s@%s has been created sucessfully.\nMigadu Returned:\n", localPart, domain)
			out, err := json.MarshalIndent(newMailbox, "", "\t")
			cobra.CheckErr(err)
			fmt.Println(string(out))
			fmt.Printf("Access the email account %s to reset the password.\n", iEmail)
		default:
			panic("mailbox create: Something went very wrong with either your password or invitation email")
		}
	},
}

func init() {
	createCmd.Flags().StringP("password", "p", "", "Password to be used for mailbox")
	createCmd.Flags().StringP("invite", "i", "", "Email to send mailbox invitation to (will also be set as the recovery email)")
	createCmd.MarkFlagsMutuallyExclusive("password", "invite")
	createCmd.MarkFlagsOneRequired("password", "invite")
}
