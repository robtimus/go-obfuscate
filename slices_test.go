package obfuscate

import (
	"math"
	"testing"
)

func TestSliceWithStrings(t *testing.T) {
	obfuscator := WithFixedLength(3)

	input := []string{"foo", "bar", "hello", "world"}
	obfuscated := Slice(input, obfuscator)

	expected := []string{"***", "***", "***", "***"}
	if slicesDiffer(obfuscated, expected) {
		t.Errorf("expected: %v, actual: %v", expected, obfuscated)
	}
}

func TestSliceWithNumbers(t *testing.T) {
	obfuscator := WithFixedLength(3)

	input := []int{1, 2, 3, math.MinInt, math.MaxInt}
	obfuscated := Slice(input, obfuscator)

	expected := []string{"***", "***", "***", "***", "***"}
	if slicesDiffer(obfuscated, expected) {
		t.Errorf("expected: %v, actual: %v", expected, obfuscated)
	}
}
