# go-obfuscate
[![Build Status](https://github.com/robtimus/go-obfuscate/actions/workflows/build.yml/badge.svg)](https://github.com/robtimus/go-obfuscate/actions/workflows/build.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-obfuscate&metric=alert_status)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-obfuscate)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=robtimus%3Ago-obfuscate&metric=coverage)](https://sonarcloud.io/summary/overall?id=robtimus%3Ago-obfuscate)
[![Known Vulnerabilities](https://snyk.io/test/github/robtimus/go-obfuscate/badge.svg)](https://snyk.io/test/github/robtimus/go-obfuscate)

Provides functionality for obfuscating text. This can be useful for logging information that contains sensitive information.

## Obfuscating strings

### Pre-defined functions

The following pre-defined functions are provided that all return an immutable obfuscator.

#### obfuscate.All / obfuscate.AllWithMask

Replaces all characters with a mask character.

```go
obfuscator := obfuscate.All()
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "***********"
```

Note: using these obfuscators still leak out information about the length of the original text. Using a fixed length or fixed value (see below) is more secure.

#### obfuscate.WithFixedLength / obfuscate.WithFixedLengthWithMask

Replaces the entire text with a fixed number of mask characters.

```go
obfuscator := obfuscate.WithFixedLength(5)
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "*****"
```

#### obfuscate.WithFixedValue

Replaces the entire text with a fixed value.

```go
obfuscator := obfuscate.WithFixedValue("foo")
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "foo"
```

#### obfuscate.Portion

While the above examples are simple, they are not very flexible. Using `obfuscate.Portion` you can build obfuscators that obfuscate only specific portions of text. Some examples:

##### Obfuscating all but the last 4 characters

Useful for obfuscating values like credit card numbers.

```go
obfuscator := obfuscate.Portion().KeepAtEnd(4).Build()
obfuscated := obfuscator.ObfuscateString("1234567890123456")
// obfuscated is "************3456"
```

It’s advised to use `AtLeastFromStart`, to make sure that values of fewer than 16 characters are still obfuscated properly:

```go
obfuscator := obfuscate.Portion().KeepAtEnd(4).AtLeastFromStart(12).Build()
obfuscated := obfuscator.ObfuscateString("1234567890")
// obfuscated is "**********" and not "******7890"
```

##### Obfuscating only the last 2 characters

Useful for obfuscating values like zip codes, where the first part is not as sensitive as the full zip code:

```go
obfuscator := obfuscate.Portion().KeepAtStart(math.MaxInt).AtLeastFromEnd(2).Build()
obfuscated := obfuscator.ObfuscateString("SW1A 2AA")
// obfuscated is "SW1A 2**"
```

Here, the `KeepAtStart` instructs the obfuscator to keep everything; however, `AtLeastFromEnd` overrides that partly to ensure that the last two characters are obfuscated regardless of the value specified by `KeepAtStart`.

##### Using a fixed length

Similar to using `obfuscate.All`, by default an obfuscator built using `obfuscate.Portion` leaks out the length of the original text. If your text has a variable length, you should consider specifying a fixed total length for the result. The length of the result will then be the same no matter how long the input is:

```go
obfuscator := obfuscate.Portion().KeepAtStart(2).KeepAtEnd(2).FixedTotalLength(6).Build()
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "He**ld"
obfuscated = obfuscator.ObfuscateString("foo")
// obfuscated is "fo**oo"
```

Note that if `KeepAtStart` and `KeepAtEnd` are both specified, parts of the input may be repeated in the result if the input’s length is less than the combined number of characters to keep. This makes it harder to find the original input. For example, if in the example `foo` would be obfuscated into `fo***o` instead, it would be clear that the input was `foo`. Instead, it can now be anything that starts with `fo` and ends with `oo`.

#### obfuscate.NewObfuscator

Converts any function that takes a string as input and returns a string into an obfuscator.

```go
obfuscator := obfuscate.NewObfuscator(strings.ToUpper)
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "HELLO WORLD"
```

#### obfuscate.None

Returns the input as-is. It can be used as default to prevent checks. For instance:

```go
obfuscator := somePossiblyNilObfuscator
if obfuscator == nil {
    obfuscator = obfuscate.None()
}
obfuscated := obfuscator.ObfuscateString("Hello World")
// obfuscated is "Hello World" if somePossiblyNilObfuscator was nil
```

### Combining obfuscators

Sometimes the obfucators in this module alone cannot perform the obfuscation you need. For instance, if you want to obfuscate credit cards, but keep the first and last 4 characters. If the credit cards are all fixed length, `obfuscate.Portion` can do just that:

```go
obfuscator := obfuscate.Portion().KeepAtStart(4).KeepAtEnd(4).Build()
obfuscated := obfuscator.ObfuscateString("1234567890123456")
// obfuscated is "1234********3456"
```

However, if you attempt to use such an obfuscator on only a part of a credit card, you could end up leaking parts of the credit card that you wanted to obfuscate:

```go
incorrectlyObfuscated := obfuscator.ObfuscateString("12345678901234")
// incorrectlyObfuscated is "1234******1234" where "1234********34" would probably be preferred
```

To overcome this issue, it’s possible to combine obfuscators. The form is as follows:
* Specify the first obfuscator, and the input length to which it should be used.
* Specify any other obfuscators, and the input lengths to which they should be used. Note that each input length should be larger than the previous input length.
* Specify the obfuscator that will be used for the remainder.

For instance, for credit card numbers of exactly 16 characters, the above can also be written like this:

```go
obfuscator := obfuscate.None().UntilLength(4).Then(obfuscate.All()).UntilLength(12).Then(obfuscate.None())
```

With this chaining, it’s now possible to keep the first and last 4 characters, but with at least 8 characters in between:

```go
obfuscator := obfuscate.None().UntilLength(4).Then(obfuscate.Portion().KeepAtEnd(4).AtLeastFromStart(8).Build())
obfuscated := obfuscator.ObfuscateString("12345678901234")
// obfuscated is "1234********34"
```

### Splitting text during obfuscation

To make it easier to create obfuscators for structured text like email addresses, use an `obfuscate.SplitPoint`. Three implementations are provided :
* `obfuscate.AtFirst(s)` splits at the first occurrence of string `s`.
* `obfuscate.AtLast(s)` splits at the last occurrence of string `s`.
* `obfuscate.AtNth(s, occurrence)` splits at the zero-based specified occurrence of string `s`.

For instance:

```go
// Keep the domain as-is
localPartObfuscator := obfuscate.Portion().KeepAtStart(1).KeepAtEnd(1).FixedTotalLength(8).Build()
domainObfuscator := obfuscate.None()
obfuscator := obfuscate.AtFirst("@").SplitTo(localPartObfuscator, domainObfuscator)
obfuscated := obfuscator.ObfuscateString("test@example.org")
// obfuscated is "t******t@example.org"
```

To obfuscate the domain except for the TLD, use a nested `obfuscate.SplitPoint`:

```go
// Keep only the TLD of the domain
localPartObfuscator := obfuscate.Portion().KeepAtStart(1).KeepAtEnd(1).FixedTotalLength(8).Build()
domainObfuscator := obfuscate.AtLast(".").SplitTo(obfuscate.All(), obfuscate.None())
obfuscator := obfuscate.AtFirst("@").SplitTo(localPartObfuscator, domainObfuscator)
obfuscated := obfuscator.ObfuscateString("test@example.org")
// obfuscated is "t******t@*******.org"
```

## Obfuscating HTTP headers

Use `obfuscate.HTTPHeaders` to create an object that can obfuscate single HTTP headers (as strings or string slices) and maps of HTTP headers (with values as strings or string slices). It will always use case insensitive matching.

```go
headerObfuscator := obfuscate.HTTPHeaders(map[string]obfuscate.Obfuscator{
    "Authorization": obfuscate.WithFixedLength(3),
})
obfuscatedAuthorization := headerObfuscator.ObfuscateHeaderValue("authorization", "Bearer someToken")
// obfuscatedAuthorization is "***"
obfuscatedAuthorizations := headerObfuscator.ObfuscateHeaderValues("authorization", []string{"Bearer someToken"})
// obfuscatedAuthorization is ["***"]
obfuscatedContentType := headerObfuscator.ObfuscateHeaderValue("Content-Type", "application/json")
// obfuscatedContentType is "application/json"
obfuscatedHeaders := headerObfuscator.ObfuscateHeaderMap(map[string]string{
    "authorization": "Bearer someToken",
    "content-type": "application/json",
})
// obfuscatedHeaders is map["authorization":"***", "content-type":"application/json"]
```

## Obfuscating HTTP parameters

Use `obfuscate.HTTPParameters` to create an object that can obfuscate HTTP query and form parameter strings, as well as separate parameter values.

```go
paramsObfuscator := obfuscate.HTTPParameters().
    WithParameter("password", obfuscate.WithFixedLength(3)).
    Build()
obfuscatedPassword := paramsObfuscator.ObfuscateParameter("password", "admin1234")
// obfuscatedPassword is "***"
obfuscatedUsername := paramsObfuscator.ObfuscateParameter("username", "admin")
// obfuscatedUsername is "admin"
obfuscatedParamString, err := paramsObfuscator.ParseAndObfuscateString("username=admin&password=admin1234")
// obfuscatedParamString is "username=admin&password=***"
```

Objects created by calling `Build()` on the result of `obfuscate.HTTPParameters` also implement `obfuscate.Obfuscator`. This works almost the same as calling `ParseAndObfuscateString`. Because parsing parameter strings can fail with an error, the builder returned by `obfuscate.HTTPParameters` can be configured to specify how to handle errors:

* `OnErrorLog(logger)` (default) will cause the error to be logged. If a non-`nil` [`log.Logger`](https://pkg.go.dev/log#Logger) is given its [`Printf`](https://pkg.go.dev/log#Logger.Printf) method will be used, otherwise [`fmt.Printf`](https://pkg.go.dev/fmt#Printf) will be used.
* `OnErrorInclude()` will cause the error to be included in the return value.
* `OnErrorDiscard()` wil cause the return value to not contain any data following the error. For security purposes parameter names and values will either be included fully or not at all.

```go
paramsObfuscator := obfuscate.HTTPParameters().
    WithParameter("password", obfuscate.WithFixedLength(3)).
    OnErrorInclude().
    Build()
obfuscatedParamString := paramsObfuscator.ObfuscateString("username=admin&password=admin1234%A")
// obfuscatedParamString is something like "username=admin&password=<error: invalid URL escape \"%A\">"
```

## Obfuscating JSON

Use `obfuscate.JSON` to create an object that can obfuscate JSON strings.

```go
jsonObfuscator := obfuscate.JSON().
    WithProperty("password", obfuscate.WithFixedLength(3), nil).
    Build()
obfuscatedJsonString, err := jsonObfuscator.ParseAndObfuscateString(`{"username": "admin", "password": "admin1234"}`)
// obfuscatedJsonString is equivalent to `{"username": "admin", "password": "***"}`
```

If a matched property is not a string or other scalar value but instead an object or array, it will by default be ignored. This behaviour can be changed in two ways:

1. Per property. Instead of providing `nil` for the third argument, pass property-specific options:
    ```go
    jsonObfuscator := obfuscate.JSON().
        WithProperty("password", obfuscate.WithFixedLength(3), &obfuscate.JSONPropertyObfuscationOptions{
            ForObjects: obfuscate.Inherit,
            ForArrays:  obfuscate.Inherit,
        }).
        Build()
    ```
2. Setting global settings on the builder:
    ```go
    jsonObfuscator := obfuscate.JSON().
        WithProperty("password", obfuscate.WithFixedLength(3), nil).
        ForObjects(obfuscate.Inherit).
        ForArrays(obfuscate.Inherit).
        Build()
    ```

In both cases, `ForObjects` and `ForArrays` can take the following values:
* `obfuscate.Exclude` to not match properties with object or array values; nested properties will be matched separately.
* `obfuscate.ExcludeAll` to not match properties with object or array values; obfuscation will exclude all nested properties as well.
* `obfuscate.Inherit` to obfuscate each nested scalar property value or array element using the given obfuscator.
* `obfuscate.InheritOverridable` to obfuscate each nested scalar property value or array element using the given obfuscator; however, if a nested property has its own obfuscator defined this will be used instead.

Objects created by calling `Build()` on the result of `obfuscate.JSON` also implement `obfuscate.Obfuscator`. This works almost the same as calling `ParseAndObfuscateString`. Because parsing JSON strings can fail with an error, the builder returned by `obfuscate.JSON` can be configured to specify how to handle errors:

* `OnErrorLog(logger)` (default) will cause the error to be logged. If a non-`nil` [`log.Logger`](https://pkg.go.dev/log#Logger) is given its [`Printf`](https://pkg.go.dev/log#Logger.Printf) method will be used, otherwise [`fmt.Printf`](https://pkg.go.dev/fmt#Printf) will be used.
* `OnErrorInclude()` will cause the error to be included in the return value.
* `OnErrorDiscard()` wil cause the return value to not contain any data following the error. For security purposes parameter names and values will either be included fully or not at all.

```go
jsonObfuscator := obfuscate.JSON().
    WithProperty("password", obfuscate.WithFixedLength(3), nil).
    OnErrorInclude().
    Build()
obfuscatedJsonString := jsonObfuscator.ObfuscateString(`{"username": admin, "password": "admin1234"}`)
// obfuscatedJsonString ends with "<error: invalid character 'a' looking for beginning of value>"
// It may or may not include any part of the JSON before the parse error
```

## Obfuscating maps

Use `obfuscate.Maps` to create an object that can obfuscate string values in maps. The key type can be anything that is comparable, but usually `string` is used.

```go
mapObfuscator := obfuscate.Maps(map[string]obfuscate.Obfuscator{
    "password": obfuscate.WithFixedLength(3),
})
obfuscatedMap := mapObfuscator.ObfuscateMap(map[string]string{
    "username": "admin",
    "password": "admin1234",
})
// obfuscatedMap is map["username":"admin", "password":"***"]
```

## Obfuscating slices

Use `obfuscate.Slice` to obfuscate all elements of a slice.

```go
input := []int{1, 2, 3}
obfuscated := obfuscate.Slice(input, obfuscate.WithFixedLength(3))
// obfuscated is ["***", "***", "***"]
```
