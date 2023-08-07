package errors

import (
	"bytes"
	"encoding/json"
	"unsafe"
)

func marshalWithoutEscapeHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	b := buf.Bytes()
	if len(b) > 0 {
		return b[:len(b)-1], nil
	}
	return b, nil
}

// marshalJSONError marshals foreign errors.
func marshalJSONError(err error) ([]byte, E) {
	var e error
	var s stack
	if stackErr, ok := err.(stackTracer); ok {
		s = stackErr.StackTrace()
	} else if stackErr, ok := err.(pkgStackTracer); ok {
		st := stackErr.StackTrace()
		s = stack(*(*[]uintptr)(unsafe.Pointer(&st)))
	}
	var jsonCause []byte
	cause := Cause(err)
	if cause != nil {
		jsonCause, e = marshalWithoutEscapeHTML(cause)
		if e != nil {
			return nil, WithStack(e)
		}
		if len(jsonCause) == 0 || bytes.Equal(jsonCause, []byte("{}")) {
			var eStack E
			jsonCause, eStack = marshalJSONError(cause)
			if eStack != nil {
				return nil, eStack
			}
		}
		if bytes.Equal(jsonCause, []byte("{}")) {
			jsonCause = []byte{}
		}
	}
	jsonErr, e := marshalWithoutEscapeHTML(&struct {
		Error string          `json:"error,omitempty"`
		Stack stack           `json:"stack,omitempty"`
		Cause json.RawMessage `json:"cause,omitempty"`
	}{
		Error: err.Error(),
		Stack: s,
		Cause: jsonCause,
	})
	if e != nil {
		return nil, WithStack(e)
	}
	return jsonErr, nil
}

// marshalJSONAnyError marshals our and foreign errors.
func marshalJSONAnyError(err error) ([]byte, E) {
	var eStack E
	jsonWrap, e := marshalWithoutEscapeHTML(err)
	if e != nil {
		return nil, WithStack(e)
	}
	if len(jsonWrap) == 0 || bytes.Equal(jsonWrap, []byte("{}")) {
		jsonWrap, eStack = marshalJSONError(err)
		if eStack != nil {
			return nil, eStack
		}
	}
	if bytes.Equal(jsonWrap, []byte("{}")) {
		jsonWrap = []byte{}
	}
	return jsonWrap, nil
}
