package obfuscate

import "strings"

// HTTPHeaderObfuscator represents an object that can obfuscate HTTP header values.
type HTTPHeaderObfuscator struct {
	obfuscators map[string]Obfuscator
}

// ObfuscateHeaderValue obfuscates the value for a single header.
func (o HTTPHeaderObfuscator) ObfuscateHeaderValue(headerName, headerValue string) string {
	if obfuscator, ok := o.obfuscators[strings.ToLower(headerName)]; ok {
		return obfuscator.ObfuscateString(headerValue)
	}
	return headerValue
}

// ObfuscateHeaderValues obfuscates multiple values for a single header.
func (o HTTPHeaderObfuscator) ObfuscateHeaderValues(headerName string, headerValues []string) []string {
	if obfuscator, ok := o.obfuscators[strings.ToLower(headerName)]; ok {
		var result []string
		for _, headerValue := range headerValues {
			result = append(result, obfuscator.ObfuscateString(headerValue))
		}
		return result
	}
	return append([]string{}, headerValues...)
}

// ObfuscateHeaderMap obfuscates all values in a map where the keys are HTTP header names and the values are the matching HTTP header values.
func (o HTTPHeaderObfuscator) ObfuscateHeaderMap(headerMap map[string]string) map[string]string {
	result := map[string]string{}
	for headerName, headerValue := range headerMap {
		result[headerName] = o.ObfuscateHeaderValue(headerName, headerValue)
	}
	return result
}

// ObfuscateHeaderMultiMap obfuscates all values in a map where the keys are HTTP header names and the values are the matching HTTP header values.
func (o HTTPHeaderObfuscator) ObfuscateHeaderMultiMap(headerMap map[string][]string) map[string][]string {
	result := map[string][]string{}
	for headerName, headerValue := range headerMap {
		result[headerName] = o.ObfuscateHeaderValues(headerName, headerValue)
	}
	return result
}

// HTTPHeaders creates a new HTTP header obfuscator.
func HTTPHeaders(obfuscators map[string]Obfuscator) HTTPHeaderObfuscator {
	obfuscatorMap := map[string]Obfuscator{}
	for headerName, obfuscator := range obfuscators {
		obfuscatorMap[strings.ToLower(headerName)] = obfuscator
	}
	return HTTPHeaderObfuscator{obfuscators: obfuscatorMap}
}
