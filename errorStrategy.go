package obfuscate

import (
	"fmt"
	"log"
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
	case OnErrorDiscard:
		return "OnErrorDiscard"
	}

	return fmt.Sprintf("ErrorStrategy(%d)", int(es))
}

const (
	// OnErrorLog will cause obfuscation to stop when an error occurs. The error will be logged.
	OnErrorLog ErrorStrategy = iota
	// OnErrorInclude will cause obfuscation to stop when an error occurs. The error will be included in the obfuscation result.
	OnErrorInclude
	// OnErrorDiscard will cause obfuscation to stop when an error occurs. The error will not be visible in any way.
	OnErrorDiscard
)

func logErrorf(logger *log.Logger, format string, v ...any) {
	if logger == nil {
		fmt.Printf(format, v...)
	} else {
		logger.Printf(format, v...)
	}
}
