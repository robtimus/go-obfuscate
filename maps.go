package obfuscate

// MapObfuscator represents an object that can obfuscate map values.
type MapObfuscator[K comparable] struct {
	obfuscators map[K]Obfuscator
}

// ObfuscateMap obfuscates all values in a map.
func (o MapObfuscator[K]) ObfuscateMap(m map[K]string) map[K]string {
	result := map[K]string{}
	for key, value := range m {
		if obfuscator, ok := o.obfuscators[key]; ok {
			result[key] = obfuscator.ObfuscateString(value)
		} else {
			result[key] = value
		}
	}
	return result
}

// ObfuscateMultiMap obfuscates all values in a map.
func (o MapObfuscator[K]) ObfuscateMultiMap(m map[K][]string) map[K][]string {
	result := map[K][]string{}
	for key, values := range m {
		if obfuscator, ok := o.obfuscators[key]; ok {
			var obfuscatedValues []string
			for _, value := range values {
				obfuscatedValues = append(obfuscatedValues, obfuscator.ObfuscateString(value))
			}
			result[key] = obfuscatedValues
		} else {
			result[key] = append([]string{}, values...)
		}
	}
	return result
}

// Maps creates a new map obfuscator.
func Maps[K comparable](obfuscators map[K]Obfuscator) MapObfuscator[K] {
	obfuscatorMap := map[K]Obfuscator{}
	for key, obfuscator := range obfuscators {
		obfuscatorMap[key] = obfuscator
	}
	return MapObfuscator[K]{obfuscators: obfuscatorMap}
}
