package logging

import (
	"fmt"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/errors"
	"github.com/sirupsen/logrus"
	"time"
)

type StructuredLoggerEntry struct {
	logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// WithError custom function
func (l *StructuredLoggerEntry) WithError(err error) *logrus.Entry {
	var (
		errMsg     string
		errCode    uint32
		stacktrace interface{}
		causeError string
	)
	clerror, ok := err.(errors.CustomError)
	if ok {
		errCode = clerror.GetCode()
		errMsg = clerror.Error()
		stacktrace = clerror.GetStackAsJSON()
		causeError = clerror.GetMessages()
	} else {
		errCode = 0
		errMsg = err.Error()
		stacktrace = ""
		causeError = err.Error()
	}
	return l.logger.WithFields(logrus.Fields{
		"error_message": errMsg,
		"error_code":    errCode,
		"caurse": 		 causeError,
		"stacktrace":         stacktrace,
	})
}

func (l *StructuredLoggerEntry) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

func (l *StructuredLoggerEntry) Info(args ...interface{}) {
	l.logger.Info(args)
}

func (l *StructuredLoggerEntry) Infoln(args ...interface{}) {
	l.logger.Infoln(args)
}

func (l *StructuredLoggerEntry) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args)
}

func (l *StructuredLoggerEntry) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args)
}


