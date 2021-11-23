package errors_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

// See: https://github.com/stretchr/testify/issues/1065
func notImplements(t *testing.T, interfaceObject interface{}, object interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()
	interfaceType := reflect.TypeOf(interfaceObject).Elem()

	if object == nil {
		return assert.Fail(t, fmt.Sprintf("Cannot check if nil does not implement %v", interfaceType), msgAndArgs...)
	}
	if !reflect.TypeOf(object).Implements(interfaceType) {
		return true
	}

	return assert.Fail(t, fmt.Sprintf("%T must not implement %v", object, interfaceType), msgAndArgs...)
}

func TestBase(t *testing.T) {
	parent := errors.Base("Foobar")
	assert.EqualError(t, parent, "Foobar")
	assert.Nil(t, errors.Unwrap(parent))
	notImplements(t, (*stackTracer)(nil), parent)

	another := errors.Basef("Foobar")
	assert.EqualError(t, another, "Foobar")
	assert.Nil(t, errors.Unwrap(another))
	notImplements(t, (*stackTracer)(nil), another)

	child := errors.Basef("Foobar (%w)", parent)
	assert.EqualError(t, child, "Foobar (Foobar)")
	assert.Equal(t, parent, errors.Unwrap(child))
	notImplements(t, (*stackTracer)(nil), child)
	assert.ErrorIs(t, child, parent)

	child2 := errors.BaseWrap(parent, "Foobar2")
	assert.EqualError(t, child2, "Foobar2")
	assert.Equal(t, parent, errors.Unwrap(child2))
	notImplements(t, (*stackTracer)(nil), child2)
	assert.ErrorIs(t, child2, parent)

	child3 := errors.BaseWrapf(parent, "Foobar3 (%s)", parent)
	assert.EqualError(t, child3, "Foobar3 (Foobar)")
	assert.Equal(t, parent, errors.Unwrap(child3))
	notImplements(t, (*stackTracer)(nil), child3)
	assert.ErrorIs(t, child3, parent)

	// We use w() to prevent static analysis.
	child4 := errors.BaseWrapf(parent, "Foobar4 ("+w()+")", parent)
	assert.EqualError(t, child4, "Foobar4 (%!w(*errors.errorString=&{Foobar}))")
	assert.ErrorIs(t, child4, parent)
}
