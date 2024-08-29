package obfuscate

import "testing"

func TestDefaultErrorStrategy(t *testing.T) {
	var actual ErrorStrategy

	expected := OnErrorLog

	assertEqual(t, expected, actual)
}

func TestErrorStrategyOnErrorLog(t *testing.T) {
	testErrorStrategy(t, OnErrorLog, 0, "OnErrorLog")
}

func TestErrorStrategyOnErrorInclude(t *testing.T) {
	testErrorStrategy(t, OnErrorInclude, 1, "OnErrorInclude")
}

func TestErrorStrategyOnErrorStop(t *testing.T) {
	testErrorStrategy(t, OnErrorStop, 2, "OnErrorStop")
}

func TestErrorStrategyOnErrorPanic(t *testing.T) {
	testErrorStrategy(t, OnErrorPanic, 3, "OnErrorPanic")
}

func TestErrorStrategyUnknown(t *testing.T) {
	testErrorStrategy(t, ErrorStrategy(255), 255, "ErrorStrategy(255)")
}

func testErrorStrategy(t *testing.T, onError ErrorStrategy, expectedValue int, expectedString string) {
	actualValue := int(onError)

	if actualValue != expectedValue {
		t.Errorf("expected: '%v', actual: '%v'", expectedValue, actualValue)
	}

	actualString := onError.String()

	if actualString != expectedString {
		t.Errorf("expected: '%v', actual: '%v'", expectedString, actualString)
	}
}
