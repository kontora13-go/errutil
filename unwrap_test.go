package errutil_test

import (
	"encoding/json"
	"fmt"
	"github.com/kontora13-go/errutil"
	"log"
	"testing"
)

func TestUnwrap(t *testing.T) {
	var err error
	var code = "INTERNAL"

	err = fmt.Errorf("first err")
	err = errutil.WithCode(err, code)
	err = newErrorStackWithDepth(err, 4)
	err = errutil.WithDevMessage(err, "test user error", "one")
	err = errutil.WithMessage(err, "user message 1")
	err = errutil.WithMessage(err, "user message 2")
	err = errutil.WithDevMessage(err, "dev 1", "dev 2")

	log.Print("err := ", err.Error())
	log.Print("---")

	log.Print("err.code := ", errutil.Code(err))
	log.Print("err.msg := ", errutil.Message(err))
	log.Print("err.msgs := ", errutil.Messages(err))
	log.Print("---")

	log.Print("err.dev := ", errutil.DevMessage(err))
	log.Print("err.devs := ", errutil.DevMessages(err))
	log.Print("err.cause := ", errutil.Cause(err))
	log.Print("---")

	log.Print("err.Stack := ", errutil.Stack(err))
	st, err := json.Marshal(errutil.StackTrace(err))
	if err != nil {
		t.Error(err)
	}
	log.Print("err.StackTrace := ", string(st))
}

func TestUnwrapMessage(t *testing.T) {
	var err error

	err = fmt.Errorf("first err")
	log.Print("err := ", err.Error())
	log.Print("---")

	log.Print("err.msg := ", errutil.Message(err))
	log.Print("err.msg (default) := ", errutil.Message(err, errutil.DefaultUserMessage))
	log.Print("err.msg (defaults) := ", errutil.Message(err, errutil.DefaultUserMessage, "Или обратитесь к администратору"))
}

func newErrorStackWithDepth(err error, depth int) error {
	if depth <= 0 {
		return errutil.WithStack(err)
	}
	depth--
	return newErrorStackWithDepth(err, depth)
}
