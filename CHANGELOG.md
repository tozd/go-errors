# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Errors returned by this package provide also optional details map accessible
  through `detailer` interface.
- `WithDetails` which wraps an error exposing access to (a potentially new layer of)
  details about the error.

## [0.3.0] - 2022-01-03

### Changed

- Change license to Apache 2.0.

## [0.2.0] - 2021-12-01

### Changed

- `errors.Cause` handles better `Cause` which returns `nil`.
- JSON marshaling of foreign errors uses `errors.Cause`.

## [0.1.0] - 2021-11-30

### Added

- First public release.

[unreleased]: https://gitlab.com/tozd/go/errors/-/compare/v0.3.0...main
[0.3.0]: https://gitlab.com/tozd/go/errors/-/compare/v0.2.0...v0.3.0
[0.2.0]: https://gitlab.com/tozd/go/errors/-/compare/v0.1.0...v0.2.0
[0.1.0]: https://gitlab.com/tozd/go/errors/-/tags/v0.1.0

<!-- markdownlint-disable-file MD024 -->
