package obfuscate

import (
	"fmt"
	"log"
	"strings"
	"testing"
)

func TestHTTPParameterDefaultErrorStrategy(t *testing.T) {
	obfuscator := newHTTPParameterObfuscator(nil)

	actual := obfuscator.onError

	expected := OnErrorLog

	if actual != expected {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}
}

func TestHTTPParameterDefaultLogging(t *testing.T) {
	obfuscator := newHTTPParameterObfuscator(nil)

	testPanic(t, "log.Panicf", func() {
		obfuscator.panicf("panic: %v", obfuscator)
	}, fmt.Sprintf("panic: %v", obfuscator))
}

func TestObfuscateParameterString(t *testing.T) {
	obfuscator := newHTTPParameterObfuscator(nil)

	input := "foo=bar&hello=world&empty=&no-value"

	actual := obfuscator.ObfuscateString(input)

	expected := "foo=***&hello=world&empty=&no-value"

	if actual != expected {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}

	actual, err := obfuscator.ObfuscateParameterString(input)

	if actual != expected {
		t.Errorf("expected: '%v', actual: '%v'", expected, actual)
	}

	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}
}

func TestObfuscateParameterStringWithError(t *testing.T) {
	output := &strings.Builder{}
	logger := log.New(output, "", 0)

	obfuscator := newHTTPParameterObfuscator(&HTTPParameterObfuscatorOptions{OnError: OnErrorPanic, Logger: logger})

	input := "foo=bar&hello=world&empty=&no-value&err=%A&err=%B"

	actualOutput, actualErr := obfuscator.ObfuscateParameterString(input)

	expectedOutput := "foo=***&hello=world&empty=&no-value&err="

	if actualOutput != expectedOutput {
		t.Errorf("expected: '%v', actual: '%v'", expectedOutput, actualOutput)
	}

	expectedError := "invalid URL escape \"%A\""

	if actualErr.Error() != expectedError {
		t.Errorf("expected: '%v', actual: '%v'", expectedError, actualErr)
	}

	actualLogged := output.String()

	expectedLogged := ""

	if actualLogged != expectedLogged {
		t.Errorf("expected: '%v', actual: '%v'", expectedLogged, actualLogged)
	}
}

func TestObfuscateParameterStringOnErrorLog(t *testing.T) {
	testObfuscateParameterStringWithErrors(t, OnErrorLog,
		"foo=***&hello=world&empty=&no-value&err=",
		"ObfuscateString error: invalid URL escape \"%A\"\n")
}

func TestObfuscateParameterStringOnErrorInclude(t *testing.T) {
	testObfuscateParameterStringWithErrors(t, OnErrorInclude,
		"foo=***&hello=world&empty=&no-value&err=<error: invalid URL escape \"%A\">",
		"")
}

func TestObfuscateParameterStringOnErrorStop(t *testing.T) {
	testObfuscateParameterStringWithErrors(t, OnErrorStop,
		"foo=***&hello=world&empty=&no-value&err=",
		"")
}

func TestObfuscateParameterStringOnErrorPanic(t *testing.T) {
	output := &strings.Builder{}
	logger := log.New(output, "", 0)

	var obfuscator Obfuscator = newHTTPParameterObfuscator(&HTTPParameterObfuscatorOptions{OnError: OnErrorPanic, Logger: logger})

	input := "foo=bar&hello=world&empty=&no-value&err=%A&err=%B"

	testPanic(t, "TestObfuscateParameterStringOnErrorPanic", func() {
		obfuscator.ObfuscateString(input)
	}, "ObfuscateString error: invalid URL escape \"%A\"")

	actualLogged := output.String()

	expectedLogged := "ObfuscateString error: invalid URL escape \"%A\"\n"

	if actualLogged != expectedLogged {
		t.Errorf("expected: '%v', actual: '%v'", expectedLogged, actualLogged)
	}
}

func testObfuscateParameterStringWithErrors(t *testing.T, onError ErrorStrategy, expectedOutput, expectedLogged string) {
	output := &strings.Builder{}
	logger := log.New(output, "", 0)

	var obfuscator Obfuscator = newHTTPParameterObfuscator(&HTTPParameterObfuscatorOptions{OnError: onError, Logger: logger})

	input := "foo=bar&hello=world&empty=&no-value&err=%A&err=%B"

	actualOutput := obfuscator.ObfuscateString(input)

	if actualOutput != expectedOutput {
		t.Errorf("expected: '%v', actual: '%v'", expectedOutput, actualOutput)
	}

	actualLogged := output.String()

	if actualLogged != expectedLogged {
		t.Errorf("expected: '%v', actual: '%v'", expectedLogged, actualLogged)
	}
}

func newHTTPParameterObfuscator(options *HTTPParameterObfuscatorOptions) HTTPParameterObfuscator {
	obfuscators := map[string]Obfuscator{
		"foo":      All(),
		"no-value": All(),
		"err":      All(),
	}
	return HTTPParameters(obfuscators, options)
}
