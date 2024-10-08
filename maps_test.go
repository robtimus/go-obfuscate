package obfuscate

import "testing"

func TestObfuscateMap(t *testing.T) {
	obfuscator := Maps(map[string]Obfuscator{
		"key1": WithFixedLength(3),
		"KEY2": Portion().KeepAtEnd(2).Build(),
	})
	m := map[string]string{
		"key0": "value0",
		"key1": "value1",
		"key2": "value2",
		"KEY0": "VALUE0",
		"KEY1": "VALUE1",
		"KEY2": "VALUE2",
	}
	expected := map[string]string{
		"key0": "value0",
		"key1": "***",
		"key2": "value2",
		"KEY0": "VALUE0",
		"KEY1": "VALUE1",
		"KEY2": "****E2",
	}

	actual := obfuscator.ObfuscateMap(m)
	if mapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateNilMap(t *testing.T) {
	obfuscator := Maps(map[string]Obfuscator{
		"key1": WithFixedLength(3),
		"KEY2": Portion().KeepAtEnd(2).Build(),
	})
	expected := map[string]string{}

	actual := obfuscator.ObfuscateMap(nil)
	if mapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateMultiMap(t *testing.T) {
	obfuscator := Maps(map[string]Obfuscator{
		"key1": WithFixedLength(3),
		"KEY2": Portion().KeepAtEnd(2).Build(),
	})
	m := map[string][]string{
		"key0": {"value00", "value01"},
		"key1": {"value10", "value11"},
		"key2": {"value20", "value21"},
		"KEY0": {"VALUE00", "VALUE01"},
		"KEY1": {"VALUE10", "VALUE11"},
		"KEY2": {"VALUE20", "VALUE21"},
	}
	expected := map[string][]string{
		"key0": {"value00", "value01"},
		"key1": {"***", "***"},
		"key2": {"value20", "value21"},
		"KEY0": {"VALUE00", "VALUE01"},
		"KEY1": {"VALUE10", "VALUE11"},
		"KEY2": {"*****20", "*****21"},
	}

	actual := obfuscator.ObfuscateMultiMap(m)
	if multiMapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateNilMultiMap(t *testing.T) {
	obfuscator := Maps(map[string]Obfuscator{
		"key1": WithFixedLength(3),
		"KEY2": Portion().KeepAtEnd(2).Build(),
	})
	expected := map[string][]string{}

	actual := obfuscator.ObfuscateMultiMap(nil)
	if multiMapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func slicesDiffer(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return true
	}
	for i, value1 := range slice1 {
		if value1 != slice2[i] {
			return true
		}
	}
	return false
}

func mapsDiffer(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return true
	}
	for key, value1 := range map1 {
		if value2, ok := map2[key]; !ok || value1 != value2 {
			return true
		}
	}
	return false
}

func multiMapsDiffer(map1, map2 map[string][]string) bool {
	if len(map1) != len(map2) {
		return true
	}
	for key, values1 := range map1 {
		if values2, ok := map2[key]; !ok || slicesDiffer(values1, values2) {
			return true
		}
	}
	return false
}
