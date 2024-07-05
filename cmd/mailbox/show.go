package mailbox

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"git.sr.ht/~salad/migagoapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show mailbox attributes",
	Long: `Show will output any and up to all the attributes of a mailbox, as
		available through the Migadu api

		When passing domains or using a config file with multiple domains
		configured, the command will default to the first listed domain.
		In the case of a config file, this can be overridden by passing
		a single domain to the --domains/-d flag`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		var mailOutput strings.Builder
		userEmail := viper.GetString("user_email")
		userToken := viper.GetString("user_token")
		domains := viper.GetStringSlice("domains")
		localPart := args[0]
		// verbosity, err := cmd.Flags().GetCount("verbosity")
		// cobra.CheckErr(err)
		domain := domains[0]
		client, err := migagoapi.NewClient(&userEmail, &userToken, nil, &domain, nil)
		cobra.CheckErr(err)
		box, err := client.GetMailbox(context.Background(), localPart)
		cobra.CheckErr(err)
		out, err := json.MarshalIndent(box, "", "\t")
		cobra.CheckErr(err)
		mailOutput.Write(out)
		fmt.Println(mailOutput.String())
	},
}

func init() {

}
