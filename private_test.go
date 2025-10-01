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

type structWithMarshalerError struct{}

func (s structWithMarshalerError) MarshalJSON() ([]byte, error) {
	return []byte(`{"error":"test"}`), nil
}

func (s *structWithMarshalerError) Error() string {
	return "test"
}

func TestUseMarshaler(t *testing.T) {
	t.Parallel()

	assert.False(t, useMarshaler(New("test")))
	assert.False(t, useMarshaler(&fundamentalError{}))
	assert.False(t, useMarshaler(Base("test")))
	assert.True(t, useMarshaler(&structWithMarshalerError{}))
	var se stringError = "test"
	assert.False(t, useMarshaler(&se))
	assert.False(t, useMarshaler(&json.MarshalerError{}))
	assert.True(t, useMarshaler(&structError{}))
	assert.True(t, useMarshaler(&structParentError{}))
}
