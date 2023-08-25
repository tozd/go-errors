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

type structWithMarshaler struct{}

func (s structWithMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`{"error":"test"}`), nil
}

func (s *structWithMarshaler) Error() string {
	return "test"
}

func TestUseMarshaler(t *testing.T) {
	t.Parallel()

	assert.False(t, useMarshaler(New("test")))
	assert.False(t, useMarshaler(&fundamentalError{})) //nolint:exhaustruct
	assert.False(t, useMarshaler(Base("test")))
	assert.True(t, useMarshaler(&structWithMarshaler{}))
	var se stringError = "test"
	assert.False(t, useMarshaler(&se))
	assert.False(t, useMarshaler(&json.MarshalerError{})) //nolint:exhaustruct
	assert.True(t, useMarshaler(&structError{}))          //nolint:exhaustruct
	assert.True(t, useMarshaler(&structParentError{}))    //nolint:exhaustruct
}
