package errors

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type stringError string

func (e *stringError) Error() string {
	return string(*e)
}

type structError struct {
	Foo string `json:"foo"`
}

func (s *structError) Error() string {
	return s.Foo
}

type structParentError struct {
	structError
}

func TestSupportsJSON(t *testing.T) {
	t.Parallel()

	assert.True(t, supportsJSON(New("test")))
	assert.True(t, supportsJSON(&fundamentalError{})) //nolint:exhaustruct
	assert.False(t, supportsJSON(Base("test")))
	var se stringError = "test"
	assert.False(t, supportsJSON(&se))
	assert.False(t, supportsJSON(&json.MarshalerError{})) //nolint:exhaustruct
	assert.True(t, supportsJSON(&structError{}))          //nolint:exhaustruct
	assert.True(t, supportsJSON(&structParentError{}))    //nolint:exhaustruct
}
