package obfuscate

import (
	"testing"
)

func TestHTTPParameterDefaultErrorStrategy(t *testing.T) {
	obfuscator := newHTTPParameterObfuscatorBuilder().Build()

	assertEqual(t, OnErrorLog, obfuscator.onError)
}

func TestObfuscateParameterString(t *testing.T) {
	var obfuscator Obfuscator = newHTTPParameterObfuscatorBuilder().Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value"

	actual := obfuscator.ObfuscateString(input)

	expected := "foo=***&hello=world&FOO=BAR&empty=&no-value"

	assertEqual(t, expected, actual)

	actual, err := obfuscator.ParseAndObfuscateString(input)

	assertEqual(t, expected, actual)
	assertEqual(t, nil, err)
}

func TestObfuscateParameterStringWithError(t *testing.T) {
	logger := newCapturingLogger()

	obfuscator := newHTTPParameterObfuscatorBuilder().OnErrorLog(logger.Logger).Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value&err=%A&err=%B"

	obfuscated, err := obfuscator.ParseAndObfuscateString(input)

	assertEqual(t, "foo=***&hello=world&FOO=BAR&empty=&no-value&err=", obfuscated)
	assertEqual(t, "invalid URL escape \"%A\"", err.Error())
	assertEqual(t, "", logger.String())
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

func TestObfuscateParameterStringOnErrorDiscard(t *testing.T) {
	builder := newHTTPParameterObfuscatorBuilder().OnErrorDiscard()

	testObfuscateParameterStringWithErrors(t, builder, nil,
		"foo=***&hello=world&FOO=BAR&empty=&no-value&err=",
		"")
}

func testObfuscateParameterStringWithErrors(t *testing.T, builder *HTTPParameterObfuscatorBuilder, logger *CapturingLogger, expectedObfuscated, expectedLogged string) {
	t.Helper()

	obfuscator := builder.Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-value&err=%A&err=%B"

	testObfuscateStringWithErrors(t, obfuscator, logger, input, expectedObfuscated, expectedLogged)
}

func TestHTTPParameterObfuscatorChaining(t *testing.T) {
	obfuscator := WithFixedLength(3)
	parameterObfuscator := newHTTPParameterObfuscatorBuilder().Build()

	input := "foo=bar&hello=world&FOO=BAR&empty=&no-valuepostfix"

	obfuscator = parameterObfuscator.UntilLength(len(input) - 7).Then(obfuscator)

	testParseAndObfuscateString(t, obfuscator, input, "foo=***&hello=world&FOO=BAR&empty=&no-value***")
}

func newHTTPParameterObfuscatorBuilder() *HTTPParameterObfuscatorBuilder {
	return HTTPParameters().
		WithParameter("foo", All()).
		WithParameter("no-value", All()).
		WithParameter("err", All())
}
