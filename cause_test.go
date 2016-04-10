package stacktrace_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/stacktrace"
)

type customError string

func (e customError) Error() string { return string(e) }

func TestRootCause(t *testing.T) {
	for _, test := range []struct {
		err       error
		rootCause error
	}{
		{
			err:       nil,
			rootCause: nil,
		},
		{
			err:       errors.New("msg"),
			rootCause: errors.New("msg"),
		},
		{
			err:       stacktrace.NewError("msg"),
			rootCause: errors.New("msg"),
		},
		{
			err:       stacktrace.Propagate(stacktrace.NewError("msg1"), "msg2"),
			rootCause: errors.New("msg1"),
		},
		{
			err:       customError("msg"),
			rootCause: customError("msg"),
		},
		{
			err:       stacktrace.Propagate(customError("msg1"), "msg2"),
			rootCause: customError("msg1"),
		},
	} {
		assert.Equal(t, test.rootCause, stacktrace.RootCause(test.err))
	}
}
