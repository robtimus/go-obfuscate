package obfuscate

import (
	"log"
	"strings"
	"testing"
)

func TestHTTPParameterDefaultErrorStrategy(t *testing.T) {
	obfuscator := newHTTPParameterObfuscatorBuilder().Build()

	actual := obfuscator.onError

	expected := OnErrorLog

	assertEqual(t, expected, actual)
}

func TestObfuscateParameterString(t *testing.T) {
	var obfuscator ParsingObfuscator = newHTTPParameterObfuscatorBuilder().Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value"

	actual := obfuscator.ObfuscateString(input)

	expected := "foo=***&hello=world&FOO=BAR&empty=&no-value"

	assertEqual(t, expected, actual)

	actual, err := obfuscator.ParseAndObfuscateString(input)

	assertEqual(t, expected, actual)

	if err != nil {
		t.Errorf("unexpected error: '%v'", err)
	}
}

func TestObfuscateParameterStringWithError(t *testing.T) {
	output := &strings.Builder{}
	logger := log.New(output, "", 0)

	obfuscator := newHTTPParameterObfuscatorBuilder().OnErrorLog(logger).Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value&err=%A&err=%B"

	actualOutput, actualErr := obfuscator.ParseAndObfuscateString(input)

	expectedOutput := "foo=***&hello=world&FOO=BAR&empty=&no-value&err="

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
	logger := newCapturingLogger()
	builder := newHTTPParameterObfuscatorBuilder().OnErrorLog(logger.Logger)

	testObfuscateParameterStringWithErrors(t, builder, logger,
		"foo=***&hello=world&FOO=BAR&empty=&no-value&err=",
		"ObfuscateString error: invalid URL escape \"%A\"\n")
}

func TestObfuscateParameterStringOnErrorInclude(t *testing.T) {
	builder := newHTTPParameterObfuscatorBuilder().OnErrorInclude()

	testObfuscateParameterStringWithErrors(t, builder, nil,
		"foo=***&hello=world&FOO=BAR&empty=&no-value&err=<error: invalid URL escape \"%A\">",
		"")
}

func TestObfuscateParameterStringOnErrorStop(t *testing.T) {
	builder := newHTTPParameterObfuscatorBuilder().OnErrorStop()

	testObfuscateParameterStringWithErrors(t, builder, nil,
		"foo=***&hello=world&FOO=BAR&empty=&no-value&err=",
		"")
}

func testObfuscateParameterStringWithErrors(t *testing.T, builder *HTTPParameterObfuscatorBuilder, logger *CapturingLogger, expectedOutput, expectedLogged string) {
	t.Helper()

	var obfuscator Obfuscator = builder.Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value&err=%A&err=%B"

	actualOutput := obfuscator.ObfuscateString(input)

	if actualOutput != expectedOutput {
		t.Errorf("expected: '%v', actual: '%v'", expectedOutput, actualOutput)
	}

	if logger != nil {
		actualLogged := logger.String()

		if actualLogged != expectedLogged {
			t.Errorf("expected: '%v', actual: '%v'", expectedLogged, actualLogged)
		}
	}
}

func newHTTPParameterObfuscatorBuilder() *HTTPParameterObfuscatorBuilder {
	return HTTPParameters().
		WithParameter("foo", All()).
		WithParameter("no-value", All()).
		WithParameter("err", All())
}
