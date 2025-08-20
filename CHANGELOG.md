# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html) except to the first release.

## [Unreleased]

### Added

### Changed

### Fixed

## [v0.1.0] - 2025-08-20

The release introduces base option types: wrappers for all types supported by
the msgpack/v5 library and a generic type.

### Added

- Implemented optional types for builtin go types int*, uint*, float*,
  bytes, string, bool.
- Implemented generic optional type for any go type.
