package obfuscate

import (
	"fmt"
	"testing"
)

func TestNoneUntilLength4ThenAllUntilLength12ThenNone(t *testing.T) {
	obfuscator := None().UntilLength(4).Then(All()).UntilLength(12).Then(None())
	parameters := []struct {
		input, expected string
	}{
		{"0", "0"},
		{"01", "01"},
		{"012", "012"},
		{"0123", "0123"},
		{"01234", "0123*"},
		{"012345", "0123**"},
		{"0123456", "0123***"},
		{"01234567", "0123****"},
		{"012345678", "0123*****"},
		{"0123456789", "0123******"},
		{"0123456789A", "0123*******"},
		{"0123456789AB", "0123********"},
		{"0123456789ABC", "0123********C"},
		{"0123456789ABCD", "0123********CD"},
		{"0123456789ABCDE", "0123********CDE"},
		{"0123456789ABCDEF", "0123********CDEF"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestNoneUntilLength4ThenWithFixedLength3(t *testing.T) {
	obfuscator := None().UntilLength(4).Then(WithFixedLength(3))
	parameters := []struct {
		input, expected string
	}{
		{"0", "0"},
		{"01", "01"},
		{"012", "012"},
		{"0123", "0123"},
		{"01234", "0123***"},
		{"012345", "0123***"},
		{"0123456", "0123***"},
		{"01234567", "0123***"},
		{"012345678", "0123***"},
		{"0123456789", "0123***"},
		{"0123456789A", "0123***"},
		{"0123456789AB", "0123***"},
		{"0123456789ABC", "0123***"},
		{"0123456789ABCD", "0123***"},
		{"0123456789ABCDE", "0123***"},
		{"0123456789ABCDEF", "0123***"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestWithFixedLength3UntilLength4ThenNone(t *testing.T) {
	obfuscator := WithFixedLength(3).UntilLength(4).Then(None())
	parameters := []struct {
		input, expected string
	}{
		{"0", "***"},
		{"01", "***"},
		{"012", "***"},
		{"0123", "***"},
		{"01234", "***4"},
		{"012345", "***45"},
		{"0123456", "***456"},
		{"01234567", "***4567"},
		{"012345678", "***45678"},
		{"0123456789", "***456789"},
		{"0123456789A", "***456789A"},
		{"0123456789AB", "***456789AB"},
		{"0123456789ABC", "***456789ABC"},
		{"0123456789ABCD", "***456789ABCD"},
		{"0123456789ABCDE", "***456789ABCDE"},
		{"0123456789ABCDEF", "***456789ABCDEF"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestWithFixedLength3UntilLength4ThenWithFixedValueXxx(t *testing.T) {
	obfuscator := WithFixedLength(3).UntilLength(4).Then(WithFixedValue("xxx"))
	parameters := []struct {
		input, expected string
	}{
		{"0", "***"},
		{"01", "***"},
		{"012", "***"},
		{"0123", "***"},
		{"01234", "***xxx"},
		{"012345", "***xxx"},
		{"0123456", "***xxx"},
		{"01234567", "***xxx"},
		{"012345678", "***xxx"},
		{"0123456789", "***xxx"},
		{"0123456789A", "***xxx"},
		{"0123456789AB", "***xxx"},
		{"0123456789ABC", "***xxx"},
		{"0123456789ABCD", "***xxx"},
		{"0123456789ABCDE", "***xxx"},
		{"0123456789ABCDEF", "***xxx"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestInvalidInputLengths(t *testing.T) {
	obfuscator := None()
	testPanic(t, "First prefix length", func() {
		obfuscator.UntilLength(0)
	}, "prefixLength: 0 < 1")

	obfuscator = None().UntilLength(1).Then(All())
	testPanic(t, "Second prefix length", func() {
		obfuscator.UntilLength(1)
	}, "prefixLength: 1 < 2")

	obfuscator = None().UntilLength(1).Then(All()).UntilLength(2).Then(None())
	testPanic(t, "Third prefix length", func() {
		obfuscator.UntilLength(2)
	}, "prefixLength: 2 < 3")

	obfuscator = None().UntilLength(1).Then(All()).UntilLength(2).Then(None()).UntilLength(3).Then(All())
	testPanic(t, "Fourth prefix length", func() {
		obfuscator.UntilLength(3)
	}, "prefixLength: 3 < 4")
}
