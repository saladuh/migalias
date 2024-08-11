package utils

import (
	"errors"
	"strings"
)

// TODO: Do something better here. This is awful. Such a hacky workaround to the whole value not implementing the generic thing
// (where the pointer does, but the value doesn't, and the slice contains values, not pointers)
func ListWithFunc[M ~[]A, A any](output *strings.Builder, l M, delimiter, starter, ender string, getS func(*A) string) {
	if len(l) == 0 {
		return
	}
	output.WriteString(starter)
	for i, o := range l {
		output.WriteString(getS(&o))
		if i != len(l)-1 {
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

func ProcessOutputLevel(verboseArg string, maxVerbosity int) (int, error) {
	var outputVerbosity int
	var err error = nil
	switch verboseArg {
	case "", "min", "minimal":
		outputVerbosity = 0
	case "extra":
		outputVerbosity = 1
	case "max", "maximum":
		outputVerbosity = 2
	default:
		outputVerbosity = 0
		err = errors.New("Output verbosiy invalid")
	}

	outputVerbosity = min(outputVerbosity, maxVerbosity)
	return outputVerbosity, err
}
