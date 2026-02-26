//go:build !go1.26

package errors

import "errors"

// Copied from errors/wrap.go available from Go 1.26 on.
func stderrorsAsType[E error](err error) (E, bool) { //nolint:ireturn
	if err == nil {
		var zero E
		return zero, false
	}
	var pe *E // lazily initialized.
	return asType(err, &pe)
}

// Copied from errors/wrap.go available from Go 1.26 on.
func asType[E error](err error, ppe **E) (_ E, _ bool) { //nolint:ireturn
	for {
		var e E
		if errors.As(err, &e) {
			return e, true
		}
		if x, ok := err.(interface{ As(any) bool }); ok { //nolint:inamedparam
			if *ppe == nil {
				*ppe = new(E)
			}
			if x.As(*ppe) {
				return **ppe, true
			}
		}
		{
			var x interface{ Unwrap() error }
			var x1 interface{ Unwrap() []error }
			switch {
			case errors.As(err, &x):
				err = x.Unwrap()
				if err == nil {
					return //nolint:nakedret
				}
			case errors.As(err, &x1):
				for _, err := range x1.Unwrap() {
					if err == nil {
						continue
					}
					if x1, ok := asType(err, ppe); ok {
						return x1, true
					}
				}
				return //nolint:nakedret
			default:
				return //nolint:nakedret
			}
		}
	}
}
