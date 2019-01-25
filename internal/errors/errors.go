package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"runtime"
	"sync"
)

// Expose additional information for error.
type CustomError interface {
	Error() string
	GetCode() uint32
	GetInner() error
	StackAddrs() string
	StackFrames() []StackFrame
	GetStack() string
	GetStackAsJSON() interface{}
	GetFullMessage() string
	GetMessages() string
}

// Represents a single stack frame.
type StackFrame struct {
	PC         uintptr
	Func       *runtime.Func
	FuncName   string
	File       string
	LineNumber int
}

type baseError struct {
	code    uint32
	message string
	inner   error

	stack       []uintptr
	framesOnce  sync.Once
	stackFrames []StackFrame
}

func (e *baseError) Error() string {
	return e.message
}


func (e *baseError) GetCode() uint32 {
	return e.code
}

// Return nested error
func (e *baseError) GetInner() error {
	return e.inner
}

// Return errors message and stacktrace as string
func (e *baseError) GetFullMessage() string {
	return extractFullErrorMessage(e, true)
}

func (e *baseError) GetMessages() string {
	return extractFullErrorMessage(e, false)
}

func (e *baseError) StackAddrs() string {
	buf := bytes.NewBuffer(make([]byte, 0, len(e.stack)*8))
	for _, pc := range e.stack {
		fmt.Fprintf(buf, "0x%x ", pc)
	}
	bufBytes := buf.Bytes()
	return string(bufBytes[:len(bufBytes)-1])
}

func (e *baseError) StackFrames() []StackFrame {
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

func (e *baseError) GetStack() string {
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

func (e *baseError) GetStackAsJSON() interface{} {
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

func extractFullErrorMessage(e CustomError, includeStack bool) string {
	var ok bool
	var lastClErr CustomError
	errMsg := bytes.NewBuffer(make([]byte, 0, 1024))

	dbxErr := e
	for {
		lastClErr = dbxErr
		errMsg.WriteString(dbxErr.Error())

		innerErr := dbxErr.GetInner()
		if innerErr == nil {
			break
		}
		dbxErr, ok = innerErr.(CustomError)
		if !ok {
			// We have reached the end and traveresed all inner errors.
			// Add last message and exit loop.
			errMsg.WriteString(" [")
			errMsg.WriteString(innerErr.Error())
			errMsg.WriteString("] ")
			break
		}
		errMsg.WriteString("\n")
	}
	if includeStack {

		errMsg.WriteString("\n")
		errMsg.WriteString("\n STACK TRACE:\n")
		errMsg.WriteString(lastClErr.GetStack())
	}
	return errMsg.String()
}

func New(code uint32, message string) error {
	return new(nil, code, message)
}

func Newf(code uint32, format string, args ...interface{}) error {
	return new(nil, code, fmt.Sprintf(format, args...))
}

func new(err error, code uint32, message string) *baseError {
	stack := make([]uintptr, 2)
	stackLength := runtime.Callers(3, stack)
	return &baseError{
		message: message,
		code:    code,
		stack:   stack[:stackLength],
		inner:   err,
	}
}

// wrap error with custom error
func Wrap(err error, code uint32, message string) error {
	return new(err, code, message)
}

func Wrapf(err error, code uint32, format string, args ...interface{}) error {
	return new(err, code, fmt.Sprintf(format, args...))
}

// get the root cause error
func RootError(ierr error) (nerr error) {
	nerr = ierr
	for i := 0; i < 20; i++ {
		terr := unwrapError(nerr)
		if terr == nil {
			return nerr
		}
		nerr = terr
	}
	return fmt.Errorf("too many iterations: %T", nerr)
}

func unwrapError(ierr error) (nerr error) {
	if clError, ok := ierr.(CustomError); ok {
		return clError.GetInner()
	}
	return nil
}
