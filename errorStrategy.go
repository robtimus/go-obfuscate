package obfuscate

import (
	"fmt"
)

// ErrorStrategy represents the strategy to follow when an error occurs while obfuscating a string.
type ErrorStrategy int

// String implements the fmt.Stringer interface.
func (es ErrorStrategy) String() string {
	switch es {
	case OnErrorLog:
		return "OnErrorLog"
	case OnErrorInclude:
		return "OnErrorInclude"
	case OnErrorStop:
		return "OnErrorStop"
	}
	return fmt.Sprintf("ErrorStrategy(%d)", int(es))
}

const (
	// OnErrorLog will cause obfuscation to stop when an error occurs. The error will be logged.
	OnErrorLog ErrorStrategy = iota
	// OnErrorInclude will cause obfuscation to stop when an error occurs. The error will be included in the obfuscation result.
	OnErrorInclude
	// OnErrorStop will cause obfuscation to stop when an error occurs. The error will not be visible in any way.
	OnErrorStop
)

var defaultPrintf = func(format string, v ...any) {
	fmt.Printf(format, v...)
}
