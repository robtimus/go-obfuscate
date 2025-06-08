package obfuscate

import (
	"fmt"
	"log"
	"maps"
	"net/url"
	"strings"
)

// HTTPParameterObfuscator represents an object that can obfuscate HTTP query and form parameter strings,
// as well as separate parameter values.
type HTTPParameterObfuscator struct {
	parameters map[string]Obfuscator
	onError    ErrorStrategy
	logger     *log.Logger
}

// ObfuscateParameter obfuscates the given value for a parameter with the given name.
func (o HTTPParameterObfuscator) ObfuscateParameter(name, value string) string {
	if obfuscator, ok := o.parameters[name]; ok {
		return obfuscator.ObfuscateString(value)
	}
	return value
}

// strings.Builder's WriteString method is documented to return a nil error, so no need to check for it in methods below

// ParseAndObfuscateString implements the [Obfuscator] interface.
func (o HTTPParameterObfuscator) ParseAndObfuscateString(s string) (string, error) {
	builder := strings.Builder{}
	err := o.obfuscateParameterString(s, &builder)
	return builder.String(), err
}

// ObfuscateString implements the [Obfuscator] interface.
//
// It is like [HTTPParameterObfuscator.ParseAndObfuscateString], but it handles any error internally according to the [ErrorStrategy]
// provided when the HTTPParameterObfuscator instance was created.
func (o HTTPParameterObfuscator) ObfuscateString(s string) string {
	builder := strings.Builder{}
	err := o.obfuscateParameterString(s, &builder)
	if err != nil {
		switch o.onError {
		case OnErrorLog:
			logError(o.logger, "ObfuscateString error: %v\n", err)
		case OnErrorInclude:
			builder.WriteString(fmt.Sprintf("<error: %v>", err))
		case OnErrorDiscard:
			break
		}
	}
	return builder.String()
}

// UntilLength implements the [Obfuscator] interface.
func (o HTTPParameterObfuscator) UntilLength(prefixLength int) ObfuscatorPrefix {
	return NewObfuscatorPrefix(o, prefixLength)
}

func (o HTTPParameterObfuscator) obfuscateParameterString(s string, builder *strings.Builder) error {
	index := strings.Index(s, "&")
	for index != -1 {
		err := o.obfuscateParameter(s[:index], builder)
		if err != nil {
			return err
		}
		builder.WriteString("&")
		s = s[index+1:]
		index = strings.Index(s, "&")
	}
	return o.obfuscateParameter(s, builder)
}

func (o HTTPParameterObfuscator) obfuscateParameter(s string, builder *strings.Builder) error {
	index := strings.Index(s, "=")
	if index == -1 {
		builder.WriteString(s)
	} else {
		name, err := url.QueryUnescape(s[:index])
		if err != nil {
			return err
		}
		builder.WriteString(s[:index+1])

		value, err := url.QueryUnescape(s[index+1:])
		if err != nil {
			return err
		}
		builder.WriteString(o.ObfuscateParameter(name, value))
	}
	return nil
}

// HTTPParameterObfuscatorBuilder is a builder for [HTTPParameterObfuscator] instances.
type HTTPParameterObfuscatorBuilder struct {
	parameters map[string]Obfuscator
	onError    ErrorStrategy
	logger     *log.Logger
}

// HTTPParameters creates a new builder for HTTP parameter obfuscators.
func HTTPParameters() *HTTPParameterObfuscatorBuilder {
	return &HTTPParameterObfuscatorBuilder{parameters: map[string]Obfuscator{}}
}

// WithParameter registers a parameter to obfuscate. It uses the given obfuscator for obfuscating any occurrence of a parameter with the given name.
func (b *HTTPParameterObfuscatorBuilder) WithParameter(parameterName string, obfuscator Obfuscator) *HTTPParameterObfuscatorBuilder {
	b.parameters[parameterName] = obfuscator
	return b
}

// OnErrorLog sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorLog].
// If the given logger is nil, [fmt.Printf] will be used instead.
func (b *HTTPParameterObfuscatorBuilder) OnErrorLog(logger *log.Logger) *HTTPParameterObfuscatorBuilder {
	b.onError = OnErrorLog
	b.logger = logger
	return b
}

// OnErrorInclude sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorInclude].
func (b *HTTPParameterObfuscatorBuilder) OnErrorInclude() *HTTPParameterObfuscatorBuilder {
	b.onError = OnErrorInclude
	b.logger = nil
	return b
}

// OnErrorDiscard sets the strategy to follow when an error occurs while obfuscating a string to [OnErrorDiscard].
func (b *HTTPParameterObfuscatorBuilder) OnErrorDiscard() *HTTPParameterObfuscatorBuilder {
	b.onError = OnErrorDiscard
	b.logger = nil
	return b
}

// Build creates a new HTTP parameter obfuscator using the contents of the builder.
func (b *HTTPParameterObfuscatorBuilder) Build() HTTPParameterObfuscator {
	return HTTPParameterObfuscator{
		parameters: maps.Clone(b.parameters),
		onError:    b.onError,
		logger:     b.logger,
	}
}
