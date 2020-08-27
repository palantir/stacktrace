// +build go1.13

package stacktrace_test

import "errors"

func unwrap(err error) error {
	return errors.Unwrap(err)
}
