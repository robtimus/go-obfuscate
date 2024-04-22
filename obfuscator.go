package obfuscate

import "log"

// Obfuscator represents an object that can obfuscate strings, making them partly or completely unreadable.
type Obfuscator interface {
	// ObfuscateString obfuscates the given string.
	ObfuscateString(s string) string

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
	obfuscate       func(s string) string
	minPrefixLength int
}

func (o obfuscator) ObfuscateString(s string) string {
	return o.obfuscate(s)
}

func (o obfuscator) UntilLength(prefixLength int) ObfuscatorPrefix {
	if prefixLength < o.minPrefixLength {
		log.Panicf("prefixLength: %d < %d", prefixLength, o.minPrefixLength)
	}
	return ObfuscatorPrefix{o, prefixLength}
}

// NewObfuscator creates a new obfuscator that delegates to the given function.
func NewObfuscator(obfuscate func(s string) string) Obfuscator {
	return newObfuscator(obfuscate, 1)
}

func newObfuscator(obfuscate func(s string) string, minPrefixLength int) Obfuscator {
	return obfuscator{obfuscate, minPrefixLength}
}

// ObfuscatorPrefix represents a prefix of a specific length that uses a specific obfuscator.
// It can be used to create combined obfuscators that obfuscate strings for the part up to the length of this prefix using the prefix' obfuscator,
// then the rest with another.
type ObfuscatorPrefix struct {
	obfuscator   Obfuscator
	prefixLength int
}

// Then returns an obfuscator that first uses the obfuscator that was used to create the receiver for the length of the receiver,
// then another obfuscator.
func (op ObfuscatorPrefix) Then(other Obfuscator) Obfuscator {
	first := op.obfuscator
	lengthForFirst := op.prefixLength
	second := other
	return newObfuscator(func(s string) string {
		if len(s) <= lengthForFirst {
			return first.ObfuscateString(s)
		}
		return first.ObfuscateString(s[:lengthForFirst]) + second.ObfuscateString(s[lengthForFirst:])
	}, lengthForFirst+1)
}
