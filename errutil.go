// Copyright 2024-2025 Kontora13. All rights reserved.
// Licensed under the Apache License, Version 2.0

// Библиотека для создания собственного типа реализующего интерфейс error.
// Позволят создавать, редактировать, объединять собственные ошибки.

package errutil

import "fmt"

var (
	DefaultCode        = CodeCritical
	DefaultUserMessage = "Упс, что-то пошло не так. Попробуйте позже..."
)

const (
	CodePanic    string = "PANIC"
	CodeCritical string = "CRITICAL"
	CodeUser     string = "USER"
)

// New - конструктор ошибки из списка строк
func New(message ...string) error {
	err := &errWithStack{
		code:       DefaultCode,
		stacktrace: newErrorStack(),
	}

	return &errWithDevMessage{
		cause: err,
		dev:   message,
	}
}

// Newf - конструктор ошибки из форматной строки с параметрами
func Newf(format string, args ...interface{}) error {
	err := &errWithStack{
		code:       DefaultCode,
		stacktrace: newErrorStack(),
	}

	return &errWithDevMessage{
		cause: err,
		dev:   []string{fmt.Sprintf(format, args...)},
	}
}

// NewWithCode - конструктор ошибки из списка строк с указанием кода ошибки
func NewWithCode(code string, message ...string) error {
	err := &errWithStack{
		code:       code,
		stacktrace: newErrorStack(),
	}

	return &errWithDevMessage{
		cause: err,
		dev:   message,
	}
}

// NewWithCodef - конструктор ошибки из форматной строки с параметрами с указанием кода ошибки
func NewWithCodef(code string, format string, args ...interface{}) error {
	err := &errWithStack{
		code:       code,
		stacktrace: newErrorStack(),
	}

	return &errWithDevMessage{
		cause: err,
		dev:   []string{fmt.Sprintf(format, args...)},
	}
}
