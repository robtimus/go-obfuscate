package obfuscate

import (
	"log"
	"strings"
)

// All returns an obfuscator that replaces all characters in strings with an asterisk (*).
func All() Obfuscator {
	return AllWithMask("*")
}

// AllWithMask returns an obfuscator that replaces all characters in strings with the given mask.
func AllWithMask(mask string) Obfuscator {
	return NewObfuscator(func(s string) string {
		return strings.Repeat(mask, len(s))
	})
}

var none = NewObfuscator(func(s string) string {
	return s
})

// None returns an obfuscator that does not obfuscate anything.
// It can be used as default value.
func None() Obfuscator {
	return none
}

// WithFixedLength returns an obfuscator that replaces strings with the given fixed length occurrences of a single asterisk (*).
//
// It panics if the fixed length is negative.
func WithFixedLength(fixedLength int) Obfuscator {
	return WithFixedLengthWithMask(fixedLength, "*")
}

// WithFixedLengthWithMask returns an obfuscator that replaces strings with the given fixed length occurrences of the given mask.
//
// It panics if the fixed length is negative or if the result of (len(mask) * fixedLength) overflows.
func WithFixedLengthWithMask(fixedLength int, mask string) Obfuscator {
	if fixedLength < 0 {
		log.Panicf("fixedLength: %d < 0", fixedLength)
	}
	fixedValue := strings.Repeat(mask, fixedLength)
	return WithFixedValue(fixedValue)
}

// WithFixedValue returns an obfuscator that replaces strings with the given fixed value.
func WithFixedValue(fixedValue string) Obfuscator {
	return NewObfuscator(func(s string) string {
		return fixedValue
	})
}
