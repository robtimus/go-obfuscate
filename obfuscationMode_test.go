package obfuscate

import "testing"

func TestDefaultObfuscationMode(t *testing.T) {
	var actual ObfuscationMode

	expected := Exclude

	assertEqual(t, expected, actual)
}

func TestObfuscationModeExclude(t *testing.T) {
	testObfuscationMode(t, Exclude, 0, "Exclude")
}

func TestObfuscationModeExcludeAll(t *testing.T) {
	testObfuscationMode(t, ExcludeAll, 1, "ExcludeAll")
}

func TestObfuscationModeInherit(t *testing.T) {
	testObfuscationMode(t, Inherit, 2, "Inherit")
}

func TestObfuscationModeInheritOverridable(t *testing.T) {
	testObfuscationMode(t, InheritOverridable, 3, "InheritOverridable")
}

func TestObfuscationModeUnknown(t *testing.T) {
	testObfuscationMode(t, ObfuscationMode(255), 255, "ObfuscationMode(255)")
}

func testObfuscationMode(t *testing.T, obfuscationMode ObfuscationMode, expectedValue int, expectedString string) {
	t.Helper()

	assertEqual(t, expectedValue, int(obfuscationMode))
	assertEqual(t, expectedString, obfuscationMode.String())
}
