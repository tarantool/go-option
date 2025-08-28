<a href="http://tarantool.org">
	<img src="https://avatars2.githubusercontent.com/u/2344919?v=2&s=250" align="right">
</a>

[![Go Reference][godoc-badge]][godoc-url]
[![Actions Status][actions-badge]][actions-url]
[![Code Coverage][coverage-badge]][coverage-url]
[![Telegram EN][telegram-badge]][telegram-en-url]
[![Telegram RU][telegram-badge]][telegram-ru-url]

[godoc-badge]: https://pkg.go.dev/badge/github.com/tarantool/go-option.svg
[godoc-url]: https://pkg.go.dev/github.com/tarantool/go-option
[actions-badge]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml/badge.svg
[actions-url]: https://github.com/tarantool/go-option/actions/workflows/testing.yaml
[coverage-badge]: https://img.shields.io/coverallsCoverage/github/tarantool/go-option
[coverage-url]: https://coveralls.io/github/tarantool/go-option?branch=master
[telegram-badge]: https://img.shields.io/badge/Telegram-join%20chat-blue.svg
[telegram-en-url]: http://telegram.me/tarantool
[telegram-ru-url]: http://telegram.me/tarantoolru

# option

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
* [Run tests](#run-tests)
* [Development](#development)

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

```Go
package main

import (
	"bytes"
	"fmt"

	"github.com/tarantool/go-option"
	msgpack "github.com/vmihailenco/msgpack/v5"
)

func main() {

	var buf bytes.Buffer

	enc := msgpack.NewEncoder(&buf)
	dec := msgpack.NewDecoder(&buf)

	someUint := option.SomeUint(12)
	err := someUint.EncodeMsgpack(enc)
	if err != nil {
		panic("encode fail")
	}

	var unmarshaled option.Uint
	err = unmarshaled.DecodeMsgpack(dec)
	if err != nil {
		panic(fmt.Errorf("decode fail: %s", err))
	}

	if !unmarshaled.IsSome() {
		panic("IsSome error")
	}

	if unmarshaled.Unwrap() != 12 {
		panic("Unwrap error")
	}
}

```


## Run tests

To run default set of tests:

```shell
go test ./... -count=1
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
