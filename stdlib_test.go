package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gitlab.com/tozd/go/errors"
)

func TestBase(t *testing.T) {
	t.Parallel()

	parent := errors.Base("Foobar")
	assert.EqualError(t, parent, "Foobar")   //nolint:testifylint
	assert.NoError(t, errors.Unwrap(parent)) //nolint:testifylint
	assert.NotImplements(t, (*stackTracer)(nil), parent)

	another := errors.Basef("Foobar")
	assert.EqualError(t, another, "Foobar")   //nolint:testifylint
	assert.NoError(t, errors.Unwrap(another)) //nolint:testifylint
	assert.NotImplements(t, (*stackTracer)(nil), another)

	child := errors.Basef("Foobar (%w)", parent)
	assert.EqualError(t, child, "Foobar (Foobar)") //nolint:testifylint
	assert.Equal(t, parent, errors.Unwrap(child))
	assert.NotImplements(t, (*stackTracer)(nil), child)
	assert.ErrorIs(t, child, parent) //nolint:testifylint

	child2 := errors.BaseWrap(parent, "Foobar2")
	assert.EqualError(t, child2, "Foobar2") //nolint:testifylint
	assert.Equal(t, parent, errors.Unwrap(child2))
	assert.NotImplements(t, (*stackTracer)(nil), child2)
	assert.ErrorIs(t, child2, parent) //nolint:testifylint

	child3 := errors.BaseWrapf(parent, "Foobar3 (%s)", parent)
	assert.EqualError(t, child3, "Foobar3 (Foobar)") //nolint:testifylint
	assert.Equal(t, parent, errors.Unwrap(child3))
	assert.NotImplements(t, (*stackTracer)(nil), child3)
	assert.ErrorIs(t, child3, parent) //nolint:testifylint

	// We use w() to prevent static analysis.
	child4 := errors.BaseWrapf(parent, "Foobar4 ("+w()+")", parent)
	assert.EqualError(t, child4, "Foobar4 (%!w(*errors.errorString=&{Foobar}))") //nolint:testifylint
	assert.ErrorIs(t, child4, parent)
}
