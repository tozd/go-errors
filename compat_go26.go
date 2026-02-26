//go:build go1.26

package errors

import (
	stderrors "errors"
)

var stderrorsAsType = stderrors.AsType //nolint:gochecknoglobals
