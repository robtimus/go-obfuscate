package obfuscate

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// HTTPParameterObfuscator represents an object that can obfuscate HTTP query and form parameter strings,
// as well as separate parameter values.
type HTTPParameterObfuscator struct {
	obfuscators map[string]Obfuscator
	onError     ErrorStrategy
	printf      func(format string, v ...any)
	panicf      func(format string, v ...any)
}

// HTTPParameterObfuscatorOptions represents the configurable options used by HTTPParameterObfuscator instances.
type HTTPParameterObfuscatorOptions struct {
	// OnError represents the strategy to follow when an error occurs while obfuscating a string.
	OnError ErrorStrategy
	// Logger represents the optional logger to use in case OnError is OnErrorLog or OnErrorPanic.
	Logger *log.Logger
}

// ObfuscateParameter obfuscates the given value for a parameter with the given name.
func (o HTTPParameterObfuscator) ObfuscateParameter(name, value string) string {
	if obfuscator, ok := o.obfuscators[name]; ok {
		return obfuscator.ObfuscateString(value)
	}
	return value
}

// ObfuscateString implements the Obfuscator interface.
func (o HTTPParameterObfuscator) ObfuscateString(s string) string {
	builder := strings.Builder{}
	index := strings.Index(s, "&")
	for index != -1 {
		if !o.obfuscateParameter(s[:index], &builder) {
			return builder.String()
		}
		builder.WriteString("&")
		s = s[index+1:]
		index = strings.Index(s, "&")
	}
	o.obfuscateParameter(s, &builder)
	return builder.String()
}

func (o HTTPParameterObfuscator) obfuscateParameter(s string, builder *strings.Builder) bool {
	index := strings.Index(s, "=")
	// strings.Builder's WriteString method is documented to return a nil error, so no need to check for it
	if index == -1 {
		builder.WriteString(s)
	} else {
		name, err := url.QueryUnescape(s[:index])
		if err != nil {
			o.handleError(err, builder)
			return false
		}
		builder.WriteString(s[:index+1])

		value, err := url.QueryUnescape(s[index+1:])
		if err != nil {
			o.handleError(err, builder)
			return false
		}
		builder.WriteString(o.ObfuscateParameter(name, value))
	}
	return true
}

func (o HTTPParameterObfuscator) handleError(err error, builder *strings.Builder) {
	switch o.onError {
	case OnErrorLog:
		o.printf("ObfuscateString error: %v\n", err)
		break
	case OnErrorInclude:
		builder.WriteString(fmt.Sprintf("<error: %v>", err))
		break
	case OnErrorStop:
		break
	default:
		o.panicf("ObfuscateString error: %v", err)
	}
}

// UntilLength implements the Obfuscator interface.
func (o HTTPParameterObfuscator) UntilLength(prefixLength int) ObfuscatorPrefix {
	return NewObfuscatorPrefix(o, prefixLength)
}

var defaultPrintf = func(format string, v ...any) {
	fmt.Printf(format, v...)
}

var defaultPanicf = func(format string, v ...any) {
	log.Panicf(format, v...)
}

// HTTPParameters creates a new HTTP parameter obfuscator.
func HTTPParameters(obfuscators map[string]Obfuscator, options *HTTPParameterObfuscatorOptions) HTTPParameterObfuscator {
	obfuscatorMap := map[string]Obfuscator{}
	for headerName, obfuscator := range obfuscators {
		obfuscatorMap[strings.ToLower(headerName)] = obfuscator
	}

	var onError ErrorStrategy
	printf := defaultPrintf
	panicf := defaultPanicf
	if options != nil {
		onError = options.OnError
		printf = options.Logger.Printf
		panicf = options.Logger.Panicf
	}

	return HTTPParameterObfuscator{obfuscators: obfuscatorMap, onError: onError, printf: printf, panicf: panicf}
}
