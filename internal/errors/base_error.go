package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors/interfaces"
	"path"
	"runtime"
	"strings"
	"sync"
)

// default implementation of error interfaces
type BaseError struct {
	// error id
	id string

	// error code
	code int

	// message holds the text of the error message. It may be empty
	// or it may be format string
	messageFormat string

	// cause holds the cause of the error as returned
	// by the Cause method.
	cause error

	// arguments for formatted string
	args []interface{}

	stack       []uintptr
	framesOnce  sync.Once
	stackFrames []interfaces.StackFrame
}

func NewError(code int, id string, args ...interface{}) *BaseError {
	stack := make([]uintptr, 2)
	stackLength := runtime.Callers(3, stack)
	err := &BaseError{
		id:    id,
		code:  code,
		cause: nil,
		args:  args,
		stack: stack[:stackLength],
	}
	return err
}

func WithMessage(code int, id string, messageTemplate string, args ...interface{}) *BaseError {
	stack := make([]uintptr, 2)
	stackLength := runtime.Callers(3, stack)
	err := &BaseError{
		id:            id,
		code:          code,
		cause:         nil,
		messageFormat: messageTemplate,
		args:          args,
		stack:         stack[:stackLength],
	}
	return err
}

func (e *BaseError) Wrap(err error) error {
	e.cause = err
	return e
}

func (e *BaseError) Error() string {
	return e.FormattedMessage()
}

func (e *BaseError) ID() string {
	return e.id
}

func (e *BaseError) Code() int {
	return e.code
}

func (e *BaseError) Message() string {
	return e.messageFormat
}

func (e *BaseError) GetArgs() []interface{} {
	return e.args
}

// Return nested error
func (e *BaseError) GetMessage() string {
	return e.messageFormat
}

func (e *BaseError) FormattedMessage() string {
	if e.messageFormat != "" {
		argsMap := make(map[string]string)
		msg := e.messageFormat
		if len(e.args) != 0 {
			previousKey := ""
			for _, v := range e.args {
				if previousKey != "" {
					argsMap[previousKey] = v.(string)
				}
				previousKey = v.(string)
			}
		}
		for k, v := range argsMap {
			msg = strings.Replace(msg, fmt.Sprintf("{{var_%s}}", k), v, -1)
		}
		return msg
	} else if e.id != "" {
		argsMap := make(map[string]string)
		if len(e.args) != 0 {
			previousKey := ""
			for _, v := range e.args {
				if previousKey != "" {
					argsMap[previousKey] = v.(string)
				}
				previousKey = v.(string)
			}
		}
		var buf bytes.Buffer
		for k, v := range argsMap {
			buf.WriteString(fmt.Sprintf("%s:%v, ", k, v))
		}

		if len(argsMap) != 0{
			buf.Truncate(buf.Len() - 2)
			return fmt.Sprintf("%s : [%s]", e.id, buf.String())
		}
		return fmt.Sprintf("%s", e.id)
	}
	return "NoErrorMessage"
}


func (w *BaseError) Cause() error { return w.cause }

func (e *BaseError) StackAddrs() string {
	buf := bytes.NewBuffer(make([]byte, 0, len(e.stack)*8))
	for _, pc := range e.stack {
		fmt.Fprintf(buf, "0x%x ", pc)
	}
	bufBytes := buf.Bytes()
	return string(bufBytes[:len(bufBytes)-1])
}

func (e *BaseError) StackFrames() []interfaces.StackFrame {
	e.framesOnce.Do(func() {
		e.stackFrames = make([]interfaces.StackFrame, len(e.stack))
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

func (e *BaseError) GetStack() string {
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

func (e *BaseError) GetStackAsJSON() interface{} {
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
