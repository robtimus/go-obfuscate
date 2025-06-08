package obfuscate

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
)

// JSONObfuscator represents an object that can obfuscate JSON strings.
type JSONObfuscator struct {
	properties map[string]jsonPropertyConfig
	forObjects ObfuscationMode
	forArrays  ObfuscationMode
	onError    ErrorStrategy
	logger     *log.Logger
}

// ParseAndObfuscateString implements the [Obfuscator] interface.
func (o JSONObfuscator) ParseAndObfuscateString(s string) (string, error) {
	var parsed any
	err := json.Unmarshal([]byte(s), &parsed)
	if err != nil {
		return "", err
	}

	obfuscated := o.obfuscateWithDefault(parsed, nil)

	obfuscatedJson, err := json.MarshalIndent(obfuscated, "", "  ")
	return string(obfuscatedJson), err
}

// ObfuscateString implements the [Obfuscator] interface.
//
// It is like [JSONObfuscator.ParseAndObfuscateString], but it handles any error internally according to the [ErrorStrategy]
// provided when the JSONObfuscator instance was created.
func (o JSONObfuscator) ObfuscateString(s string) string {
	obfuscated, err := o.ParseAndObfuscateString(s)
	if err != nil {
		switch o.onError {
		case OnErrorLog:
			logError(o.logger, "ObfuscateString error: %v\n", err)
		case OnErrorInclude:
			obfuscated = fmt.Sprintf("%s<error: %v>", obfuscated, err)
		case OnErrorDiscard:
			break
		}
	}
	return obfuscated
}

// UntilLength implements the [Obfuscator] interface.
func (o JSONObfuscator) UntilLength(prefixLength int) ObfuscatorPrefix {
	return NewObfuscatorPrefix(o, prefixLength)
}

func (o JSONObfuscator) obfuscateScalar(value any, obfuscator Obfuscator) any {
	if obfuscator == nil {
		return value
	}
	s := fmt.Sprintf("%v", value)
	return obfuscator.ObfuscateString(s)
}

func (o JSONObfuscator) obfuscateScalars(object any, obfuscator Obfuscator) any {
	if s, ok := object.([]any); ok {
		result := make([]any, len(s))
		for index, value := range s {
			result[index] = o.obfuscateScalars(value, obfuscator)
		}
		return result
	}
	if m, ok := object.(map[string]any); ok {
		result := map[string]any{}
		for name, value := range m {
			result[name] = o.obfuscateScalars(value, obfuscator)
		}
		return result
	}
	return o.obfuscateScalar(o, obfuscator)
}

func (o JSONObfuscator) obfuscateWithDefault(object any, defaultObfuscator Obfuscator) any {
	if s, ok := object.([]any); ok {
		result := make([]any, len(s))
		for index, value := range s {
			result[index] = o.obfuscateWithDefault(value, defaultObfuscator)
		}
		return result
	}
	if m, ok := object.(map[string]any); ok {
		result := map[string]any{}
		for name, value := range m {
			var config *jsonPropertyConfig
			obfuscator := defaultObfuscator
			if conf, ok := o.properties[name]; ok {
				config = &conf
				obfuscator = conf.obfuscator
			}
			result[name] = o.obfuscateMapEntry(value, obfuscator, defaultObfuscator, config)
		}
		return result
	}
	return o.obfuscateScalar(object, defaultObfuscator)
}

func (o JSONObfuscator) obfuscateMapEntry(value any, obfuscator, defaultObfuscator Obfuscator, config *jsonPropertyConfig) any {
	if obfuscator == nil {
		return o.obfuscateWithDefault(value, defaultObfuscator)
	}
	if config == nil {
		// obfuscator == defaultObfuscator
		return o.obfuscateWithDefault(value, obfuscator)
	}
	if _, ok := value.([]any); ok {
		obfuscationMode := obfuscationModeOrDefault(config.forArrays, o.forArrays)
		return o.obfuscateValue(value, obfuscator, defaultObfuscator, obfuscationMode)
	}
	if _, ok := value.(map[string]any); ok {
		obfuscationMode := obfuscationModeOrDefault(config.forObjects, o.forObjects)
		return o.obfuscateValue(value, obfuscator, defaultObfuscator, obfuscationMode)
	}
	return o.obfuscateScalar(value, obfuscator)
}

func (o JSONObfuscator) obfuscateValue(value any, obfuscator, defaultObfuscator Obfuscator, obfuscationMode ObfuscationMode) any {
	switch obfuscationMode {
	case Exclude:
		return o.obfuscateWithDefault(value, defaultObfuscator)
	case ExcludeAll:
		return value
	case Inherit:
		return o.obfuscateScalars(value, obfuscator)
	case InheritOverridable:
		return o.obfuscateWithDefault(value, obfuscator)
	default:
		log.Panicf("Unsupported ObfuscationMode: %s", obfuscationMode)
		return value
	}
}

// JSONObfuscatorBuilder is a builder for [JSONObfuscator] instances.
type JSONObfuscatorBuilder struct {
	properties map[string]jsonPropertyConfig
	forObjects ObfuscationMode
	forArrays  ObfuscationMode
	onError    ErrorStrategy
	logger     *log.Logger
}

// JSON creates a new builder for JSON obfuscators.
func JSON() *JSONObfuscatorBuilder {
	return &JSONObfuscatorBuilder{properties: map[string]jsonPropertyConfig{}}
}

// WithProperty registers a property to obfuscate. It uses the given obfuscator for obfuscating any occurrence of a property with the given name.
//
// When the value of a property with the given name is an object or array, by default it will be ignored.
// Additional options can be given that determine how to handle these. See [ObfuscationMode] for more information.
func (b *JSONObfuscatorBuilder) WithProperty(propertyName string, obfuscator Obfuscator, options *JSONPropertyObfuscationOptions) *JSONObfuscatorBuilder {
	config := jsonPropertyConfig{obfuscator: obfuscator}
	if options != nil {
		config.forObjects = &options.ForObjects
		config.forArrays = &options.ForArrays
	}
	b.properties[propertyName] = config
	return b
}

// ForObjects sets the default obfuscation mode for objects. This is used for any property for which no explicit obfuscation mode has been given.
func (b *JSONObfuscatorBuilder) ForObjects(mode ObfuscationMode) *JSONObfuscatorBuilder {
	b.forObjects = mode
	return b
}

// ForArrays sets the default obfuscation mode for arrays. This is used for any property for which no explicit obfuscation mode has been given.
func (b *JSONObfuscatorBuilder) ForArrays(mode ObfuscationMode) *JSONObfuscatorBuilder {
	b.forArrays = mode
	return b
}

// OnErrorLog sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorLog].
// If the given logger is nil, [fmt.Printf] will be used instead.
func (b *JSONObfuscatorBuilder) OnErrorLog(logger *log.Logger) *JSONObfuscatorBuilder {
	b.onError = OnErrorLog
	b.logger = logger
	return b
}

// OnErrorInclude sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorInclude].
func (b *JSONObfuscatorBuilder) OnErrorInclude() *JSONObfuscatorBuilder {
	b.onError = OnErrorInclude
	b.logger = nil
	return b
}

// OnErrorDiscard sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorDiscard].
func (b *JSONObfuscatorBuilder) OnErrorDiscard() *JSONObfuscatorBuilder {
	b.onError = OnErrorDiscard
	b.logger = nil
	return b
}

// Build creates a new JSON obfuscator using the contents of the builder.
func (b *JSONObfuscatorBuilder) Build() JSONObfuscator {
	return JSONObfuscator{
		properties: maps.Clone(b.properties),
		forObjects: b.forObjects,
		forArrays:  b.forArrays,
		onError:    b.onError,
		logger:     b.logger,
	}
}

// JSONPropertyObfuscationOptions represents obfuscation options for a single property.
type JSONPropertyObfuscationOptions struct {
	// ForObjects defines how to handle matched object properties.
	ForObjects ObfuscationMode
	// ForArrays defines how to handle matched array properties.
	ForArrays ObfuscationMode
}

type jsonPropertyConfig struct {
	obfuscator Obfuscator
	forObjects *ObfuscationMode
	forArrays  *ObfuscationMode
}
