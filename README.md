# Errors with a stack trace

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/go/errors)](https://pkg.go.dev/gitlab.com/tozd/go/errors)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/go/errors)](https://goreportcard.com/report/gitlab.com/tozd/go/errors)
[![pipeline status](https://gitlab.com/tozd/go/errors/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/go/errors/-/pipelines)
[![coverage report](https://gitlab.com/tozd/go/errors/badges/main/coverage.svg)](https://gitlab.com/tozd/go/errors/-/graphs/main/charts)

A Go package providing errors with a stack trace.

Features:

* Based of [`github.com/pkg/errors`](https://github.com/pkg/errors) with similar API, addressing many its
  [open issues](https://github.com/pkg/errors/issues). In many cases it can be used as a drop-in replacement.
  At the same time compatible with [`github.com/pkg/errors`](https://github.com/pkg/errors) errors.
* Uses standard error wrapping (available since Go 1.13).
* Provides [`errors.Errorf`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Errorf) which supports `%w` format verb
  to both wrap and record a stack trace at the same time (if not already recorded).
* Provides [`errors.E`](https://pkg.go.dev/gitlab.com/tozd/go/errors#E) type to be used instead of standard `error`
  to annotate which functions return errors with a stack trace.
* Clearly defines what are differences and expected use cases for:
  * [`errors.Errorf`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Errorf): creating a new error and recording a stack
    trace, optionally wrapping an existing error
  * [`errors.WithStack`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WithStack):
    adding a stack trace to an error without one
  * [`errors.WithMessage`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WithMessage):
    adding a prefix to the error message
  * [`errors.Wrap`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Wrap): creating a new error but recording its cause
* Provides [`errors.Base`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Base) function to create errors without
  a stack trace to be used as base errors for [`errors.Is`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Is)
  and [`errors.As`](https://pkg.go.dev/gitlab.com/tozd/go/errors#As).
* Differentiates between wrapping and recording a cause: only [`errors.Wrap`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Wrap)
  records a cause, while other functions are error transformers, wrapping the original.
* Novice friendly formatting of a stack trace when error is formatted using `%+v`:
  tells what is the order of the stack trace and what is the relation between
  wrapped errors.
* Makes sure a stack trace is not recorded multiple times unnecessarily.
* Provide optional details map on all errors returned by this package.
* Errors implement `MarshalJSON` and can be marshaled into JSON.

## Installation

This is a Go package. You can add it to your project using `go get`:

```sh
go get gitlab.com/tozd/go/errors
```

There is also a [read-only GitHub mirror available](https://github.com/tozd/go-errors),
if you need to fork the project there.

## Usage

See full package documentation with examples on [pkg.go.dev](https://pkg.go.dev/gitlab.com/tozd/go/errors#section-documentation).

## Why a new Go errors package?

[`github.com/pkg/errors`](https://github.com/pkg/errors) package is archived and not developed anymore,
with [many issues](https://github.com/pkg/errors/issues) not addressed (primarily because many require some
backward incompatible change). At the same time it has been made before
Go 1.13 added official support for wrapping errors and it does not (and cannot, in backwards compatible way)
fully embrace it. This package takes what is best from `github.com/pkg/errors`, but breaks things a bit to address
many of the open issues community has identified since then and to modernize it to today's Go:

* Message formatting `WithMessage` vs. `Wrap`: [#114](https://github.com/pkg/errors/pull/114)
* Do not re-add stack trace if one is already there: [#122](https://github.com/pkg/errors/pull/122)
* Be explicit when you want to record a stack trace again vs. do not if it already exists:
  [#75](https://github.com/pkg/errors/issues/75) [#158](https://github.com/pkg/errors/issues/158)
  [#242](https://github.com/pkg/errors/issues/242)
* `StackTrace()` should return `[]uintptr`: [#79](https://github.com/pkg/errors/issues/79)
* Do not assume `Cause` cannot return `nil`: [#89](https://github.com/pkg/errors/issues/89)
* Obtaining only message from `Wrap`: [#93](https://github.com/pkg/errors/issues/93)
* `WithMessage` always prefixes the message: [#102](https://github.com/pkg/errors/issues/102)
* Differentiate between "wrapping" and "causing": [#112](https://github.com/pkg/errors/issues/112)
* Support for base errors: [#130](https://github.com/pkg/errors/issues/130) [#160](https://github.com/pkg/errors/issues/160)
* Support for a different delimiter by supporting `Errorf`: [#207](https://github.com/pkg/errors/issues/207) [#226](https://github.com/pkg/errors/issues/226)
* Support for `Errorf` wrapping an error: [#244](https://github.com/pkg/errors/issues/244)
* Having each function wrap only once: [#223](https://github.com/pkg/errors/issues/223)

## What are main differences from `github.com/pkg/errors`?

* The `stackTracer` interface's `StackTrace()` method returns `[]uintptr` and not custom type `StackTrace`.
* All error-wrapping functions return errors which implement the standard `unwrapper` interface,
  but only `errors.Wrap` records a cause error and returns an error which implements the `causer` interface.
* All error-wrapping functions wrap the error into only one new error.
* `Errorf` supports `%w`.
* Errors formatted using `%+v` include lines `stack trace (most recent call first):` and
  `the above error was caused by the following error:` to make it clearer how is the stack
  trace formatted and how are multiple errors related to each other.
* Only `errors.Wrap` always records the stack trace while other functions do
  not record if it is already present.
* `errors.Cause` repeatedly unwraps the error until it finds one which implements the `causer` interface,
  and then return its cause.

## It looks like `Wrap` should be named `Cause`. Why it is not?

For legacy reasons because this package builds on shoulders of `github.com/pkg/errors`.
Every modification to errors made through this package is done through wrapping
so that original error is always available. `Wrap` wraps the error to records the cause.

## Related projects

* [cockroachdb/errors](https://github.com/cockroachdb/errors) – Go errors
  with every possible feature you might ever need in your large project.
  This package aims to stay lean and be more or less just a drop-in replacement
  for core Go errors, but with stack traces (and few utility functions for common
  cases).
