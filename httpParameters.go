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

// strings.Builder's WriteString method is documented to return a nil error, so no need to check for it in methods below

// ObfuscateParameterString obfuscates the given string.
//
// It is like [HTTPParameterObfuscator.ObfuscateString], but it returns any error instead of handling it internally.
func (o HTTPParameterObfuscator) ObfuscateParameterString(s string) (string, error) {
	builder := strings.Builder{}
	err := o.obfuscateParameterString(s, &builder)
	return builder.String(), err
}

// ObfuscateString implements the [Obfuscator] interface.
//
// It is like [HTTPParameterObfuscator.ObfuscateParameterString], but it handles any error internally according to the [ErrorStrategy]
// provided when the HTTPParameterObfuscator instance was created.
func (o HTTPParameterObfuscator) ObfuscateString(s string) string {
	builder := strings.Builder{}
	err := o.obfuscateParameterString(s, &builder)
	if err != nil {
		switch o.onError {
		case OnErrorLog:
			o.printf("ObfuscateString error: %v\n", err)
		case OnErrorInclude:
			builder.WriteString(fmt.Sprintf("<error: %v>", err))
		case OnErrorStop:
			break
		default:
			o.panicf("ObfuscateString error: %v", err)
		}
	}
	return builder.String()
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

// UntilLength implements the [Obfuscator] interface.
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
