package utils

import (
	"git.sr.ht/~salad/migagoapi"
	"strings"
)

type Wrapped[T any] struct {
	Value T
	Err   error
}

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

func WrapUp[T any](value T, err error) Wrapped[T] {
	return Wrapped[T]{
		Value: value,
		Err:   err,
	}
}

func (wrap *Wrapped[T]) IsErr() bool {
	return wrap.Err != nil
}

func (wrap *Wrapped[T]) Get() T {
	return wrap.Value
}

func (wrap *Wrapped[T]) GetOrPanic() T {
	if err := wrap.Err; err != nil {
		panic(err)
	}
	return wrap.Value
}

func ProcessVerboseArgs(verboseArg string, verbosity int, maxVerbosity int) int {
	var outputVerbosity int
	if verboseArg == "" {
		outputVerbosity = 0
	} else {
		switch verboseArg {
		case "min", "minimal":
			outputVerbosity = 0
		case "extra":
			outputVerbosity = 1
		case "max", "maximum":
			outputVerbosity = 2
		default:
			panic("What the frick\n")
		}
	}

	outputVerbosity = max(outputVerbosity, min(verbosity, maxVerbosity))
	return outputVerbosity
}
