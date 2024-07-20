package utils

import (
	"fmt"
	"strings"
)

// TODO: Do something better here. This is awful. Such a hacky workaround to the whole value not implementing the generic thing
// (where the pointer does, but the value doesn't, and the slice contains values, not pointers)
func ListWithFunc[M ~[]A, A any](output *strings.Builder, mailObs M, getFunc func(*A) string, delimiter, starter, ender string) {
	if len(mailObs) == 0 {
		return
	}
	output.WriteString(starter)
	for i, o := range mailObs {
		output.WriteString(getFunc(&o))
		if i != len(mailObs)-1 {
			output.WriteString(delimiter)
		}
	}
	output.WriteString(ender)
}

type Wrapped[T any] struct {
	Value T
	Err   error
}

func WrapUp[T any](value T, err error) *Wrapped[T] {
	return &Wrapped[T]{
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
			panic(fmt.Sprintf("Information density verbosity argument invalid: %s\n", verboseArg))
		}
	}

	outputVerbosity = min(max(outputVerbosity, verbosity), maxVerbosity)
	return outputVerbosity
}
