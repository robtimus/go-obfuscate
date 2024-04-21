package obfuscate

import (
	"log"
	"strings"
)

// Portion returns a builder for obfuscators that obfuscates a specific portion of their input.
func Portion() *PortionBuilder {
	return &PortionBuilder{}
}

// PortionBuilder represents a builder for obfuscators that obfuscates a specific portion of their input.
type PortionBuilder struct {
	keepAtStart      *int
	keepAtEnd        *int
	atLeastFromStart *int
	atLeastFromEnd   *int
	fixedTotalLength *int
	mask             *string
}

// KeepAtStart sets the number of characters at the start that created obfuscators will skip when obfuscating.
// Defaults to 0.
//
// It panics if the given value is negative.
func (pb *PortionBuilder) KeepAtStart(value int) *PortionBuilder {
	validateNonNegative(value, "KeepAtStart")
	pb.keepAtStart = &value
	return pb
}

// KeepAtEnd sets the number of characters at the end that created obfuscators will skip when obfuscating.
// Defaults to 0.
//
// It panics if the given value is negative.
func (pb *PortionBuilder) KeepAtEnd(value int) *PortionBuilder {
	validateNonNegative(value, "KeepAtEnd")
	pb.keepAtEnd = &value
	return pb
}

// AtLeastFromStart sets the minimum number of characters from the start that need to be obfuscated.
// This will overrule any value for KeepAtStart or KeepAtEnd.
// Defaults to 0.
//
// It panics if the given value is negative.
func (pb *PortionBuilder) AtLeastFromStart(value int) *PortionBuilder {
	validateNonNegative(value, "AtLeastFromStart")
	pb.atLeastFromStart = &value
	return pb
}

// AtLeastFromEnd sets the minimum number of characters from the end that need to be obfuscated.
// This will overrule any value for KeepAtStart or KeepAtEnd.
// Defaults to 0.
//
// It panics if the given value is negative.
func (pb *PortionBuilder) AtLeastFromEnd(value int) *PortionBuilder {
	validateNonNegative(value, "AtLeastFromEnd")
	pb.atLeastFromEnd = &value
	return pb
}

// FixedTotalLength sets the fixed total length to use for obfuscated contents.
// When obfuscating, the result will have the mask added until this total length has been reached.
//
// Note: when used in combination with KeepAtStart and/or KeepAtEnd, this total length must be at least the sum of both other values.
// When used in combination with both, parts of the input may be repeated in the obfuscated content if the input's length is less than the combined
// number of characters to keep.
//
// It panics if the given value is negative.
func (pb *PortionBuilder) FixedTotalLength(value int) *PortionBuilder {
	validateNonNegative(value, "FixedTotalLength")
	pb.fixedTotalLength = &value
	return pb
}

// Mask sets the string to use for masking. Defaults to a single asterisk (*).
//
// It panics if the given mask is empty.
func (pb *PortionBuilder) Mask(mask string) *PortionBuilder {
	if len(mask) == 0 {
		log.Panicf("mask must not be empty")
	}
	pb.mask = &mask
	return pb
}

// Build returns an obfuscator that obfuscates a specific portion of their input, using the values configured in this builder.
//
// It panics if a FixedTotalLength is set that is smaller than the sum of the (default) values for KeepAtStart and KeepAtEnd.
func (pb *PortionBuilder) Build() Obfuscator {
	keepAtStart := valueOrDefault(pb.keepAtStart, 0)
	keepAtEnd := valueOrDefault(pb.keepAtEnd, 0)
	atLeastFromStart := valueOrDefault(pb.atLeastFromStart, 0)
	atLeastFromEnd := valueOrDefault(pb.atLeastFromEnd, 0)
	fixedTotalLength := pb.fixedTotalLength
	mask := valueOrDefault(pb.mask, "*")

	if fixedTotalLength != nil && *fixedTotalLength < keepAtStart+keepAtEnd {
		log.Panicf("FixedTotalLength (%d) < KeepAtStart (%d) + KeepAtEnd (%d)", *fixedTotalLength, keepAtStart, keepAtEnd)
	}

	allowDuplicates := pb.fixedTotalLength != nil

	return NewObfuscator(func(s string) string {
		originalLength := len(s)
		length := originalLength
		fromStart := calculateFromStart(length, keepAtStart, atLeastFromStart, atLeastFromEnd)
		fromEnd := calculateFromEnd(length, fromStart, keepAtEnd, atLeastFromStart, atLeastFromEnd, allowDuplicates)
		// 0 <= fromStart <= length == end - start, so start <= start + fromStart <= end
		// 0 <= fromEnd <= length == end - start, so 0 <= length - fromEnd and start <= end - fromEnd

		if fixedTotalLength != nil {
			length = *fixedTotalLength
		}
		obfuscatedLength := length - fromEnd - fromStart

		// Result: 0 to fromStart non-obfuscated, then obfuscated, then from end - fromEnd non-obfuscated
		return s[:fromStart] + strings.Repeat(mask, obfuscatedLength) + s[originalLength-fromEnd:]
	})
}

func validateNonNegative(value int, name string) {
	if value < 0 {
		log.Panicf("%s: %d < 0", name, value)
	}
}

func valueOrDefault[T interface{}](value *T, defaultValue T) T {
	if value == nil {
		return defaultValue
	}
	return *value
}

func calculateFromStart(length, keepAtStart, atLeastFromStart, atLeastFromEnd int) int {
	if atLeastFromStart > 0 {
		// The first characters need to be obfuscated so ignore keepAtStart
		return 0
	}
	// 0 <= keepAtMost <= length, the maximum number of characters to not obfuscate taking into account atLeastFromEnd
	// 0 <= result <= length, the minimum of what we want to obfuscate and what we can obfuscate
	keepAtMost := max(0, length-atLeastFromEnd)
	return min(keepAtStart, keepAtMost)
}

func calculateFromEnd(length, keepFromStart, keepAtEnd, atLeastFromStart, atLeastFromEnd int, allowDuplicates bool) int {
	if atLeastFromEnd > 0 {
		// The last characters need to be obfuscated so ignore keepAtEnd
		return 0
	}
	// 0 <= available <= length, the number of characters not already handled by fromStart (to prevent characters being appended twice)
	//                           if allowDuplicates then available == length
	// 0 <= keepAtMost <= length, the maximum number of characters to not obfuscate taking into account atLeastFromStart
	// 0 <= result <= length, the minimum of what we want to obfuscate and what we can obfuscate
	var available int
	if allowDuplicates {
		available = length
	} else {
		available = length - keepFromStart
	}
	keepAtMost := max(0, length-atLeastFromStart)
	return min(keepAtEnd, min(available, keepAtMost))
}
