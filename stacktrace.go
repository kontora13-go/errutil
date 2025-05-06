// Copyright 2024-2025 Kontora13. All rights reserved.
// Licensed under the Apache License, Version 2.0

package errutil

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"os"
	"runtime"
	"strings"
)

// MaxStackTraceDepth - максимальная глубина стека
var MaxStackTraceDepth = 50

// StackTraceAppPrefix - префикс модуля, по которому определяется свойство InApp фрейма
var StackTraceAppPrefix = "github.com/kontora13-go"

var goRoot = strings.ReplaceAll(build.Default.GOROOT, "\\", "/")

// StackFrame - структура для описания фрейма стека
type StackFrame struct {
	// Путь к файлу, содержащему этот ProgramCounter
	File string `json:"file,omitempty"`

	// Номер строки в этом файле
	LineNumber int `json:"line_number,omitempty"`

	// Имя функции, в котором произошёл этот вызов
	Function string `json:"function,omitempty"`

	// Package, содержащий эту функцию
	Package string `json:"package,omitempty"`

	InApp bool `json:"in_app,omitempty"`

	PC uintptr `json:"pc,omitempty"`
}

// String возвращает фрейм стека, отформатированный так же, как это делает go в runtime/debug.Stack()
func (frame *StackFrame) String() string {
	if frame.IsEmpty() {
		return ""
	}

	str := fmt.Sprintf("%s:%d (0x%x)\n", frame.File, frame.LineNumber, frame.PC)

	source, err := frame.SourceLine()
	if err != nil {
		return str
	}

	return str + fmt.Sprintf("\t%s: %s\n", frame.Function, source)
}

// IsEmpty возвращает признак заполненности фрейма
func (frame *StackFrame) IsEmpty() bool {
	return frame.PC == 0
}

// SourceLine возвращает строку кода из исходного файла
func (frame *StackFrame) SourceLine() (string, error) {
	if frame.LineNumber <= 0 {
		return "...", nil
	}

	file, err := os.Open(frame.File)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	currentLine := 1
	for scanner.Scan() {
		if currentLine == frame.LineNumber {
			return string(bytes.Trim(scanner.Bytes(), " \t")), nil
		}
		currentLine++
	}
	if err = scanner.Err(); err != nil {
		return "", err
	}

	return "...", nil
}

// NewStackTrace создает стектрейс []StackFrame с использованием runtime.Callers.
func NewStackTrace(skip int) []StackFrame {
	pcs := make([]uintptr, MaxStackTraceDepth)
	n := runtime.Callers(skip+1, pcs)

	if n == 0 {
		return nil
	}

	runtimeFrames := extractFrames(pcs[:n])
	return createFrames(runtimeFrames)
}

// extractFrames выполняет распаковку слайса uintptr в слайс runtime.Frame
func extractFrames(pcs []uintptr) []runtime.Frame {
	n := len(pcs)
	var frames = make([]runtime.Frame, n)
	callersFrames := runtime.CallersFrames(pcs)

	i := n - 1
	for {
		callerFrame, more := callersFrames.Next()
		frames[i] = callerFrame
		i--

		if !more || i == 0 {
			break
		}
	}

	return frames
}

// createFrames создаёт слайс StackFrame, отфильтровывая лишние фреймы
// например, runtime и testing
func createFrames(frames []runtime.Frame) []StackFrame {
	if len(frames) == 0 {
		return nil
	}

	result := make([]StackFrame, 0, len(frames))

	for _, frame := range frames {
		function := frame.Function
		var pkg string
		if function != "" {
			pkg, function = splitPackageAndFunction(function)
		}

		if !shouldSkipFrame(pkg) {
			result = append(result, newFrame(pkg, function, frame.File, frame.Line, frame.PC))
		}
	}

	return result
}

// newFrame создаёт объект фрейма стека
func newFrame(pkg string, function string, file string, line int, pc uintptr) StackFrame {
	frame := StackFrame{
		LineNumber: line,
		Package:    pkg,
		Function:   function,
		File:       file,
		PC:         pc,
	}

	if len(frame.File) == 0 {
		frame.File = "unknown"
	}

	setInAppFrame(&frame)

	return frame
}

// splitPackageAndFunction разделяет имя пакета и имя функции
func splitPackageAndFunction(name string) (pkg string, fun string) {
	pkg = packageName(name)
	if len(pkg) > 0 {
		fun = name[len(pkg)+1:]
	}
	fun = strings.Replace(fun, "·", ".", -1)

	return
}

// packageName выделяет имя пакета, если отсутствует, то возвращается пустая строка.
// Повторяет https://golang.org/pkg/debug/gosym/#Sym.PackageName, избегая зависимости от debug/gosym.
func packageName(name string) string {

	// В версиях Go 1.20 и выше префикс "type:" и "go:" - это сгенерированный
	// компилятором символ, который не относится ни к какому пакету.
	// см, cmd/compile/internal/gc/subr.go
	if strings.HasPrefix(name, "go:") || strings.HasPrefix(name, "type:") {
		return ""
	}

	pathend := strings.LastIndex(name, "/")
	if pathend < 0 {
		pathend = 0
	}

	if i := strings.Index(name[pathend:], "."); i != -1 {
		return name[:pathend+i]
	}
	return ""
}

// shouldSkipFrame проверяет нужно ли пропустить текущий фрейм по имени пакета.
func shouldSkipFrame(pkg string) bool {
	// Пропускаем внутренние пакеты Go
	if pkg == "runtime" || pkg == "testing" {
		return true
	}

	// Пропускаем внутренние фреймы пакета errutil, за исключением _test (для тестирования).
	if strings.HasPrefix(pkg, "github.com/kontora13-go/errutil") &&
		!strings.HasSuffix(pkg, "_test") {
		return true
	}

	return false
}

// setInAppFrame устанавливает признак вызова внутри приложения
func setInAppFrame(frame *StackFrame) {
	if strings.HasPrefix(frame.File, goRoot) ||
		strings.Contains(frame.Package, "vendor") ||
		strings.Contains(frame.Package, "third_party") {
		frame.InApp = false
	} else if strings.HasPrefix(frame.Package, StackTraceAppPrefix) {
		frame.InApp = true
	} else {
		frame.InApp = false
	}
}
