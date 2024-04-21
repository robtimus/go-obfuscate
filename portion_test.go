package obfuscate

import (
	"fmt"
	"math"
	"testing"
)

func TestKeepAtStart(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo"},
		{"foobar", "foob**"},
		{"hello", "hell*"},
		{"hello world", "hell*******"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEnd(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo"},
		{"foobar", "**obar"},
		{"hello", "*ello"},
		{"hello world", "*******orld"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndKeepAtEnd(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).KeepAtEnd(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo"},
		{"foobar", "foobar"},
		{"hello", "hello"},
		{"hello world", "hell***orld"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).KeepAtEnd(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndAtLeastFromEnd(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).AtLeastFromEnd(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "***"},
		{"foobar", "fo****"},
		{"hello", "h****"},
		{"hello world", "hell*******"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).AtLeastFromEnd(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEndAndAtLeastFromStart(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).AtLeastFromStart(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "***"},
		{"foobar", "****ar"},
		{"hello", "****o"},
		{"hello world", "*******orld"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4).AtLeastFromStart(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).FixedTotalLength(9).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo******"},
		{"foobar", "foob*****"},
		{"hello", "hell*****"},
		{"hello world", "hell*****"},
		{"", "*********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).FixedTotalLength(9) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEndAndFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).FixedTotalLength(9).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "******foo"},
		{"foobar", "*****obar"},
		{"hello", "*****ello"},
		{"hello world", "*****orld"},
		{"", "*********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4).FixedTotalLength(9) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndKeepAtEndAndFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).KeepAtEnd(4).FixedTotalLength(9).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo***foo"},
		{"foobar", "foob*obar"},
		{"hello", "hell*ello"},
		{"hello world", "hell*orld"},
		{"", "*********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).KeepAtEnd(4).FixedTotalLength(9) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndAtLeastFromEndAndFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).AtLeastFromEnd(4).FixedTotalLength(9).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "*********"},
		{"foobar", "fo*******"},
		{"hello", "h********"},
		{"hello world", "hell*****"},
		{"", "*********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).AtLeastFromEnd(4).FixedTotalLength(9) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEndAndAtLeastFromStartAndFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).AtLeastFromStart(4).FixedTotalLength(9).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "*********"},
		{"foobar", "*******ar"},
		{"hello", "********o"},
		{"hello world", "*****orld"},
		{"", "*********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4).AtLeastFromStart(4).FixedTotalLength(9) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndEqualFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).FixedTotalLength(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo*"},
		{"foobar", "foob"},
		{"hello", "hell"},
		{"hello world", "hell"},
		{"", "****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).FixedTotalLength(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEndAndEqualFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).FixedTotalLength(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "*foo"},
		{"foobar", "obar"},
		{"hello", "ello"},
		{"hello world", "orld"},
		{"", "****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4).FixedTotalLength(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndKeepAtEndAndEqualFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).KeepAtEnd(4).FixedTotalLength(8).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "foo**foo"},
		{"foobar", "foobobar"},
		{"hello", "hellello"},
		{"hello world", "hellorld"},
		{"", "********"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).KeepAtEnd(4).FixedTotalLength(8) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtStartAndAtLeastFromEndAndEqualFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtStart(4).AtLeastFromEnd(4).FixedTotalLength(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "****"},
		{"foobar", "fo**"},
		{"hello", "h***"},
		{"hello world", "hell"},
		{"", "****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(4).AtLeastFromEnd(4).FixedTotalLength(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestKeepAtEndAndAtLeastFromStartAndEqualFixedTotalLength(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(4).AtLeastFromStart(4).FixedTotalLength(4).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "****"},
		{"foobar", "**ar"},
		{"hello", "***o"},
		{"hello world", "orld"},
		{"", "****"},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(4).AtLeastFromStart(4).FixedTotalLength(4) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestLastTwoCharactersOnly(t *testing.T) {
	obfuscator := Portion().KeepAtStart(math.MaxInt).AtLeastFromEnd(2).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "f**"},
		{"foobar", "foob**"},
		{"hello", "hel**"},
		{"hello world", "hello wor**"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtStart(MAX).AtLeastFromEnd(2) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestFirstTwoCharactersOnly(t *testing.T) {
	obfuscator := Portion().KeepAtEnd(math.MaxInt).AtLeastFromStart(2).Build()
	parameters := []struct {
		input, expected string
	}{
		{"foo", "**o"},
		{"foobar", "**obar"},
		{"hello", "**llo"},
		{"hello world", "**llo world"},
		{"", ""},
	}
	for i := range parameters {
		input := parameters[i].input
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("KeepAtEnd(MAX).AtLeastFromStart(2) applied to '%s'", input), func(t *testing.T) {
			testObfuscateString(t, obfuscator, input, expected)
		})
	}
}

func TestInvalidInput(t *testing.T) {
	testPanic(t, "KeepAtStart < 0", func() {
		Portion().KeepAtStart(-1)
	}, "KeepAtStart: -1 < 0")

	testPanic(t, "KeepAtEnd < 0", func() {
		Portion().KeepAtEnd(-1)
	}, "KeepAtEnd: -1 < 0")

	testPanic(t, "AtLeastFromStart < 0", func() {
		Portion().AtLeastFromStart(-1)
	}, "AtLeastFromStart: -1 < 0")

	testPanic(t, "AtLeastFromEnd < 0", func() {
		Portion().AtLeastFromEnd(-1)
	}, "AtLeastFromEnd: -1 < 0")

	testPanic(t, "KeepAtStart > FixedTotalLength", func() {
		Portion().KeepAtStart(4).FixedTotalLength(3).Build()
	}, "FixedTotalLength (3) < KeepAtStart (4) + KeepAtEnd (0)")

	testPanic(t, "KeepAtEnd > FixedTotalLength", func() {
		Portion().KeepAtEnd(4).FixedTotalLength(3).Build()
	}, "FixedTotalLength (3) < KeepAtStart (0) + KeepAtEnd (4)")

	testPanic(t, "KeepAtStart + KeepAtEnd > FixedTotalLength", func() {
		Portion().KeepAtStart(4).KeepAtEnd(4).FixedTotalLength(7).Build()
	}, "FixedTotalLength (7) < KeepAtStart (4) + KeepAtEnd (4)")
}
