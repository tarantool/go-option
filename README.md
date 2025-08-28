<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

[![Go Reference][godoc-badge]][godoc-url]
[![Actions Status][actions-badge]][actions-url]
[![Code Coverage][coverage-badge]][coverage-url]
[![Telegram EN][telegram-badge]][telegram-en-url]
[![Telegram RU][telegram-badge]][telegram-ru-url]

# go-option: library to work with optional types

Package `option` implements effective and useful instruments to work
with optional types in Go. It eliminates code doubling and provides
high performance due to:
  - no memory allocations
  - serialization without reflection (at least for pre-generated types)
  - support for basic and custom types

## Table of contents

* [Installation](#installation)
* [Documentation](#documentation)
* [Quick start](#quick-start)
  * [Using pre-generated optional types](#using-pre-generated-optional-types)
  * [Usage with go-tarantool](#usage-with-go-tarantool)
* [Gentype Utility](#gentype-utility)
  * [Overview](#overview)
  * [Features](#features)
  * [Gentype installation](#gentype-installation)
  * [Generating Optional Types](#generating-optional-types)
  * [Using Generated Types](#using-generated-types)
* [Development](#development)
  * [Run tests](#run-tests)
* [License](#license)

## Installation

```shell
go install github.com/tarantool/go-option@latest
```

## Documentation

You could run the `godoc` server on `localhost:6060` with the command:

```shell
make godoc_run
```

And open the generated documentation in another terminal or use the
[link][godoc-url]:

```shell
make godoc_open
```

## Quick start

### Using pre-generated optional types

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

### Usage with go-tarantool

It may be necessary to use an optional type in a structure. For example,
to distinguish between a nil value and a missing one.

```Go
package main

import (
  "github.com/tarantool/go-option"
  tarantool "github.com/tarantool/go-tarantool/v2"
)

type User struct {
	// may be used in conjunciton with 'msgpack:",omitempty"' directive to skip fields
    _msgpack struct{} `msgpack:",asArray"` //nolint: structcheck,unused
    Name     string
    Phone    option.String
}

func main() {
    var conn *tarantool.Doer
    // Initialize tarantool connection

    // Imagine you get a slice of users from Tarantool.
    users := []User{
        {
            Name: "Nosipho Nnenne",
            Phone: option.SomeString("+15056463408"),
        },
        {
            Name:  "Maryamu Efe",
            Phone: option.NoneString(),
        },
        {
            Name: "Venera Okafor",
        },
    }

    for id, user := range users {
        conn.Do(
            tarantool.NewInsertRequest("users").Tuple(user),
        )
    }
}
```

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

### Gentype installation

```bash
go install github.com/tarantool/go-option/cmd/gentypes@latest
# OR (for go version 1.24+)
go get -tool github.com/tarantool/go-option/cmd/gentypes@latest
```

### Generating Optional Types

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

### Using Generated Types

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

## Development

You could use our Makefile targets:

```shell
make codespell
make test
make testrace
make coveralls-deps
make coveralls
make coverage
```

### Run tests

To run default set of tests directly:

```shell
go test ./... -count=1
```

## License

BSD 2-Clause License

[godoc-badge]: https://pkg.go.dev/badge/github.com/tarantool/go-option.svg
[godoc-url]: https://pkg.go.dev/github.com/tarantool/go-option
[actions-badge]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml/badge.svg
[actions-url]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml
[coverage-badge]: https://img.shields.io/coverallsCoverage/github/tarantool/go-option
[coverage-url]: https://coveralls.io/github/tarantool/go-option?branch=master
[telegram-badge]: https://img.shields.io/badge/Telegram-join%20chat-blue.svg
[telegram-en-url]: http://telegram.me/tarantool
[telegram-ru-url]: http://telegram.me/tarantoolru