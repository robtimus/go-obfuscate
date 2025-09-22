package obfuscate

import (
	"fmt"
	"log"
)

// Obfuscator represents an object that can obfuscate strings, making them partly or completely unreadable.
type Obfuscator interface {
	// ObfuscateString obfuscates the given string.
	//
	// If the parser needs to do any parsing that can result in an error, this should be handled according to one of the possible [ErrorStrategy] constants.
	ObfuscateString(s string) string

	// ParseAndObfuscateString parses the given string and obfuscates the parsed result.
	//
	// If the obfuscator does not need any parsing, this method will do the same as [ObfuscateString], and the error will be nil.
	ParseAndObfuscateString(s string) (string, error)

	// UntilLength creates a prefix that can be used to chain another obfuscator to this obfuscator.
	// For the part up to the given prefix length, this obfuscator will be used; for any remaining content another obfuscator will be used.
	// This makes it possible to easily create complex obfuscators that would otherwise be impossible using any of the other obfuscators provided by this module.
	//
	// The prefix length needs to be at least 1, and larger than all previous lengths in a method chain.
	// In other words, each prefix length must be larger than its direct predecessor.
	// This method panics if this pre-condition is not met.
	UntilLength(prefixLength int) ObfuscatorPrefix
}

type obfuscator struct {
	obfuscate         func(s string) string
	parseAndObfuscate func(s string) (string, error)
	minPrefixLength   int
}

func (o obfuscator) ObfuscateString(s string) string {
	return o.obfuscate(s)
}

func (o obfuscator) ParseAndObfuscateString(s string) (string, error) {
	return o.parseAndObfuscate(s)
}

func (o obfuscator) UntilLength(prefixLength int) ObfuscatorPrefix {
	return NewObfuscatorPrefix(o, prefixLength)
}

// NewObfuscator creates a new obfuscator that delegates to the given function.
//
// Calling [Obfuscator.ParseAndObfuscateString] on the result will always return a nil error.
func NewObfuscator(obfuscate func(s string) string) Obfuscator {
	parseAndObfuscate := func(s string) (string, error) {
		return obfuscate(s), nil
	}

	return obfuscator{obfuscate: obfuscate, parseAndObfuscate: parseAndObfuscate, minPrefixLength: 1}
}

// NewObfuscatorOnErrorLog creates a new obfuscator that delegates to the given function.
//
// When calling [Obfuscator.ObfuscateString] on the result with a string that would cause the given function to return an error, this error is logged.
// If the given logger is nil, [fmt.Printf] will be used instead.
func NewObfuscatorOnErrorLog(parseAndObfuscate func(s string) (string, error), logger *log.Logger) Obfuscator {
	obfuscate := func(s string) string {
		result, err := parseAndObfuscate(s)
		if err != nil {
			logErrorf(logger, "ObfuscateString error: %v\n", err)
		}

		return result
	}

	return obfuscator{obfuscate: obfuscate, parseAndObfuscate: parseAndObfuscate, minPrefixLength: 1}
}

// NewObfuscatorOnErrorInclude creates a new obfuscator that delegates to the given function.
//
// When calling [Obfuscator.ObfuscateString] on the result with a string that would cause the given function to return an error, this error is appended to the result.
func NewObfuscatorOnErrorInclude(parseAndObfuscate func(s string) (string, error)) Obfuscator {
	obfuscate := func(s string) string {
		result, err := parseAndObfuscate(s)
		if err != nil {
			result = fmt.Sprintf("%s<error: %v>", result, err)
		}

		return result
	}

	return obfuscator{obfuscate: obfuscate, parseAndObfuscate: parseAndObfuscate, minPrefixLength: 1}
}

// NewObfuscatorOnErrorDiscard creates a new obfuscator that delegates to the given function.
//
// When calling [Obfuscator.ObfuscateString] on the result with a string that would cause the given function to return an error, this error is discarded.
func NewObfuscatorOnErrorDiscard(parseAndObfuscate func(s string) (string, error)) Obfuscator {
	obfuscate := func(s string) string {
		result, _ := parseAndObfuscate(s)

		return result
	}

	return obfuscator{obfuscate: obfuscate, parseAndObfuscate: parseAndObfuscate, minPrefixLength: 1}
}

// ObfuscatorPrefix represents a prefix of a specific length that uses a specific obfuscator.
// It can be used to create combined obfuscators that obfuscate strings for the part up to the length of this prefix using the prefix' obfuscator,
// then the rest with another.
type ObfuscatorPrefix struct {
	obfuscator   Obfuscator
	prefixLength int
}

// NewObfuscatorPrefix creates a new ObfuscatorPrefix.
//
// This function is intended only for implementing custom obfuscators.
// It should not be used directly.
func NewObfuscatorPrefix(o Obfuscator, prefixLength int) ObfuscatorPrefix {
	minPrefixLength := 1
	if obf, ok := o.(obfuscator); ok {
		minPrefixLength = obf.minPrefixLength
	}

	if prefixLength < minPrefixLength {
		log.Panicf("prefixLength: %d < %d", prefixLength, minPrefixLength)
	}

	return ObfuscatorPrefix{o, prefixLength}
}

// Then returns an obfuscator that first uses the obfuscator that was used to create the receiver for the length of the receiver,
// then another obfuscator.
func (op ObfuscatorPrefix) Then(other Obfuscator) Obfuscator {
	first := op.obfuscator
	lengthForFirst := op.prefixLength
	second := other

	obfuscate := func(s string) string {
		if len(s) <= lengthForFirst {
			return first.ObfuscateString(s)
		}

		return first.ObfuscateString(s[:lengthForFirst]) + second.ObfuscateString(s[lengthForFirst:])
	}
	parseAndObfuscate := func(s string) (string, error) {
		if len(s) <= lengthForFirst {
			return first.ParseAndObfuscateString(s)
		}

		result1, err := first.ParseAndObfuscateString(s[:lengthForFirst])
		if err != nil {
			return result1, err
		}

		result2, err := second.ParseAndObfuscateString(s[lengthForFirst:])

		return result1 + result2, err
	}

	return obfuscator{obfuscate: obfuscate, parseAndObfuscate: parseAndObfuscate, minPrefixLength: lengthForFirst + 1}
}
