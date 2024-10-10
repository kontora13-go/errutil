// Copyright 2024 Kontora13. All rights reserved.
// Licensed under the Apache License, Version 2.0

// Описания различных врапперов ошибок, используемых
// в пакете для работы с ошибками 

package errutil

import (
	"bytes"
	"strings"
)

// errWithCode - ошибка с контекстом кода ошибки
type errWithCode struct {
	code  string
	cause error
}

// Error - получение текстового представления ошибки
func (e *errWithCode) Error() string {
	return errorString(e)
}

// Cause - распаковка исходной ошибки
func (e *errWithCode) Cause() error {
	return e.cause
}

// Code - получение кода ошибки
func (e *errWithCode) Code() string {
	return e.code
}

/*
----------
*/

// errWithStack - ошибка с Callers trace ошибки
type errWithStack struct {
	code       string
	stacktrace []StackFrame
	cause      error
}

// Error - получение текстового представления ошибки
func (e *errWithStack) Error() string {
	return errorString(e)
}

// Stack - получение Callers trace ошибки
func (e *errWithStack) Stack() string {
	buf := bytes.Buffer{}

	for _, frame := range e.stackFrames() {
		buf.WriteString(frame.String())
	}

	return string(buf.Bytes())
}

func (e *errWithStack) StackTrace() []StackFrame {
	return e.stacktrace
}

// Cause - распаковка исходной ошибки
func (e *errWithStack) Cause() error {
	return e.cause
}

// Code - получение кода ошибки
func (e *errWithStack) Code() string {
	return e.code
}

// stackFrames - возвращает массив фреймов, содержащих информацию о стеке.
func (e *errWithStack) stackFrames() []StackFrame {
	return e.stacktrace
}

func newErrorStack() []StackFrame {
	return NewStackTrace(3)
}

/*
----------
*/

// errWithMessage - ошибка, содержащая сообщение для пользователя
type errWithMessage struct {
	msg   string
	cause error
}

// Message - получение сообщения, содержащегося в ошибке
func (e *errWithMessage) Message() string {
	return e.msg
}

// Cause - распаковка исходной ошибки
func (e *errWithMessage) Cause() error {
	return e.cause
}

// Error - получение текстового представления ошибки
func (e *errWithMessage) Error() string {
	return errorString(e)
}

/*
----------
*/

// errWithDevMessage - ошибка, содержащая сообщения для разработчика
type errWithDevMessage struct {
	dev   []string
	cause error
}

// DevMessage - получение dev-сообщения, содержащихся в ошибке
func (e *errWithDevMessage) DevMessage() string {
	return strings.Join(e.dev, ": ")
}

// DevMessages - получение dev-сообщения, содержащихся в ошибке в виде слайса
func (e *errWithDevMessage) DevMessages() []string {
	return e.dev
}

// Cause - распаковка исходной ошибки
func (e *errWithDevMessage) Cause() error {
	return e.cause
}

// Error - получение текстового представления ошибки
func (e *errWithDevMessage) Error() string {
	return errorString(e)
}
