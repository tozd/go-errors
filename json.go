package errors

import (
	"bytes"
	"encoding/json"
	"unsafe"
)

// marshalJSONError marshals foreign errors.
func marshalJSONError(err error) ([]byte, E) {
	var e error
	var s stack
	if stackErr, ok := err.(stackTracer); ok {
		s = stack(stackErr.StackTrace())
	} else if stackErr, ok := err.(pkgStackTracer); ok {
		st := stackErr.StackTrace()
		s = stack(*(*[]uintptr)(unsafe.Pointer(&st)))
	}
	var jsonCause []byte
	u, ok := err.(interface {
		Cause() error
	})
	if ok {
		cause := u.Cause()
		jsonCause, e = json.Marshal(cause)
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
	jsonErr, e := json.Marshal(&struct {
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

func (f fundamental) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: f.msg,
		Stack: f.stack,
	})
}

func (w errorf) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.msg,
		Stack: stack(w.StackTrace()),
	})
}

func (w errorfWithStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.msg,
		Stack: w.stack,
	})
}

func (w withStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.Error(),
		Stack: w.stack,
	})
}

func (w withPkgStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.Error(),
		Stack: w.StackTrace(),
	})
}

func (w wrapped) MarshalJSON() ([]byte, error) {
	jsonWrap, err := json.Marshal(w.error)
	if err != nil {
		return nil, WithStack(err)
	}
	if len(jsonWrap) == 0 || bytes.Equal(jsonWrap, []byte("{}")) {
		jsonWrap, err = marshalJSONError(w.error)
		if err != nil {
			return nil, err
		}
	}
	if bytes.Equal(jsonWrap, []byte("{}")) {
		jsonWrap = []byte{}
	}
	return json.Marshal(&struct {
		Error string          `json:"error,omitempty"`
		Stack stack           `json:"stack,omitempty"`
		Cause json.RawMessage `json:"cause,omitempty"`
	}{
		Error: w.msg,
		Stack: w.stack,
		Cause: jsonWrap,
	})
}

func (w withMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.Error(),
		Stack: stack(w.StackTrace()),
	})
}

func (w withMessageAndStack) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Error string `json:"error,omitempty"`
		Stack stack  `json:"stack,omitempty"`
	}{
		Error: w.Error(),
		Stack: w.stack,
	})
}
