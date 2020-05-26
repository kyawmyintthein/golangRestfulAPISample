package interfaces

import (
	"runtime"
)

type StackTracer interface {
	Error() string
	StackAddrs() string
	StackFrames() []StackFrame
	GetStack() string
	GetStackAsJSON() interface{}
}

// Represents a single stack frame.
type StackFrame struct {
	PC         uintptr
	Func       *runtime.Func
	FuncName   string
	File       string
	LineNumber int
}
