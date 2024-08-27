package obfuscate

import (
	"fmt"
	"testing"
)

func TestObfuscateHeaderValue(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
	})
	parameters := []struct {
		headerName, headerValue, expected string
	}{
		{"Content-Type", "application/json", "application/json"},
		{"Content-Length", "13", "13"},
		{"Authorization", "Bearer someToken", "***"},
	}
	for i := range parameters {
		headerName := parameters[i].headerName
		headerValue := parameters[i].headerValue
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s': '%s'", headerName, headerValue), func(t *testing.T) {
			actual := obfuscator.ObfuscateHeaderValue(headerName, headerValue)
			if actual != expected {
				t.Errorf("expected: '%s', actual: '%s'", expected, actual)
			}
		})
	}
}

func TestObfuscateHeaderValues(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
		"MULTIVALUED":   Portion().KeepAtEnd(2).Build(),
	})
	parameters := []struct {
		headerName             string
		headerValues, expected []string
	}{
		{"Content-Type", []string{"application/json"}, []string{"application/json"}},
		{"Content-Length", []string{"13"}, []string{"13"}},
		{"Authorization", []string{"Bearer someToken"}, []string{"***"}},
		{"MultiValued", []string{"value1", "value2"}, []string{"****e1", "****e2"}},
	}
	for i := range parameters {
		headerName := parameters[i].headerName
		headerValues := parameters[i].headerValues
		expected := parameters[i].expected
		t.Run(fmt.Sprintf("applied to '%s': '%s'", headerName, headerValues), func(t *testing.T) {
			actual := obfuscator.ObfuscateHeaderValues(headerName, headerValues)
			if arraysDiffer(actual, expected) {
				t.Errorf("expected: '%s', actual: '%s'", expected, actual)
			}
		})
	}
}

func TestObfuscateHeaderMap(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
	})
	headerMap := map[string]string{
		"Content-Type":   "application/json",
		"Content-Length": "13",
		"Authorization":  "Bearer someToken",
	}
	expected := map[string]string{
		"Content-Type":   "application/json",
		"Content-Length": "13",
		"Authorization":  "***",
	}

	actual := obfuscator.ObfuscateHeaderMap(headerMap)
	if mapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateNilHeaderMap(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
	})
	expected := map[string]string{}

	actual := obfuscator.ObfuscateHeaderMap(nil)
	if mapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateHeaderMultiMap(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
		"MULTIVALUED":   Portion().KeepAtEnd(2).Build(),
	})
	headerMap := map[string][]string{
		"Content-Type":   {"application/json"},
		"Content-Length": {"13"},
		"Authorization":  {"Bearer someToken"},
		"MultiValued":    {"value1", "value2"},
	}
	expected := map[string][]string{
		"Content-Type":   {"application/json"},
		"Content-Length": {"13"},
		"Authorization":  {"***"},
		"MultiValued":    {"****e1", "****e2"},
	}

	actual := obfuscator.ObfuscateHeaderMultiMap(headerMap)
	if multiMapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}

func TestObfuscateNilHeaderMultiMap(t *testing.T) {
	obfuscator := HTTPHeaders(map[string]Obfuscator{
		"AUTHORIZATION": WithFixedLength(3),
	})
	expected := map[string][]string{}

	actual := obfuscator.ObfuscateHeaderMultiMap(nil)
	if multiMapsDiffer(actual, expected) {
		t.Errorf("expected: '%s', actual: '%s'", expected, actual)
	}
}
