package errors

import (
	"bytes"
	"encoding/json"
	"reflect"
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
		data["stack"] = StackFormatter{st}
	} else {
		placeholderErr, ok := err.(placeholderStackTracer) //nolint:errorlint
		if ok {
			placeholderSt := placeholderErr.StackTrace()
			if len(placeholderSt) > 0 {
				data["stack"] = placeholderSt
			}
		}
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

func hasJSONTag(typ reflect.Type) bool {
	if typ.Kind() == reflect.Struct {
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if field.Tag.Get("json") != "" {
				return true
			}
			if field.Anonymous && hasJSONTag(field.Type) {
				return true
			}
		}
	}

	return false
}

// Does the error not implement our interfaces but implement MarshalJSON or uses any JSON struct tags?
func useMarshaler(err error) bool {
	if useKnownInterface(err) {
		return false
	}

	// We check for this interface without unwrapping because it does not work with wrapping anyway.
	_, ok := err.(json.Marshaler) //nolint:errorlint
	if ok {
		return true
	}

	typ := reflect.TypeOf(err).Elem()
	return hasJSONTag(typ)
}

// marshalJSONAnyError marshals our and foreign errors.
func marshalJSONAnyError(err error) ([]byte, E) {
	if err == nil {
		return []byte("null"), nil
	}

	// This short-circuits our errors as well to directly call marshalJSONError
	// and do not call it indirectly through marshalWithoutEscapeHTML.
	if !useMarshaler(err) {
		return marshalJSONError(err)
	}

	// Does the error marshal to something useful?
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

// MarshalJSON marshals the error as JSON according to the json.Marshaler interface.
//
// The error does not have to necessary come from this package and it will be marshaled
// in the same way if it implements interfaces used by this package (e.g., stackTracer
// or detailer interfaces). Only if those interfaces are not implemented,
// but json.Marshaler interface is or the error is a struct with JSON struct tags,
// marshaling will be delegated to the error itself.
//
// Errors which do come from this package can be directly marshaled in the same way as
// this function does as they implement json.Marshaler interface.
// If you are not sure about the source of the error, it is safe to call this function
// on them as well.
func (f Formatter) MarshalJSON() ([]byte, error) {
	return marshalJSONAnyError(f.Error)
}
