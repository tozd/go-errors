package errors

// AsType finds the first error in err's tree that matches the type E, and
// if one is found, returns that error value and true. Otherwise, it
// returns the zero value of E and false.
//
// The tree consists of err itself, followed by the errors obtained by
// repeatedly calling its Unwrap() error or Unwrap() []error method. When
// err wraps multiple errors, AsType examines err followed by a
// depth-first traversal of its children.
//
// An error err matches the type E if the type assertion err.(E) holds,
// or if the error has a method As(any) bool such that err.As(target)
// returns true when target is a non-nil *E. In the latter case, the As
// method is responsible for setting target.
//
// This function is a proxy for standard errors.AsType.
func AsType[E error](err error) (E, bool) { //nolint:ireturn
	return stderrorsAsType[E](err)
}
