package errorx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sync"
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

type ErrorStacktrace struct {
	stack       []uintptr
	framesOnce  sync.Once
	stackFrames []StackFrame
}

func NewErrorWithStackTrace(stackLen int, callerLen int) *ErrorStacktrace {
	stack := make([]uintptr, stackLen)
	stackLength := runtime.Callers(callerLen, stack)
	return &ErrorStacktrace{
		stack: stack[:stackLength],
	}
}

func (e *ErrorStacktrace) StackAddrs() string {
	buf := bytes.NewBuffer(make([]byte, 0, len(e.stack)*8))
	for _, pc := range e.stack {
		fmt.Fprintf(buf, "0x%x ", pc)
	}
	bufBytes := buf.Bytes()
	return string(bufBytes[:len(bufBytes)-1])
}

func (e *ErrorStacktrace) StackFrames() []StackFrame {
	e.framesOnce.Do(func() {
		e.stackFrames = make([]StackFrame, len(e.stack))
		for i, pc := range e.stack {
			frame := &e.stackFrames[i]
			frame.PC = pc
			frame.Func = runtime.FuncForPC(pc)
			if frame.Func != nil {
				frame.FuncName = frame.Func.Name()
				frame.File, frame.LineNumber = frame.Func.FileLine(frame.PC - 1)
			}
		}
	})
	return e.stackFrames
}

func (e *ErrorStacktrace) GetStack() string {
	stackFrames := e.StackFrames()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	for _, frame := range stackFrames {
		_, _ = buf.WriteString(frame.FuncName)
		_, _ = buf.WriteString("\n")
		fmt.Fprintf(buf, "\t%s:%d +0x%x\n",
			frame.File, frame.LineNumber, frame.PC)
	}
	return buf.String()
}

func (e *ErrorStacktrace) GetStackAsJSON() interface{} {
	stackFrames := e.StackFrames()
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	var (
		data []byte
		i    interface{}
	)
	data = append(data, '[')
	for i, frame := range stackFrames {
		if i != 0 {
			data = append(data, ',')
		}
		name := path.Base(frame.FuncName)
		frameBytes := []byte(fmt.Sprintf(`{"filepath": "%s", "name": "%s", "line": %d}`, frame.File, name, frame.LineNumber))
		data = append(data, frameBytes...)
	}
	data = append(data, ']')
	buf.Write(data)
	_ = json.Unmarshal(data, &i)
	return i
}
