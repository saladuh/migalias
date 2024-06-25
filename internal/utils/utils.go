package utils

import (
	"git.sr.ht/~salad/migagoapi"
	"strings"
)

func ListAddresses[M ~[]O, O migagoapi.Addresser](output *strings.Builder, mailObs M, delimiter, starter, ender string) *strings.Builder {
	if len(mailObs) == 0 {
		return output
	}
	output.WriteString(starter)
	for i, o := range mailObs {
		output.WriteString(o.GetAddress())
		if i != len(mailObs)-1 {
			output.WriteString(delimiter)
		}
	}
	output.WriteString(ender)
	return output
}

func ListAddressesWithIdentities(
	output *strings.Builder, mailObs []migagoapi.Mailbox, delimiter, starter, ender string) *strings.Builder {
	if len(mailObs) == 0 {
		return output
	}
	output.WriteString(starter)
	for i, o := range mailObs {
		output.WriteString(o.GetAddress())
		ListAddresses(output, o.Identities, "\n\t\t", "\n\t\t", "")
		if i != len(mailObs)-1 {
			output.WriteString(delimiter)
		}
	}
	output.WriteString(ender)
	return output
}
