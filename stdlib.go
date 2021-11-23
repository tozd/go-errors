package errors

import (
	stderrors "errors"
	"fmt"
)

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
//
// An error type might provide an Is method so it can be treated as equivalent
// to an existing error. For example, if MyError defines
//
//	func (m MyError) Is(target error) bool { return target == fs.ErrExist }
//
// then Is(MyError{}, fs.ErrExist) returns true. See syscall.Errno.Is for
// an example in the standard library.
//
// This function is a proxy for standard errors.Is.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true. Otherwise, it returns false.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// An error type might provide an As method so it can be treated as if it were a
// different error type.
//
// As panics if target is not a non-nil pointer to either a type that implements
// error, or to any interface type.
//
// This function is a proxy for standard errors.As.
func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

// Unwrap returns the result of calling the Unwrap method on err, if err's
// type contains an Unwrap method returning error.
// Otherwise, Unwrap returns nil.
//
// This function is a proxy for standard errors.Unwrap.
func Unwrap(err error) error {
	return stderrors.Unwrap(err)
}

// Base returns an error with the supplied message.
// Each call to Base returns a distinct error value even if the message is identical.
// It does not record a stack trace.
//
// Use this for a constant base error you convert to an actual error you return with
// WithStack. This base error you can then use in Is and As calls.
//
// This function is a proxy for standard errors.New.
func Base(message string) error {
	return stderrors.New(message)
}

// Basef returns an error with the supplied message
// formatted according to a format specifier.
// Each call to Basef returns a distinct error value even if the message is identical.
// It does not record a stack trace. It supports %w format verb to wrap an existing error.
//
// Use this for a constant base error you convert to an actual error you return with
// WithStack. This base error you can then use in Is and As calls. Use  %w format verb
// when you want to create a hierarchy of base errors.
//
// This function is a proxy for standard fmt.Errorf.
func Basef(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// BaseWrap returns an error with the supplied message, wrapping an existing error
// err.
// Each call to BaseWrap returns a distinct error value even if the message is identical.
// It does not record a stack trace.
//
// Use this when you want to create a hierarchy of base errors and you want to fully
// control the error message.
func BaseWrap(err error, message string) error {
	return &base{
		message,
		err,
	}
}

// BaseWrapf returns an error with the supplied message
// formatted according to a format specifier.
// Each call to BaseWrapf returns a distinct error value even if the message is identical.
// It does not record a stack trace. It does not support %w format verb. Use Basef if you need it.
//
// Use this when you want to create a hierarchy of base errors and you want to fully
// control the error message.
func BaseWrapf(err error, format string, args ...interface{}) error {
	return &base{
		fmt.Sprintf(format, args...),
		err,
	}
}

type base struct {
	msg string
	err error
}

func (b *base) Unwrap() error {
	return b.err
}

func (b *base) Error() string {
	return b.msg
}
