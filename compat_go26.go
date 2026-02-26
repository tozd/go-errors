//go:build go1.26

package errors

import (
	stderrors "errors"
)

func stderrorsAsType[E error](err error) (E, bool) { //nolint:ireturn
	return stderrors.AsType[E](err)
}
