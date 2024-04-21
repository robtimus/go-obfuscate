package obfuscate

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	obfuscator := All()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "***"},
		{"foobar", "******"},
		{"hello", "*****"},
		{"hello world", "***********"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestNone(t *testing.T) {
	obfuscator := None()
	inputs := []string{"foo", "foobar", "hello", "hello world"}
	for _, input := range inputs {
		expected := input
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestWithFixedLength(t *testing.T) {
	inputs := []string{"foo", "foobar", "hello", "hello world"}
	parameters := []struct {
		fixedLength int
		expected    string
	}{
		{0, ""},
		{1, "*"},
		{3, "***"},
		{8, "********"},
	}
	for i := range parameters {
		fixedLength := parameters[i].fixedLength
		obfuscator := WithFixedLength(fixedLength)
		expected := parameters[i].expected
		for _, input := range inputs {
			t.Run(fmt.Sprintf("WithFixedLength(%d) applied to '%s'", fixedLength, input), func(t *testing.T) {
				testObfuscateString(t, obfuscator, input, expected)
			})
		}
	}

	testPanic(t, "fixedLength < 0", func() {
		WithFixedLength(-1)
	}, "fixedLength: -1 < 0")
}

func TestWithFixedValue(t *testing.T) {
	inputs := []string{"foo", "foobar", "hello", "hello world"}
	fixedValues := []string{"", "obfuscated", "***"}
	for _, fixedValue := range fixedValues {
		obfuscator := WithFixedValue(fixedValue)
		for _, input := range inputs {
			expected := fixedValue
			t.Run(fmt.Sprintf("WithFixedValue('%s') applied to '%s'", fixedValue, input), func(t *testing.T) {
				testObfuscateString(t, obfuscator, input, expected)
			})
		}
	}

	testPanic(t, "fixedLength < 0", func() {
		WithFixedLength(-1)
	}, "fixedLength: -1 < 0")
}

func testObfuscateString(t *testing.T, obfuscator Obfuscator, input, expected string) {
	actual := obfuscator.ObfuscateString(input)
	if actual != expected {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func testPanic(t *testing.T, name string, action func(), expectedMessage string) {
	t.Run(name, func(t *testing.T) {
		defer func() {
			r := recover()
			if r == nil {
				t.Errorf("expected: %s, actual: nil", expectedMessage)
			} else if s, ok := r.(string); !ok || s != expectedMessage {
				t.Errorf("expected: %s, actual: %s", expectedMessage, r)
			}
		}()
		action()
		t.Errorf("expected an error")
	})
}
