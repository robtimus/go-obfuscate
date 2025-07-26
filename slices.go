package obfuscate

import "fmt"

// Slice obfuscates all elements of a slice using the given obfuscator.
// Values are converted to string using fmt.Sprintf("%v", value).
// This allows this function to be used not just for string slices but also other scalar types.
func Slice[T any](s []T, obfuscator Obfuscator) []string {
	result := make([]string, len(s))
	for index, value := range s {
		result[index] = obfuscator.ObfuscateString(fmt.Sprintf("%v", value))
	}

	return result
}
