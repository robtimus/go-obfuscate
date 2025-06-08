package obfuscate

import (
	"errors"
	"testing"
)

var errParse = errors.New("error")

func TestObfuscateStringOnErrorLog(t *testing.T) {
	logger := newCapturingLogger()
	obfuscator := NewObfuscatorOnErrorLog(func(s string) (string, error) {
		return s, errParse
	}, logger.Logger)

	testObfuscateStringWithErrors(t, obfuscator, logger, "foo", "foo", "ObfuscateString error: error\n")
}

func TestObfuscateStringOnErrorInclude(t *testing.T) {
	obfuscator := NewObfuscatorOnErrorInclude(func(s string) (string, error) {
		return s, errParse
	})

	testObfuscateStringWithErrors(t, obfuscator, nil, "foo", "foo<error: error>", "")
}

func TestObfuscateStringOnErrorDiscard(t *testing.T) {
	obfuscator := NewObfuscatorOnErrorDiscard(func(s string) (string, error) {
		return s, errParse
	})

	testObfuscateStringWithErrors(t, obfuscator, nil, "foo", "foo", "")
}

func testObfuscateStringWithErrors(t *testing.T, obfuscator Obfuscator, logger *CapturingLogger, input, expectedObfuscated, expectedLogged string) {
	t.Helper()

	testObfuscateString(t, obfuscator, input, expectedObfuscated)

	if logger != nil {
		assertEqual(t, expectedLogged, logger.String())
	}
}
