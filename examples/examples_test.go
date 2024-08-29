package examples

import (
	"math"
	"strings"
	"testing"

	"github.com/robtimus/go-obfuscate"
)

func TestObfuscateAllExample(t *testing.T) {
	obfuscator := obfuscate.All()
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "***********", obfuscated)
}

func TestObfuscateWithFixedLengthExample(t *testing.T) {
	obfuscator := obfuscate.WithFixedLength(5)
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "*****", obfuscated)
}

func TestObfuscateWithFixedValueExample(t *testing.T) {
	obfuscator := obfuscate.WithFixedValue("foo")
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "foo", obfuscated)
}

func TestObfuscatePortionAllButLast4Example(t *testing.T) {
	obfuscator := obfuscate.Portion().KeepAtEnd(4).Build()
	obfuscated := obfuscator.ObfuscateString("1234567890123456")
	assertEqual(t, "************3456", obfuscated)
}

func TestObfuscatePortionAllButLast4AtLeast12FromStartExample(t *testing.T) {
	obfuscator := obfuscate.Portion().KeepAtEnd(4).AtLeastFromStart(12).Build()
	obfuscated := obfuscator.ObfuscateString("1234567890")
	assertEqual(t, "**********", obfuscated)
}

func TestObfuscatePortionLast2Example(t *testing.T) {
	obfuscator := obfuscate.Portion().KeepAtStart(math.MaxInt).AtLeastFromEnd(2).Build()
	obfuscated := obfuscator.ObfuscateString("SW1A 2AA")
	assertEqual(t, "SW1A 2**", obfuscated)
}

func TestObfuscatePortionWithFixedLengthExample(t *testing.T) {
	obfuscator := obfuscate.Portion().KeepAtStart(2).KeepAtEnd(2).FixedTotalLength(6).Build()
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "He**ld", obfuscated)

	obfuscated = obfuscator.ObfuscateString("foo")
	assertEqual(t, "fo**oo", obfuscated)
}

func TestObfuscateUpperCaseExample(t *testing.T) {
	obfuscator := obfuscate.NewObfuscator(func(text string) string { return strings.ToUpper(text) })
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "HELLO WORLD", obfuscated)
}

func TestObfuscateNoneExample(t *testing.T) {
	somePossiblyNilObfuscator := getNilObfuscator()
	obfuscator := somePossiblyNilObfuscator
	if obfuscator == nil {
		obfuscator = obfuscate.None()
	}
	obfuscated := obfuscator.ObfuscateString("Hello World")
	assertEqual(t, "Hello World", obfuscated)
}

func getNilObfuscator() obfuscate.Obfuscator {
	return nil
}

func TestCombiningObfuscatorsExample(t *testing.T) {
	obfuscator := obfuscate.Portion().KeepAtStart(4).KeepAtEnd(4).Build()
	obfuscated := obfuscator.ObfuscateString("1234567890123456")
	assertEqual(t, "1234********3456", obfuscated)

	incorrectlyObfuscated := obfuscator.ObfuscateString("12345678901234")
	assertEqual(t, "1234******1234", incorrectlyObfuscated)
}

func TestUntilLengthExample(t *testing.T) {
	obfuscator := obfuscate.None().UntilLength(4).Then(obfuscate.All()).UntilLength(12).Then(obfuscate.None())
	obfuscated := obfuscator.ObfuscateString("1234567890123456")
	assertEqual(t, "1234********3456", obfuscated)
}

func TestUntilLengthWithPortionExample(t *testing.T) {
	obfuscator := obfuscate.None().UntilLength(4).Then(obfuscate.Portion().KeepAtEnd(4).AtLeastFromStart(8).Build())
	obfuscated := obfuscator.ObfuscateString("12345678901234")
	assertEqual(t, "1234********34", obfuscated)
}

func TestSplitPointExample(t *testing.T) {
	// Keep the domain as-is
	localPartObfuscator := obfuscate.Portion().KeepAtStart(1).KeepAtEnd(1).FixedTotalLength(8).Build()
	domainObfuscator := obfuscate.None()
	obfuscator := obfuscate.AtFirst("@").SplitTo(localPartObfuscator, domainObfuscator)
	obfuscated := obfuscator.ObfuscateString("test@example.org")
	assertEqual(t, "t******t@example.org", obfuscated)
}

func TestNestedSplitPointExample(t *testing.T) {
	// Keep only the TLD of the domain
	localPartObfuscator := obfuscate.Portion().KeepAtStart(1).KeepAtEnd(1).FixedTotalLength(8).Build()
	domainObfuscator := obfuscate.AtLast(".").SplitTo(obfuscate.All(), obfuscate.None())
	obfuscator := obfuscate.AtFirst("@").SplitTo(localPartObfuscator, domainObfuscator)
	obfuscated := obfuscator.ObfuscateString("test@example.org")
	assertEqual(t, "t******t@*******.org", obfuscated)
}

func TestHTTPHeadersExample(t *testing.T) {
	headerObfuscator := obfuscate.HTTPHeaders(map[string]obfuscate.Obfuscator{
		"Authorization": obfuscate.WithFixedLength(3),
	})
	obfuscatedAuthorization := headerObfuscator.ObfuscateHeaderValue("authorization", "Bearer someToken")
	assertEqual(t, "***", obfuscatedAuthorization)
	obfuscatedAuthorizations := headerObfuscator.ObfuscateHeaderValues("authorization", []string{"Bearer someToken"})
	assertEqualSlices(t, []string{"***"}, obfuscatedAuthorizations)
	obfuscatedContentType := headerObfuscator.ObfuscateHeaderValue("Content-Type", "application/json")
	assertEqual(t, "application/json", obfuscatedContentType)
	obfuscatedHeaders := headerObfuscator.ObfuscateHeaderMap(map[string]string{
		"authorization": "Bearer someToken",
		"content-type":  "application/json",
	})
	assertEqualMaps(t, map[string]string{"authorization": "***", "content-type": "application/json"}, obfuscatedHeaders)
}

func TestHTTPParamsExample(t *testing.T) {
	paramsObfuscator := obfuscate.HTTPParameters(map[string]obfuscate.Obfuscator{
		"password": obfuscate.WithFixedLength(3),
	}, nil)
	obfuscatedPassword := paramsObfuscator.ObfuscateParameter("password", "admin1234")
	assertEqual(t, "***", obfuscatedPassword)
	obfuscatedUsername := paramsObfuscator.ObfuscateParameter("username", "admin")
	assertEqual(t, "admin", obfuscatedUsername)
	obfuscatedParamString, err := paramsObfuscator.ObfuscateParameterString("username=admin&password=admin1234")
	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}
	assertEqual(t, "username=admin&password=***", obfuscatedParamString)
}

func TestHTTPParamsWithHandlerExample(t *testing.T) {
	paramsObfuscator := obfuscate.HTTPParameters(map[string]obfuscate.Obfuscator{
		"password": obfuscate.WithFixedLength(3),
	}, &obfuscate.HTTPParameterObfuscatorOptions{OnError: obfuscate.OnErrorInclude})
	obfuscatedParamString := paramsObfuscator.ObfuscateString("username=admin&password=admin1234%A")
	assertEqual(t, "username=admin&password=<error: invalid URL escape \"%A\">", obfuscatedParamString)
}

func TestMapsExample(t *testing.T) {
	mapObfuscator := obfuscate.Maps(map[string]obfuscate.Obfuscator{
		"password": obfuscate.WithFixedLength(3),
	})
	obfuscatedMap := mapObfuscator.ObfuscateMap(map[string]string{
		"username": "admin",
		"password": "admin1234",
	})
	assertEqualMaps(t, map[string]string{
		"username": "admin",
		"password": "***",
	}, obfuscatedMap)
}

func assertEqual[T comparable](t *testing.T, expected, actual T) {
	if actual != expected {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}
}

func assertEqualSlices[T comparable](t *testing.T, expected, actual []T) {
	if len(expected) != len(actual) {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}
	for i, expectedValue := range expected {
		if expectedValue != actual[i] {
			t.Errorf("expected: '%v', actual: '%v'", expected, actual)
		}
	}
}

func assertEqualMaps[K comparable, V comparable](t *testing.T, expected, actual map[K]V) {
	if len(expected) != len(actual) {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}
	for key, expectedValue := range expected {
		if actualValue, ok := actual[key]; !ok || expectedValue != actualValue {
			t.Errorf("expected: '%v', actual: '%v'", expected, actual)
		}
	}
}
