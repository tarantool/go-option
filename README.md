<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

[![Go Reference][godoc-badge]][godoc-url]
[![Actions Status][actions-badge]][actions-url]
[![Code Coverage][coverage-badge]][coverage-url]
[![Telegram EN][telegram-badge]][telegram-en-url]
[![Telegram RU][telegram-badge]][telegram-ru-url]

# go-option: library to work with optional types

## Pre-generated basic optional types

## Gentype Utility

A Go code generator for creating optional types with MessagePack
serialization support.

### Overview

Gentype generates wrapper types for various Go primitives and
custom types that implement optional (some/none) semantics with
full MessagePack serialization capabilities. These generated types 
are useful for representing values that may or may not be present,
while ensuring proper encoding and decoding when using MessagePack.

### Features

- Generates optional types for built-in types (bool, int, float, string, etc.)
- Supports custom types with MessagePack extension serialization
- Provides common optional type operations:
  - `SomeXxx(value)` - Create an optional with a value
  - `NoneXxx()` - Create an empty optional
  - `Unwrap()`, `UnwrapOr()`, `UnwrapOrElse()` - Value extraction
  - `IsSome()`, `IsNone()` - Presence checking
- Full MessagePack `CustomEncoder` and `CustomDecoder` implementation
- Type-safe operations

### Installation

```bash
go install github.com/tarantool/go-option/cmd/gentypes@latest
# OR (for go version 1.24+)
go get -tool github.com/tarantool/go-option/cmd/gentypes@latest
```

### Usage

#### Generating Optional Types

To generate optional types for existing types in a package:

```bash
gentypes -package ./path/to/package -ext-code 123
# OR (for go version 1.24+)
go tool gentypes -package ./path/to/package -ext-code 123
```

Or you can use it to generate file from go:
```go
//go:generate go run github.com/tarantool/go-option/cmd/gentypes@latest -ext-code 123
// OR (for go version 1.24+)
//go:generate go tool gentypes -ext-code 123
```

Flags:

 • `-package`: Path to the Go package containing types to wrap (default: `"."`)
 • `-ext-code`: MessagePack extension code to use for custom types (must be between
   -128 and 127, no default value)
 • `-verbose`: Enable verbose output (default: `false`)

#### Using Generated Types

Generated types follow the pattern Optional<TypeName> and provide methods for working
with optional values:

```go
// Create an optional with a value.
opt := SomeOptionalString("hello")

// Check if a value is present.
if opt.IsSome() {
    value := opt.Unwrap()
    fmt.Println(value)
}

// Use a default value if none.
value := opt.UnwrapOr("default")

// Encode to MessagePack.
err := opt.EncodeMsgpack(encoder)
```

[godoc-badge]: https://pkg.go.dev/badge/github.com/tarantool/go-option.svg
[godoc-url]: https://pkg.go.dev/github.com/tarantool/go-option
[actions-badge]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml/badge.svg
[actions-url]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml
[coverage-badge]: https://img.shields.io/coverallsCoverage/github/tarantool/go-option
[coverage-url]: https://coveralls.io/github/tarantool/go-option?branch=master
[telegram-badge]: https://img.shields.io/badge/Telegram-join%20chat-blue.svg
[telegram-en-url]: http://telegram.me/tarantool
[telegram-ru-url]: http://telegram.me/tarantoolru
