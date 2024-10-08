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

	"git.sr.ht/~salad/migagoapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete local_part",
		Short: "",
		Long:  ``,
		Run:   deleteRun,
	}
	return cmd
}

func deleteRun(cmd *cobra.Command, args []string) {
	fmt.Println("delete called")
	userEmail := viper.GetString("user_email")
	userToken := viper.GetString("user_token")
	domain := viper.GetStringSlice("domains")[0]
	localPart := args[0]
	client, err := migagoapi.NewClient(userEmail, userToken, "", domain, nil)
	cobra.CheckErr(err)

	fmt.Printf("Deleting mailbox %s@%s ...\n", localPart, domain)
	err = client.DeleteMailbox(context.Background(), localPart)
	if err != nil {
		fmt.Printf("An error occured while deleting %s@%s, see below.\n", localPart, domain)
		panic(err)
	} else {
		fmt.Printf("Successfully deleted mailbox %s@%s\n", localPart, domain)
	}
}

func init() {

}
