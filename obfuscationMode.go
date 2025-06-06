package obfuscate

import "fmt"

type ObfuscationMode int

// String implements the fmt.Stringer interface.
func (om ObfuscationMode) String() string {
	switch om {
	case Exclude:
		return "Exclude"
	case ExcludeAll:
		return "ExcludeAll"
	case Inherit:
		return "Inherit"
	case InheritOverridable:
		return "InheritOverridable"
	}
	return fmt.Sprintf("ObfuscationMode(%d)", int(om))
}

const (
	// Exclude means obfuscators will not be used for complex structures like JSON objects or arrays.
	// Obfuscation will instead traverse in them, and for any nested properties their own obfuscation rules will be applied.
	Exclude ObfuscationMode = iota
	// ExcludeAll means obfuscators will not be used for complex structures like JSON objects or arrays.
	// Obfuscation will not traverse in them, so any nested properties will not be obfuscated.
	ExcludeAll
	// Inherit means obfuscators will be used for any nested property of complex structures like JSON objects or arrays.
	// Any obfuscation rules for nested properties will be ignored.
	Inherit
	// InheritOverridable means obfuscators will be used for any nested property of complex structures like JSON objects or arrays.
	// For any nested properties their own obfuscation rules will be applied.
	InheritOverridable
)

func obfuscationModeOrDefault(obfuscationMode *ObfuscationMode, defaultObfuscationMode ObfuscationMode) ObfuscationMode {
	if obfuscationMode == nil {
		return defaultObfuscationMode
	}
	return *obfuscationMode
}
