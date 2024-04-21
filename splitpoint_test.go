package obfuscate

import (
	"fmt"
	"strings"
	"testing"
)

func TestAtFirst(t *testing.T) {
	obfuscator := AtFirst("@").SplitTo(All(), WithFixedLength(5))
	parameters := []struct {
		input, expected string
	}{
		{"test", "****"},
		{"test@", "****@*****"},
		{"test@example.org", "****@*****"},
		{"test@example.org@", "****@*****"},
		{"test@example.org@localhost", "****@*****"},
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

func TestAtLast(t *testing.T) {
	obfuscator := AtLast("@").SplitTo(All(), WithFixedLength(5))
	parameters := []struct {
		input, expected string
	}{
		{"test", "****"},
		{"test@", "****@*****"},
		{"test@example.org", "****@*****"},
		{"test@example.org@", "****************@*****"},
		{"test@example.org@localhost", "****************@*****"},
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

func TestAtNthNegativeOccurrence(t *testing.T) {
	testPanic(t, "Negative occurrence", func() {
		AtNth(".", -1)
	}, "occurrence: -1 < 0")
}

func TestAtNthWithOccurrence0(t *testing.T) {
	obfuscator := AtNth(".", 0).SplitTo(WithFixedValue("xxx"), WithFixedLength(5))
	parameters := []struct {
		input, expected string
	}{
		{"alpha", "xxx"},
		{"alpha.", "xxx.*****"},
		{"alpha.bravo", "xxx.*****"},
		{"alpha.bravo.charlie", "xxx.*****"},
		{"alpha.bravo.charlie.", "xxx.*****"},
		{"alpha.bravo.charlie.delta", "xxx.*****"},
		{"alpha.bravo.charlie.delta.echo", "xxx.*****"},
		{"........", "xxx.*****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestAtNthWithOccurrence1(t *testing.T) {
	obfuscator := AtNth(".", 1).SplitTo(WithFixedValue("xxx"), WithFixedLength(5))
	parameters := []struct {
		input, expected string
	}{
		{"alpha", "xxx"},
		{"alpha.bravo", "xxx"},
		{"alpha.bravo.", "xxx.*****"},
		{"alpha.bravo.charlie", "xxx.*****"},
		{"alpha.bravo.charlie.", "xxx.*****"},
		{"alpha.bravo.charlie.delta", "xxx.*****"},
		{"alpha.bravo.charlie.delta.echo", "xxx.*****"},
		{"........", "xxx.*****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestAtNthWithOccurrence2(t *testing.T) {
	obfuscator := AtNth(".", 2).SplitTo(WithFixedValue("xxx"), WithFixedLength(5))
	parameters := []struct {
		input, expected string
	}{
		{"alpha", "xxx"},
		{"alpha.bravo", "xxx"},
		{"alpha.bravo.charlie", "xxx"},
		{"alpha.bravo.charlie.", "xxx.*****"},
		{"alpha.bravo.charlie.delta", "xxx.*****"},
		{"alpha.bravo.charlie.delta.echo", "xxx.*****"},
		{"........", "xxx.*****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestNewSplitPointWithNegativeSplitLength(t *testing.T) {
	testPanic(t, "Negative splitLength", func() {
		NewSplitPoint(func(s string) int {
			return -1
		}, -1)
	}, "splitLength: -1 < 0")
}

func TestNewSplitPointWithSplitLength0(t *testing.T) {
	splitPoint := NewSplitPoint(func(s string) int {
		return strings.Index(s, "@")
	}, 0)
	obfuscator := splitPoint.SplitTo(All(), WithFixedLengthWithMask(5, "x"))
	parameters := []struct {
		input, expected string
	}{
		{"test", "****"},
		{"test@", "****xxxxx"},
		{"test@example.org", "****xxxxx"},
		{"test@example.org@", "****xxxxx"},
		{"test@example.org@localhost", "****xxxxx"},
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

func TestNewSplitPointWithSplitAtStart(t *testing.T) {
	splitPoint := NewSplitPoint(func(s string) int {
		return 0
	}, 0)
	obfuscator := splitPoint.SplitTo(WithFixedLength(3), AllWithMask("x"))
	parameters := []struct {
		input, expected string
	}{
		{"test", "***xxxx"},
		{"test@", "***xxxxx"},
		{"test@example.org", "***xxxxxxxxxxxxxxxx"},
		{"test@example.org@", "***xxxxxxxxxxxxxxxxx"},
		{"test@example.org@localhost", "***xxxxxxxxxxxxxxxxxxxxxxxxxx"},
		{"", "***"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestNewSplitPointWithSplitAtEnd(t *testing.T) {
	splitPoint := NewSplitPoint(func(s string) int {
		return len(s)
	}, 0)
	obfuscator := splitPoint.SplitTo(WithFixedLength(3), WithFixedValue("xxx"))
	parameters := []struct {
		input, expected string
	}{
		{"test", "***xxx"},
		{"test@", "***xxx"},
		{"test@example.org", "***xxx"},
		{"test@example.org@", "***xxx"},
		{"test@example.org@localhost", "***xxx"},
		{"", "***xxx"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}
