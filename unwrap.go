package errutil

import (
	"fmt"
	"slices"
	"strings"
)

type causer interface {
	Cause() error
}

type coder interface {
	Code() string
}

type messager interface {
	Message() string
}

type devMessager interface {
	DevMessage() string
	DevMessages() []string
}

type tracer interface {
	Stack() string
	StackTrace() []StackFrame
}

func errorString(err error) string {
	var e, v string

	v = DevMessage(err)
	if v != "" {
		e = v
	}

	v = Message(err)
	if v != "" {
		if e != "" {
			e = fmt.Sprintf("%v (%v)", e, v)
		} else {
			e = fmt.Sprintf("(%v)", v)
		}
	}

	v = Code(err)
	if v != "" {
		e = fmt.Sprintf("[%v] %v", v, e)
	}

	return e
}

func Cause(err error) error {
	if err == nil {
		return nil
	}

	cause, ok := err.(causer)
	if !ok {
		return err
	}

	if cause.Cause() == nil {
		return err
	}

	return Cause(cause.Cause())
}

func Code(err error) string {
	if err == nil {
		return DefaultCode
	}

	e, ok := err.(coder)
	if ok && e.Code() != "" {
		return e.Code()
	}

	cause, ok := err.(causer)
	if !ok {
		return DefaultCode
	}

	return Code(cause.Cause())
}

func Stack(err error) string {
	if err == nil {
		return ""
	}

	trace, ok := err.(tracer)
	if ok {
		return trace.Stack()
	}

	cause, ok := err.(causer)
	if !ok {
		return ""
	}

	return Stack(cause.Cause())
}

func StackTrace(err error) []StackFrame {
	if err == nil {
		return nil
	}

	trace, ok := err.(tracer)
	if ok {
		return trace.StackTrace()
	}

	cause, ok := err.(causer)
	if !ok {
		return nil
	}

	return StackTrace(cause.Cause())
}

func Message(err error, defaultMessage ...string) string {
	var msg string

	if err == nil {
		if len(defaultMessage) > 1 {
			msg = strings.Join(defaultMessage, " ")
		} else if len(defaultMessage) > 0 {
			msg = defaultMessage[0]
		}

		return msg
	}

	cause, ok := err.(causer)
	if ok {
		msg = messageRecursive(cause.Cause())
	}

	e, ok := err.(messager)
	if ok {
		if e.Message() != "" && msg != "" {
			msg = fmt.Sprintf("%s: %s", e.Message(), msg)
		} else {
			msg = e.Message()
		}
	}

	if msg == "" {
		if len(defaultMessage) > 1 {
			msg = strings.Join(defaultMessage, " ")
		} else if len(defaultMessage) > 0 {
			msg = defaultMessage[0]
		}
	}

	return msg
}

func messageRecursive(err error) string {
	if err == nil {
		return ""
	}

	var msg string

	cause, ok := err.(causer)
	if ok {
		msg = messageRecursive(cause.Cause())
	}

	e, ok := err.(messager)
	if ok {
		if e.Message() != "" && msg != "" {
			msg = fmt.Sprintf("%s: %s", e.Message(), msg)
		} else {
			msg = e.Message()
		}
	}

	return msg
}

func Messages(err error) []string {
	msg := make([]string, 0)

	messagesRecursive(err, &msg)

	return msg
}

func messagesRecursive(err error, msg *[]string) {
	if err == nil {
		return
	}

	e, ok := err.(messager)
	if ok {
		if e.Message() != "" {
			*msg = append(*msg, e.Message())
		}
	}

	cause, ok := err.(causer)
	if ok {
		messagesRecursive(cause.Cause(), msg)
	}

	return
}

func DevMessage(err error) string {
	if err == nil {
		return ""
	}

	var msg string

	cause, isCauser := err.(causer)
	if isCauser {
		msg = DevMessage(cause.Cause())
	}

	e, isMessager := err.(devMessager)
	if isMessager {
		if e.DevMessage() != "" && msg != "" {
			msg = fmt.Sprintf("%s, %s", e.DevMessage(), msg)
		} else {
			msg = e.DevMessage()
		}
	}
	if !isMessager && !isCauser {
		return err.Error()
	}

	return msg
}

func DevMessages(err error) []string {
	msg := make([]string, 0)

	devMessagesRecursive(err, &msg)

	return msg
}

func devMessagesRecursive(err error, msg *[]string) {
	if err == nil {
		return
	}

	cause, isCauser := err.(causer)

	e, isMessager := err.(devMessager)
	if isMessager {
		*msg = slices.Concat(*msg, e.DevMessages())
	}
	if !isMessager && !isCauser {
		*msg = append(*msg, err.Error())
	}

	if isCauser {
		devMessagesRecursive(cause.Cause(), msg)
	}

	return
}
