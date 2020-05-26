package cllogging

import (
	"bitbucket.org/libertywireless/circles-framework/clerrors"
	"fmt"
	"math"
	"net/http"
	"time"
	"github.com/sirupsen/logrus"
)

type StructuredLoggerEntry struct {
	logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) AsKVLogger() KVLogger {
	return l
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}){
	l.logger = l.logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.logger.Infoln("request complete")
}

func (l *StructuredLoggerEntry) WriteError(status, bytes int, elapsed time.Duration) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": int(math.Ceil(float64(elapsed.Nanoseconds()) / 1000000.0)),
	})
	l.logger.Errorln("request complete")
}

func (l *StructuredLoggerEntry) WriteWarn(status, bytes int, elapsed time.Duration) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": int(math.Ceil(float64(elapsed.Nanoseconds()) / 1000000.0)),
	})
	l.logger.Warningln("request complete")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.logger = l.logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

func (l *StructuredLoggerEntry) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

// Implementation of ClErrorLogger interface{}
func (l *StructuredLoggerEntry) WithError(err error) *logrus.Entry {
	type causer interface {
		Cause() error
	}

	var stacks interface{}
	var rootCause interface{}
	errCode := 0
	errMsg := err.Error()
	messages := ""
	errTitle := ""
	statusCode := 0
	errStacktrace, ok := err.(clerrors.StackTracer)
	if ok {
		stacks = errStacktrace.GetStackAsJSON()
	}

	errWithCode, ok := err.(clerrors.ErrorCode)
	if ok {
		errCode = errWithCode.Code()
	}

	errorFormatter, ok := err.(clerrors.ErrorFormatter)
	if ok {
		errMsg = errorFormatter.FormattedMessage()
	}

	errorCauser, ok := err.(causer)
	if ok {
		rootCause = errorCauser.Cause()
	}

	errWithName, ok := err.(clerrors.ErrorName)
	if ok {
		errTitle = errWithName.Name()
	}

	httpError, ok := err.(clerrors.HttpError)
	if ok {
		statusCode = httpError.StatusCode()
	}

	messages = clerrors.GetErrorMessages(err)
	fields := logrus.Fields{
		"error_message": errMsg,
		"root_cause":    rootCause,
		"stacktrace":    stacks,
		"messages":      messages,
	}

	if errCode != 0 {
		fields["error_code"] = errCode
	}

	if statusCode != 0 {
		fields["status_code"] = statusCode
	}

	if errTitle != "" {
		fields["error_title"] = errTitle
	}

	return l.logger.WithFields(fields)
}

func (logger *StructuredLoggerEntry) WithFields(fields logrus.Fields) *logrus.Entry {
	return logger.logger.WithFields(fields)
}

func (logger *StructuredLoggerEntry) Fatal(args ...interface{}) {
	logger.logger.Fatal(args...)
}

func (logger *StructuredLoggerEntry) Fatalln(args ...interface{}) {
	logger.logger.Fatalln(args...)
}

func (logger *StructuredLoggerEntry) Fatalf(fmt string, args ...interface{}) {
	logger.logger.Fatalf(fmt, args...)
}

/******************************************* Implementation of KVLogger Interface ************************************************/
func (logger *StructuredLoggerEntry) Error(err error, args ...interface{}) {
	logger.logger.WithError(err).Errorln(args...)
}

func (logger *StructuredLoggerEntry) Errorf(err error, fmt string, args ...interface{}) {
	logger.logger.WithError(err).Errorf(fmt, args...)
}

func (logger *StructuredLoggerEntry) ErrorKV(err error, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).WithError(err).Errorln(args...)
}

func (logger *StructuredLoggerEntry) ErrorKVf(err error, kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).WithError(err).Errorf(fmt, args...)
}

func (logger *StructuredLoggerEntry) Warn(args ...interface{}) {
	logger.logger.Warn(args...)
}

func (logger *StructuredLoggerEntry) Warnln(args ...interface{}) {
	logger.logger.Warnln(args...)
}

func (logger *StructuredLoggerEntry) Warnf(fmt string, args ...interface{}) {
	logger.logger.Warnf(fmt, args)
}

func (logger *StructuredLoggerEntry) WarnKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Warnln(args...)
}

func (logger *StructuredLoggerEntry) WarnKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Warnf(fmt, args...)
}

func (logger *StructuredLoggerEntry) Info(args ...interface{}) {
	logger.logger.Info(args...)
}

func (logger *StructuredLoggerEntry) Infoln(args ...interface{}) {
	logger.logger.Infoln(args...)
}

func (logger *StructuredLoggerEntry) Infof(fmt string, args ...interface{}) {
	logger.logger.Infof(fmt, args...)
}

func (logger *StructuredLoggerEntry) InfoKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Infoln(args...)
}

func (logger *StructuredLoggerEntry) InfoKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Infof(fmt, args...)
}

func (logger *StructuredLoggerEntry) Debug(args ...interface{}) {
	logger.logger.Debug(args...)
}

func (logger *StructuredLoggerEntry) Debugln(args ...interface{}) {
	logger.logger.Debugln(args...)
}

func (logger *StructuredLoggerEntry) Debugf(fmt string, args ...interface{}) {
	logger.logger.Debugf(fmt, args...)
}

func (logger *StructuredLoggerEntry) DebugKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Debugln(args...)
}

func (logger *StructuredLoggerEntry) DebugKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logger.WithFields(fields).Debugf(fmt, args...)
}

func convertKvsToLogrusFields(kv KV) logrus.Fields {
	fields := make(logrus.Fields)
	for k, v := range kv {
		fields[k] = v
	}
	return fields
}
