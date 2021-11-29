package errors

import (
	"bytes"
	"encoding/json"
)

// marshalJSONError marshals foreign errors.
func marshalJSONError(err error) ([]byte, E) {
	var e error
	var s stack
	if stackErr, ok := err.(stackTracer); ok {
		s = stack(stackErr.StackTrace())
	}
	var jsonWrap []byte
	u, ok := err.(interface {
		Unwrap() error
	})
	if ok {
		unwrap := u.Unwrap()
		jsonWrap, e = json.Marshal(unwrap)
		if e != nil {
			return nil, WithStack(e)
		}
		if len(jsonWrap) == 0 || bytes.Equal(jsonWrap, []byte("{}")) {
			var eStack E
			jsonWrap, eStack = marshalJSONError(unwrap)
			if eStack != nil {
				return nil, eStack
			}
		}
		if bytes.Equal(jsonWrap, []byte("{}")) {
			jsonWrap = []byte{}
		}
	}
	jsonErr, e := json.Marshal(&struct {
		Error string          `json:"error,omitempty"`
		Stack stack           `json:"stack,omitempty"`
		Wrap  json.RawMessage `json:"wrap,omitempty"`
	}{
		Error: err.Error(),
		Stack: s,
		Wrap:  jsonWrap,
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
		Wrap  json.RawMessage `json:"wrap,omitempty"`
	}{
		Error: w.msg,
		Stack: w.stack,
		Wrap:  jsonWrap,
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
