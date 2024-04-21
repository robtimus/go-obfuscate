package obfuscate

import "strings"

// SplitPoint represents a point in a string to split obfuscation.
// Like Obfuscator.UntilLength this can be used to combine obfuscators. For instance, to obfuscate email addresses:
//
//	localPartObfuscator := obfuscate.Portion()
//	    .KeepAtStart(1)
//	    .KeepAtEnd(1)
//	    .FixedTotalLength(8)
//	    .Build()
//	domainObfuscator := obfuscate.None()
//	obfuscator := obfuscate.AtFirst("@").SplitTo(localPrtObfuscator, domainObfuscator)
//	// Everything before @ will be obfuscated using localPartObfuscator, everything after @ will not be obfuscated.
//	// Example input: test@example.org
//	// Example output: t******t@example.org
//
// Unlike Obfuscator.UntilLength it's not possible to chain splitting, but it's of course possible to nest it:
//
//	localPartObfuscator := obfuscate.Portion()
//	    .KeepAtStart(1)
//	    .KeepAtEnd(1)
//	    .FixedTotalLength(8)
//	    .Build()
//	domainObfuscator := obfuscate.AtLast(".").SplitTo(obfuscate.All(), obfuscate.None())
//	obfuscator := obfuscate.AtFirst("@").SplitTo(localPartObfuscator, domainObfuscator)
//	// Everything before @ will be obfuscated using localPartObfuscator, everything after @ will be obfuscated until the last dot
//	// Example input: test@example.org
//	// Example output: t******t@*******.org
type SplitPoint struct {
	splitStart  func(s string) int
	splitLength int
}

// SplitTo creates an obfuscator that splits obfuscation at this split point.
// The part of the string before the split point will be obfuscated by the first obfuscator, the part after the split point by the second.
func (sp SplitPoint) SplitTo(beforeSplitPoint, afterSplitPoint Obfuscator) Obfuscator {
	return NewObfuscator(func(s string) string {
		splitStartIndex := sp.splitStart(s)
		if splitStartIndex < 0 {
			return beforeSplitPoint.ObfuscateString(s)
		}
		splitEndIndex := splitStartIndex + sp.splitLength
		return beforeSplitPoint.ObfuscateString(s[:splitStartIndex]) +
			s[splitStartIndex:splitEndIndex] +
			afterSplitPoint.ObfuscateString(s[splitEndIndex:])
	})
}

// NewSplitPoint creates a new split point.
//
// splitStart is a function that takes a string and returns the 0-based index where to split, or a negative value if obfuscation should not be split.
// This could for example be caused by a string to split on not being found.
//
// splitLength is the length of the split point. If not 0, the substring with this length starting at the calculated split start will not be obfuscated.
// This function panics if splitLength is negative.
func NewSplitPoint(splitStart func(s string) int, splitLength int) SplitPoint {
	validateNonNegative(splitLength, "splitLength")
	return SplitPoint{splitStart, splitLength}
}

// AtFirst creates a new split point that splits at the first occurrence of a string.
// This split point is exclusive; the string itself will not be obfuscated.
func AtFirst(s string) SplitPoint {
	return NewSplitPoint(func(input string) int {
		return strings.Index(input, s)
	}, len(s))
}

// AtLast creates a new split point that splits at the last occurrence of a string.
// This split point is exclusive; the string itself will not be obfuscated.
func AtLast(s string) SplitPoint {
	return NewSplitPoint(func(input string) int {
		return strings.LastIndex(input, s)
	}, len(s))
}

// AtNth creates a new split point that splits at a specific occurrence of a string.
// This split point is exclusive; the string itself will not be obfuscated.
//
// AtNth panics if the given zero-based occurrence is negative.
func AtNth(s string, occurrence int) SplitPoint {
	validateNonNegative(occurrence, "occurrence")
	return NewSplitPoint(func(input string) int {
		return nthIndex(input, s, occurrence)
	}, len(s))
}

func nthIndex(s, substr string, occurrence int) int {
	// Go doesn't have an index-lookup function with a start index so use string slicing instead
	source := s
	sourceStart := 0
	index := strings.Index(source, substr)
	for i := 1; i <= occurrence && index != -1; i++ {
		source = source[index+1:]
		sourceStart = sourceStart + index + 1
		index = strings.Index(source, substr)
	}
	if index == -1 {
		return index
	}
	return sourceStart + index
}
