package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestBase(t *testing.T) {
	t.Parallel()

	parent := errors.Base("Foobar")
	assert.EqualError(t, parent, "Foobar")
	assert.Nil(t, errors.Unwrap(parent))
	assert.NotImplements(t, (*stackTracer)(nil), parent)

	another := errors.Basef("Foobar")
	assert.EqualError(t, another, "Foobar")
	assert.Nil(t, errors.Unwrap(another))
	assert.NotImplements(t, (*stackTracer)(nil), another)

	child := errors.Basef("Foobar (%w)", parent)
	assert.EqualError(t, child, "Foobar (Foobar)")
	assert.Equal(t, parent, errors.Unwrap(child))
	assert.NotImplements(t, (*stackTracer)(nil), child)
	assert.ErrorIs(t, child, parent)

	child2 := errors.BaseWrap(parent, "Foobar2")
	assert.EqualError(t, child2, "Foobar2")
	assert.Equal(t, parent, errors.Unwrap(child2))
	assert.NotImplements(t, (*stackTracer)(nil), child2)
	assert.ErrorIs(t, child2, parent)

	child3 := errors.BaseWrapf(parent, "Foobar3 (%s)", parent)
	assert.EqualError(t, child3, "Foobar3 (Foobar)")
	assert.Equal(t, parent, errors.Unwrap(child3))
	assert.NotImplements(t, (*stackTracer)(nil), child3)
	assert.ErrorIs(t, child3, parent)

	// We use w() to prevent static analysis.
	child4 := errors.BaseWrapf(parent, "Foobar4 ("+w()+")", parent)
	assert.EqualError(t, child4, "Foobar4 (%!w(*errors.errorString=&{Foobar}))")
	assert.ErrorIs(t, child4, parent)
}
