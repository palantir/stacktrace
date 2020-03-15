package stacktrace

import (
	"errors"
)

/*
RootCause unwraps the original error that caused the current one.

	_, err := f()
	if perr, ok := stacktrace.RootCause(err).(*ParsingError); ok {
		showError(perr.Line, perr.Column, perr.Text)
	}
*/
func RootCause(err error) error {
	for {
		st, ok := err.(*stacktrace)
		if !ok {
			return err
		}
		if st.cause == nil {
			return errors.New(st.message)
		}
		err = st.cause
	}
}

/*
RootMessage returns message which was associated with root cause error.
	func f(file string) error {
		....
		....
		file, err := os.Open(file)
		if err != nil {
			return stacktrace.Propagate(err, "formatted message for user", ...specific args)
		}
		....
		....
	}

	err := f("./foo/bar")
	if err != nil {
		err := stacktrace.RootCause(err)
		// do something with root error

		// return formatted message for user
		showToUser(stacktrace.RootMessage(err))
	}

*/
func RootMessage(err error) string {
	var message string
	for {
		st, ok := err.(*stacktrace)
		if !ok {
			return message
		}
		if st.cause == nil {
			return st.message
		}
		err = st.cause
		message = st.message
	}
}
