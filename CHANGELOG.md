# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- `StackFormat` to format the provided stack trace as text.
- `StackMarshalJSON` marshals the provided stack trace as JSON.
- Support wrapping multiple errors with `errors.Join`.
  [#4](https://gitlab.com/tozd/go/errors/-/issues/4)

## Changed

- `Wrap` behaves like `New` and `Wrapf` like `Errorf` if provided error is nil
  instead of returning `nil`.
  [#2](https://gitlab.com/tozd/go/errors/-/issues/2)
- Package is tested only on Go 1.16 and newer.
- Lines `stack trace (most recent call first):` and
  `the above error was caused by the following error:` changed to lower case.

## [0.4.1] - 2022-04-21

### Fixed

- Initialize details when calling `WithDetails` to prevent race conditions.

## [0.4.0] - 2022-04-20

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

[unreleased]: https://gitlab.com/tozd/go/errors/-/compare/v0.4.1...main
[0.4.1]: https://gitlab.com/tozd/go/errors/-/compare/v0.4.0...v0.4.1
[0.4.0]: https://gitlab.com/tozd/go/errors/-/compare/v0.3.0...v0.4.0
[0.3.0]: https://gitlab.com/tozd/go/errors/-/compare/v0.2.0...v0.3.0
[0.2.0]: https://gitlab.com/tozd/go/errors/-/compare/v0.1.0...v0.2.0
[0.1.0]: https://gitlab.com/tozd/go/errors/-/tags/v0.1.0

<!-- markdownlint-disable-file MD024 -->
