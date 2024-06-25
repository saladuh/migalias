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
	"fmt"

	"git.sr.ht/~salad/migalias/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// mailboxCmd represents the mailbox command
var mailboxCmd = &cobra.Command{
	Use:   "mailbox",
	Short: "Mailbox related commands",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

migalias mailbox [options]`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mailbox called")
		if domains := viper.GetStringSlice("domains"); domains != nil {
			fmt.Println("There are some domains!")
			fmt.Println(domains)
		}
	},
}

func init() {
	cmd.RootCmd.AddCommand(mailboxCmd)
	mailboxCmd.AddCommand(listCmd)

	// mailboxCmd.PersistentFlags().BoolP("test", "w", false, "test")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mailboxCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mailboxCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
