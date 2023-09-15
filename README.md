# Errors with a stack trace

[![pkg.go.dev](https://pkg.go.dev/badge/gitlab.com/tozd/go/errors)](https://pkg.go.dev/gitlab.com/tozd/go/errors)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/tozd/go/errors)](https://goreportcard.com/report/gitlab.com/tozd/go/errors)
[![pipeline status](https://gitlab.com/tozd/go/errors/badges/main/pipeline.svg?ignore_skipped=true)](https://gitlab.com/tozd/go/errors/-/pipelines)
[![coverage report](https://gitlab.com/tozd/go/errors/badges/main/coverage.svg)](https://gitlab.com/tozd/go/errors/-/graphs/main/charts)

A Go package providing errors with a stack trace and optional structured details.

Features:

- Based of [`github.com/pkg/errors`](https://github.com/pkg/errors) with compatible API, addressing many its
  [open issues](https://github.com/pkg/errors/issues). In can be used as a drop-in replacement and even mixed with
  `github.com/pkg/errors`.
- Uses standard error wrapping (available since Go 1.13) and wrapping of multiple errors (available since Go 1.20).
- Provides [`errors.Errorf`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Errorf) which supports `%w` format verb
  to both wrap and record a stack trace at the same time (if not already recorded).
- Provides [`errors.E`](https://pkg.go.dev/gitlab.com/tozd/go/errors#E) type to be used instead of standard `error`
  to annotate which functions return errors with a stack trace and details.
- Clearly defines what are differences and expected use cases for:
  - [`errors.Errorf`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Errorf): creating a new error and recording a stack
    trace, optionally wrapping an existing error
  - [`errors.WithStack`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WithStack):
    adding a stack trace to an error without one
  - [`errors.WithMessage`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WithMessage):
    adding a prefix to the error message
  - [`errors.Wrap`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Wrap) and
    [`errors.WrapWith`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Wrap): creating a new error but recording its cause
  - [`errors.Join`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Join): joining multiple
    errors which happened during execution (e.g., additional errors which happened during
    cleanup)
- Provides [`errors.Base`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Base) function to create errors without
  a stack trace to be used as base errors for [`errors.Is`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Is)
  and [`errors.As`](https://pkg.go.dev/gitlab.com/tozd/go/errors#As).
- Differentiates between wrapping and recording a cause: [`errors.Wrap`](https://pkg.go.dev/gitlab.com/tozd/go/errors#Wrap)
  and [`errors.WrapWith`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WrapWith)
  record a cause, while other functions are error transformers, wrapping the original.
- Makes sure a stack trace is not recorded multiple times unnecessarily.
- Provides optional details map on all errors returned by this package. Helper
  [`errors.WithDetails`](https://pkg.go.dev/gitlab.com/tozd/go/errors#WithDetails) allows both recording a stack trace
  and annotating an error with details at the same time.
- [Errors](https://pkg.go.dev/gitlab.com/tozd/go/errors#Formatter) and
  [stack traces](https://pkg.go.dev/gitlab.com/tozd/go/errors#StackFormatter) support configurable formatting
  and can be marshaled into JSON.
  Both formatting and JSON marshaling is supported also for errors not made using this package.
- Limited [JSON unmarshal of errors](https://pkg.go.dev/gitlab.com/tozd/go/errors#UnmarshalJSON) is supported to
  enable formatting of JSON errors.

## Installation

This is a Go package. You can add it to your project using `go get`:

```sh
go get gitlab.com/tozd/go/errors
```

## Usage

See full package documentation with examples on [pkg.go.dev](https://pkg.go.dev/gitlab.com/tozd/go/errors#section-documentation).

## Why a new Go errors package?

[`github.com/pkg/errors`](https://github.com/pkg/errors) package is archived and not developed anymore,
with [many issues](https://github.com/pkg/errors/issues) not addressed (primarily because many require some
backward incompatible change). At the same time it has been made before
Go 1.13 added official support for wrapping errors and it does not (and cannot, in backwards compatible way)
fully embrace it. This package takes what is best from `github.com/pkg/errors`, but breaks things a bit to address
many of the open issues community has identified since then and to modernize it to today's Go:

- Message formatting `WithMessage` vs. `Wrap`: [#114](https://github.com/pkg/errors/pull/114)
- Do not re-add stack trace if one is already there: [#122](https://github.com/pkg/errors/pull/122)
- Be explicit when you want to record a stack trace again vs. do not if it already exists:
  [#75](https://github.com/pkg/errors/issues/75) [#158](https://github.com/pkg/errors/issues/158)
  [#242](https://github.com/pkg/errors/issues/242)
- `StackTrace()` should return `[]uintptr`: [#79](https://github.com/pkg/errors/issues/79)
- Do not assume `Cause` cannot return `nil`: [#89](https://github.com/pkg/errors/issues/89)
- Obtaining only message from `Wrap`: [#93](https://github.com/pkg/errors/issues/93)
- `WithMessage` always prefixes the message: [#102](https://github.com/pkg/errors/issues/102)
- Differentiate between "wrapping" and "causing": [#112](https://github.com/pkg/errors/issues/112)
- Support for base errors: [#130](https://github.com/pkg/errors/issues/130) [#160](https://github.com/pkg/errors/issues/160)
- Support for a different delimiter by supporting `Errorf`: [#207](https://github.com/pkg/errors/issues/207) [#226](https://github.com/pkg/errors/issues/226)
- Support for `Errorf` wrapping an error: [#244](https://github.com/pkg/errors/issues/244)
- Having each function wrap only once: [#223](https://github.com/pkg/errors/issues/223)

## What are main differences from `github.com/pkg/errors`?

In most cases this package can be used as a drop-in replacement for `github.com/pkg/errors`,
but there are some small (behavioral) differences (i.e., improvements):

- The `stackTracer` interface's `StackTrace()` method returns `[]uintptr` and not custom type `StackTrace`.
- All error-wrapping functions return errors which implement the standard `unwrapper` interface,
  but only `errors.Wrap` records a cause error and returns an error which implements the `causer` interface.
- All error-wrapping functions wrap the error into only one new error.
- Only `errors.Wrap` always records the stack trace while other functions do
  not record if it is already present.
- `errors.Cause` repeatedly unwraps the error until it finds one which implements the `causer` interface,
  and then return its cause.

Main additions are:

- `Errorf` supports `%w`.
- This package supports annotating errors with additional key-value details.
- This package provides more configurable formatting and JSON marshaling of stack traces and errors.
- Support for base errors (e.g., `errors.Base` and `errors.WrapWith`)
  instead of just operating on error messages.

## How to migrate from `github.com/pkg/errors`?

1. Replace your `github.com/pkg/errors` imports with `gitlab.com/tozd/go/errors` imports.
2. Run `go get gitlab.com/tozd/go/errors`.

In most cases, this should be it. Consider:

- Migrating to using `errors.Errorf` (which supports `%w`)
  as it subsumes most of `github.com/pkg/errors`'s error constructors.
- Using `errors.E` return type instead of plain `error`. This makes Go type system
  help you to not return an error without a stack trace.
- Using structured details with `errors.WithDetails`.
- [Using base errors](#suggestion) instead of
  message-based error constructors, especially if you want callers of your functions to
  handle different errors differently.

## How to migrate from standard `errors` and `fmt.Errorf`?

1. Replace your `errors` imports with `gitlab.com/tozd/go/errors` imports.
2. Replace your `fmt.Errorf` calls with `errors.Errorf` calls.
3. If you have top-level `errors.New` or `fmt.Errorf`/`errors.Errorf` calls to create
   variables with base errors, use `errors.Base` and `errors.Basef` instead.
   Otherwise your
   [stack traces will not be recorded correctly](#i-use-base-errors-but-stack-traces-do-not-look-right).

This is it. Now all your errors automatically record stack traces. You can now
print those stack traces when [formatting errors](https://pkg.go.dev/gitlab.com/tozd/go/errors#Formatter)
using `fmt.Printf("%+v", err)` or marshal them to JSON.

Consider:

- Using structured details with `errors.WithDetails`.
- Using `errors.E` return type instead of plain `error`. This makes Go type system
  help you to not return an error without a stack trace.
- [Using base errors](#suggestion) instead of
  message-based error constructors, especially if you want callers of your functions to
  handle different errors differently.

This package provides everything `errors` package does (but with stack traces and optional details) in
a compatible way, often simply proxying calls to the standard `errors` package, but there is a difference
how `errors.Join` operates: if only one non-nil error is provided, `errors.Join`
returns it as-is without wrapping it. Furthermore, `errors.Join` records a stack trace at the
point it was called.

## How should this package be used?

Patterns for errors in Go have evolved through time, with Go 1.13 introducing error wrapping
(and Go 1.20 wrapping of multiple errors) enabling standard annotation of errors with
additional information. For example, standard `fmt.Errorf` allows annotating the error
message of a base error with additional information, e.g.:

```go
import (
  "errors"
  "fmt"
)

var MyErr = errors.New("my error")

// later on

err := fmt.Errorf(`user "%s" made an error: %w`, username, MyErr)
```

This is great because one can later on extract the cause of the error using `errors.Is` without
having to parse the error message itself:

```go
if errors.Is(err, MyErr) {
  ...
}
```

But if `MyErr` can happen at multiple places it is hard to debug without additional
information where this error happened. One can add it to the error message, but this is tedious:

```go
err := fmt.Errorf(`user "%s" made an error during sign-in: %w`, username, MyErr)

// somewhere else

err := fmt.Errorf(`user "%s" made an error during sign-out: %w`, username, MyErr)
```

Furthermore, if you need to extract the username you again have to parse the error message.

Instead, this package provides support for recording a stack trace optional
structured details:

```go
import (
  "gitlab.com/tozd/go/errors"
)

var MyErr = errors.Base("my error")
var UserErr = errors.Basef("user made an error: %w", MyErr)

// later on

err := errors.WithDetails(UserErr, "username", username)
```

Here, `err` contains a descriptive error message (without potentially sensitive
information), a stack trace, and easy to extract username.
You can display all that to the developer using `fmt.Printf`:

```go
fmt.Printf("error: %#+v", err)
```

Extracting username is easy:

```go
errors.AllDetails(err)["username"]
```

And if you need to know exactly which error happened (e.g., to show a translated
error message to the end user), you can use `errors.Is` or similar logic (to map
between errors and their user-friendly translated error messages).

The package provides `errors.Errorf`, `errors.Wrap` and other message-based
error constructors, but they are provided primarily for compatibility
(and to support patterns different from the suggested one below).

### Suggestion

- Use `errors.Base`, `errors.Basef`, `errors.BaseWrap` and `errors.BaseWrapf`
  to create a tree of constant base errors.
- Do not use `errors.Errorf` but use `errors.WithDetails` to record both a stack
  trace and any additional structured details at point where the error occurs.
- You can use `errors.WithStack` or `errors.WithDetails` to do the same with
  errors coming from outside of your codebase, as soon as possible.
- If you want to map one error to another while recording the cause,
  use `errors.WrapWith`.
- If errors coming from outside of your codebase do not provide adequate base errors,
  use `errors.WrapWith` as well to provide them yourself.
- If during handling of an error another error occurs (e.g., in `defer` during cleanup)
  use `errors.Join` or `errors.Errorf` (but only to control how messages are joined)
  to join them all.
- Use `errors.E` return type instead of plain `error`. This makes Go type system
  help you to not return an error without a stack trace.

Your functions should return only errors for which you provide base errors as well.
Those base errors become part of your API.
All additional (non-constant) information about the particular error goes into its
stack trace and details.

Remember, error messages, stack traces, and details are for developers not end users.
Be mindful if and how you expose them to end users.

## I use base errors, but stack traces do not look right?

Let's assume you use `errors.New` to create a base error:

```go
var myBaseErr = errors.New("error")
```

Later on, you want to annotate it with a stack trace and return it from the function
where the error occurred:

```go
func run() error.E {
  // ... do something ...
  return errors.WithStack(myBaseErr)
}
```

Sadly, this does not work as expected. If you use `errors.New`
(or `errors.Errorf`) to create a base error, a stack trace is recorded at the
point it was called. `errors.WithStack` then does nothing because it detects that
a stack trace already exists.
So the stack trace you get is from the point you created the base error.

You should create base errors using `errors.Base`, `errors.Basef`, `errors.BaseWrap` and
`errors.BaseWrapf`. They create errors without stack traces. And `errors.WithStack`
adds a stack trace at the point it was called.

## It looks like `Wrap` should be named `Cause` or `WithCause`. Why it is not?

For legacy reasons because this package builds on shoulders of `github.com/pkg/errors`.
Every modification to errors made through this package is done through wrapping
so that original error is always available. `Wrap` wraps the error to records the cause.
`Cause` exist as a helper to return the recorded cause.

## Related projects

- [github.com/cockroachdb/errors](https://github.com/cockroachdb/errors) – Go errors
  with every possible feature you might ever need in your large project.
  This package aims to stay lean and be more or less just a drop-in replacement
  for core Go errors and archived `github.com/pkg/errors`, but with stack traces
  and structured details (and few utility functions for common cases).
- [github.com/friendsofgo/errors](https://github.com/friendsofgo/errors) – A fork of
  `github.com/pkg/errors` but beyond updating the code to error wrapping introduced in
  Go 1.13 it does not seem maintained.
- [github.com/go-errors/errors](https://github.com/go-errors/errors) – Another small error
  package with stack trace recording support, but with different API than
  `github.com/pkg/errors`. It does not support structured details, extended formatting
  nor JSON marshal.

## GitHub mirror

There is also a [read-only GitHub mirror available](https://github.com/tozd/go-errors),
if you need to fork the project there.
