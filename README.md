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

 * `-package`: Path to the Go package containing types to wrap (default: `"."`)
 * `-ext-code`: MessagePack extension code to use for custom types (must be between
   -128 and 127, no default value)
 * `-verbose`: Enable verbose output (default: `false`)
 * `-force`: Ignore absence of marshal/unmarshal methods on type (default: `false`).
   Helpful for types from third-party modules.
 * `-imports`: Add imports to generated file (default is empty).
   Helpful for types from third-party modules.
 * `-marshal-func`: func that should do marshaling (default is `MarshalMsgpack` method).
   Helpful for types from third-party modules.
   Should be func of type `func(v T) ([]byte, error)` and should
   be located in the same dir or should be imported.
 * `-unmarshal-func`: func that should do unmarshalling (default is `UnmarshalMsgpack` method).
   Helpful for types from third-party modules.
   Should be func of type `func(v *T, data []byte) error` and should
   be located in the same dir or should be imported.

#### Generating Optional Types for Third-Party Modules

Sometimes you need to generate an optional type for a type from a third-party module,
and you can't add `MarshalMsgpack`/`UnmarshalMsgpack` methods to it.
In this case, you can use the `-force`, `-imports`, `-marshal-func`, and `-unmarshal-func` flags.

For example, to generate an optional type for `github.com/google/uuid.UUID`:

1.  Create a file with marshal and unmarshal functions for the third-party type.
    For example, `uuid.go`:

    ```go
    package main

    import (
        "errors"

        "github.com/google/uuid"
    )

    func encodeUUID(uuid uuid.UUID) ([]byte, error) {
        return uuid[:], nil
    }

    var (
        ErrInvalidLength = errors.New("invalid length")
    )

    func decodeUUID(uuid *uuid.UUID, data []byte) error {
        if len(data) != len(uuid) {
            return ErrInvalidLength
        }
        copy(uuid[:], data)
        return nil
    }
    ```

2.  Use the following `go:generate` command:

    ```go
    //go:generate go run github.com/tarantool/go-option/cmd/gentypes@latest -package . -imports "github.com/google/uuid" -type UUID -marshal-func "encodeUUID" -unmarshal-func "decodeUUID" -force -ext-code 100
    ```

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

## Benchmarking

Along with the approach supplied with `go-option` library pointer-based and slice-based approaches were benchmarked as well.

#### Init + Get (empty value)

```
# int
BenchmarkNoneInt/Typed-8        	560566400	 2.200 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneInt/Generic-8      	543332625	 2.193 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneInt/GenericPtr-8   	487631254	 2.474 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneInt/GenericSlice-8 	441513422	 2.608 ns/op	 0 B/op	  0 allocs/op
# string
BenchmarkNoneString/Typed-8        	170894025	 6.545 ns/op	 0 B/op	 0 allocs/op
BenchmarkNoneString/Generic-8      	185572758	 6.451 ns/op	 0 B/op	 0 allocs/op
BenchmarkNoneString/GenericPtr-8   	159143874	 7.459 ns/op	 0 B/op	 0 allocs/op
BenchmarkNoneString/GenericSlice-8 	173419598	 6.708 ns/op	 0 B/op	 0 allocs/op
# struct
BenchmarkNoneStruct/Typed-8        	384845384	 3.107 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneStruct/Generic-8      	415633797	 2.884 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneStruct/GenericPtr-8   	331620082	 3.580 ns/op	 0 B/op	  0 allocs/op
BenchmarkNoneStruct/GenericSlice-8 	387593746	 3.115 ns/op	 0 B/op	  0 allocs/op
```

#### Init + Get (non-empty value)

```
# int
BenchmarkSomeInt/Typed-8        	499550200	 2.231 ns/op	 0 B/op	  0 allocs/op
BenchmarkSomeInt/Generic-8      	321369986	 3.491 ns/op	 0 B/op	  0 allocs/op
BenchmarkSomeInt/GenericPtr-8   	 64221356	 16.03 ns/op	 8 B/op	  1 allocs/op
BenchmarkSomeInt/GenericSlice-8 	 71858188	 16.53 ns/op	 8 B/op	  1 allocs/op
# string
BenchmarkSomeString/Typed-8        	192472155	 5.840 ns/op	  0 B/op	 0 allocs/op
BenchmarkSomeString/Generic-8      	197161162	 6.471 ns/op	  0 B/op	 0 allocs/op
BenchmarkSomeString/GenericPtr-8   	 16207524	 98.67 ns/op	 16 B/op	 1 allocs/op
BenchmarkSomeString/GenericSlice-8 	 12426998	 100.4 ns/op	 16 B/op	 1 allocs/op
# struct
BenchmarkSomeStruct/Typed-8          	358631294	 3.407 ns/op	  0 B/op	   0 allocs/op
BenchmarkSomeStruct/Generic-8        	241312274	 4.978 ns/op	  0 B/op	   0 allocs/op
BenchmarkSomeStruct/GenericPtr-8     	 32534370	 33.28 ns/op	 24 B/op	 1 allocs/op
BenchmarkSomeStruct/GenericSlice-8   	 34119435	 33.08 ns/op	 24 B/op	 1 allocs/op
```

At this point we can see already that the alternatives (based on pointer and slice) require allocations while the approach implemented in `go-option` doesn't.

Now let's check encoding and decoding.

## Encode + Decode

```
# int
BenchmarkEncodeDecodeInt/Typed-8        	46089481	 22.66 ns/op	  0 B/op	 0 allocs/op
BenchmarkEncodeDecodeInt/Generic-8      	10070619	 119.6 ns/op	 32 B/op	 2 allocs/op
BenchmarkEncodeDecodeInt/GenericPtr-8   	20202076	 58.14 ns/op	 16 B/op	 2 allocs/op
BenchmarkEncodeDecodeInt/GenericSlice-8 	17400481	 66.24 ns/op	 24 B/op	 3 allocs/op
# string
BenchmarkEncodeDecodeString/Typed-8        	 6053182	 191.4 ns/op	  8 B/op	 1 allocs/op
BenchmarkEncodeDecodeString/Generic-8      	 1891269	 668.3 ns/op	 56 B/op	 3 allocs/op
BenchmarkEncodeDecodeString/GenericPtr-8   	 1645518	 659.2 ns/op	 56 B/op	 4 allocs/op
BenchmarkEncodeDecodeString/GenericSlice-8 	 1464177	 775.4 ns/op	 72 B/op	 5 allocs/op
# struct
BenchmarkEncodeDecodeStruct/Typed-8        	12816339	 90.85 ns/op	  3 B/op	 1 allocs/op
BenchmarkEncodeDecodeStruct/Generic-8      	 2304001	 532.5 ns/op	 67 B/op	 3 allocs/op
BenchmarkEncodeDecodeStruct/GenericPtr-8   	 2071520	 570.2 ns/op	 75 B/op	 4 allocs/op
BenchmarkEncodeDecodeStruct/GenericSlice-8 	 2007445	 587.4 ns/op	 99 B/op	 5 allocs/op
```

As it can be seen generic implementation ~3-4 times slower than the typed one. Thus it is recommended to use pre-generated optionals for basic types supplied with `go-option` (`option.Int`, `option.String` etc.).

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
