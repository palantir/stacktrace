// +build !go1.13

package stacktrace_test

func unwrap(err error) error {
	t, ok := err.(interface{ Unwrap() error })
	if !ok {
		return nil
	}
	return t.Unwrap()
}
