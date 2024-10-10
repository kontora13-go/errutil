package errutil

import (
	"fmt"
	"strings"
)

func WithCode(err error, code string) error {
	if err == nil {
		return &errWithStack{
			cause:      err,
			code:       code,
			stacktrace: newErrorStack(),
		}
	}

	return &errWithCode{
		cause: err,
		code:  code,
	}
}

func WithStack(err error) error {
	return &errWithStack{
		cause:      err,
		stacktrace: newErrorStack(),
	}
}

func WithMessage(err error, msg ...string) error {
	if err == nil {
		err = &errWithStack{
			code:       DefaultCode,
			cause:      err,
			stacktrace: newErrorStack(),
		}
	}

	return &errWithMessage{
		cause: err,
		msg:   strings.Join(msg, ": "),
	}
}

func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		err = &errWithStack{
			code:       DefaultCode,
			cause:      err,
			stacktrace: newErrorStack(),
		}
	}

	return &errWithMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

func WithDevMessage(err error, msg ...string) error {
	if err == nil {
		err = &errWithStack{
			code:       DefaultCode,
			cause:      err,
			stacktrace: newErrorStack(),
		}
	}

	return &errWithDevMessage{
		cause: err,
		dev:   msg,
	}
}

func WithDevMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		err = &errWithStack{
			code:       DefaultCode,
			cause:      err,
			stacktrace: newErrorStack(),
		}
	}

	return &errWithDevMessage{
		cause: err,
		dev:   []string{fmt.Sprintf(format, args...)},
	}
}
