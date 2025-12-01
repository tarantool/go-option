# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html) except to the first release.

## [Unreleased]

### Added

### Changed

### Fixed

## [v1.1.0] - 2025-12-02

The release introduces the `option.Any` type.

### Added

- The `option.Any` type is a wrapper for the `any` type (#14).

## [v1.0.0] - 2025-09-09

This release introduces code generator `gentypes`.

### Added

- **Support for third-party types**: A new code generator `gentypes` is
  introduced. It allows generating optional types for any user-defined or
  third-party types. The generator can be invoked using `go:generate`.

## [v0.1.0] - 2025-08-20

The release introduces base option types: wrappers for all types supported by
the msgpack/v5 library and a generic type.

### Added

- Implemented optional types for builtin go types int*, uint*, float*,
  bytes, string, bool.
- Implemented generic optional type for any go type.
