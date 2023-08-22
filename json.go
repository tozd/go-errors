package errors

import (
	"bytes"
	"encoding/json"
)

func marshalWithoutEscapeHTML(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	b := buf.Bytes()
	if len(b) > 0 {
		// Remove trailing \n which is added by Encode.
		return b[:len(b)-1], nil
	}
	return b, nil
}

// marshalJSONError marshals errors using interfaces.
func marshalJSONError(err error) ([]byte, E) {
	details, cause, errs := allDetailsUntilCauseOrJoined(err)

	data := map[string]interface{}{}

	// We start with details so that other "standard"
	// fields can override conflicting fields from details.
	for key, value := range details {
		data[key] = value
	}

	msg := err.Error()
	if msg != "" {
		data["error"] = msg
	}

	st := getExistingStackTrace(err)
	if len(st) > 0 {
		data["stack"] = StackFormatter(st)
	}

	for _, er := range errs {
		// er should never be nil, but we still check.
		if er != nil {
			jsonEr, e := marshalJSONAnyError(er)
			if e != nil {
				return nil, e
			}
			if len(jsonEr) != 0 && !bytes.Equal(jsonEr, []byte("{}")) {
				if data["errors"] == nil {
					data["errors"] = []json.RawMessage{json.RawMessage(jsonEr)}
				} else {
					data["errors"] = append(data["errors"].([]json.RawMessage), json.RawMessage(jsonEr)) //nolint:forcetypeassert
				}
			}
		}
	}

	if cause != nil {
		jsonCause, e := marshalJSONAnyError(cause)
		if e != nil {
			return nil, e
		}
		if len(jsonCause) != 0 && !bytes.Equal(jsonCause, []byte("{}")) {
			data["cause"] = json.RawMessage(jsonCause)
		}
	}

	jsonErr, e := marshalWithoutEscapeHTML(data)
	if e != nil {
		return nil, WithStack(e)
	}
	return jsonErr, nil
}

// marshalJSONAnyError marshals our and foreign errors.
func marshalJSONAnyError(err error) ([]byte, E) {
	if err == nil {
		return []byte("null"), nil
	}

	// We short-circuit our errors to directly call marshalJSONError
	// and do not call it indirectly through marshalWithoutEscapeHTML.
	switch err.(type) { //nolint:errorlint
	case *fundamental, *msgWithStack, *msgWithoutStack, *msgJoined, *withStack, *withoutStack, *cause:
		return marshalJSONError(err)
	}

	// Does the error marshal to something useful on its own?
	// We do not call MarshalJSON here but invoke Go JSON marshal because
	// maybe error relays on the default JSON marshal (of structs, for example).
	jsonErr, e := marshalWithoutEscapeHTML(err)
	if e != nil {
		return nil, WithStack(e)
	}
	if len(jsonErr) == 0 || bytes.Equal(jsonErr, []byte("{}")) {
		// No it does not, we call marshalJSONError.
		return marshalJSONError(err)
	}

	// It does, we return it.
	return jsonErr, nil
}

func (f Formatter) MarshalJSON() ([]byte, error) {
	return marshalJSONAnyError(f.Error)
}
