package errors

import "encoding/json"

func (e fundamental) MarshalJSON() ([]byte, error) {
	return marshalWithoutEscapeHTML(&struct {
		Error string         `json:"error,omitempty"`
		Stack StackFormatter `json:"stack,omitempty"`
	}{
		Error: e.msg,
		Stack: e.stack,
	})
}

func (e msgWithStack) MarshalJSON() ([]byte, error) {
	return marshalWithoutEscapeHTML(&struct {
		Error string         `json:"error,omitempty"`
		Stack StackFormatter `json:"stack,omitempty"`
	}{
		Error: e.msg,
		Stack: e.StackTrace(),
	})
}

func (e msgWithoutStack) MarshalJSON() ([]byte, error) {
	return marshalWithoutEscapeHTML(&struct {
		Error string         `json:"error,omitempty"`
		Stack StackFormatter `json:"stack,omitempty"`
	}{
		Error: e.msg,
		Stack: e.stack,
	})
}

func (e msgJoined) MarshalJSON() ([]byte, error) {
	jsonWraps := make([]json.RawMessage, 0, len(e.errs))
	for _, err := range e.errs {
		jsonWrap, e := marshalJSONAnyError(err)
		if e != nil {
			return nil, e
		}
		jsonWraps = append(jsonWraps, jsonWrap)
	}
	return marshalWithoutEscapeHTML(&struct {
		Error  string            `json:"error,omitempty"`
		Errors []json.RawMessage `json:"errors,omitempty"`
		Stack  StackFormatter    `json:"stack,omitempty"`
	}{
		Error:  e.msg,
		Errors: jsonWraps,
		Stack:  e.stack,
	})
}

func (e withStack) MarshalJSON() ([]byte, error) {
	return marshalWithoutEscapeHTML(&struct {
		Error string         `json:"error,omitempty"`
		Stack StackFormatter `json:"stack,omitempty"`
	}{
		Error: e.Error(),
		Stack: e.StackTrace(),
	})
}

func (e withoutStack) MarshalJSON() ([]byte, error) {
	return marshalWithoutEscapeHTML(&struct {
		Error string         `json:"error,omitempty"`
		Stack StackFormatter `json:"stack,omitempty"`
	}{
		Error: e.Error(),
		Stack: e.stack,
	})
}

func (e cause) MarshalJSON() ([]byte, error) {
	jsonWrap, err := marshalJSONAnyError(e.err)
	if err != nil {
		return nil, err
	}
	return marshalWithoutEscapeHTML(&struct {
		Error string          `json:"error,omitempty"`
		Stack StackFormatter  `json:"stack,omitempty"`
		Cause json.RawMessage `json:"cause,omitempty"`
	}{
		Error: e.msg,
		Stack: e.stack,
		Cause: jsonWrap,
	})
}
