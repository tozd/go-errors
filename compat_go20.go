//go:build go1.20
// +build go1.20

package errors

import (
	"fmt"
)

var formatString = fmt.FormatString //nolint:gochecknoglobals
