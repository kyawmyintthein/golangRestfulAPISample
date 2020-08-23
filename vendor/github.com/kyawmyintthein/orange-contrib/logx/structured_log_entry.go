package logx

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type discardLoggingKey struct{}
type StructuredLoggerEntry struct {
	logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Error(ctx context.Context, err error, args ...interface{}) {
	l.logger.WithFields(getErrorFields(err)).Errorln(args...)
}

func (l *StructuredLoggerEntry) Warn(ctx context.Context, args ...interface{}) {
	l.logger.Warnln(args...)
}

func (l *StructuredLoggerEntry) Info(ctx context.Context, args ...interface{}) {
	l.logger.Infoln(args...)
}

func (l *StructuredLoggerEntry) Debug(ctx context.Context, args ...interface{}) {
	l.logger.Debugln(args...)
}

func (l *StructuredLoggerEntry) Errorf(ctx context.Context, err error, message string, args ...interface{}) {
	l.logger.WithFields(getErrorFields(err)).Errorf(message, args...)
}

func (l *StructuredLoggerEntry) Warnf(ctx context.Context, message string, args ...interface{}) {
	l.logger.Warnf(message, args...)
}

func (l *StructuredLoggerEntry) Infof(ctx context.Context, message string, args ...interface{}) {
	l.logger.Infof(message, args...)
}

func (l *StructuredLoggerEntry) Debugf(ctx context.Context, message string, args ...interface{}) {
	l.logger.Debugf(message, args...)
}

func (l *StructuredLoggerEntry) ErrorKV(ctx context.Context, err error, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	errorFields := getErrorFields(err)
	for k, v := range errorFields {
		fields[k] = v
	}
	l.logger.WithFields(fields).Errorln(args...)
}

func (l *StructuredLoggerEntry) WarnKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Warnln(args...)
}

func (l *StructuredLoggerEntry) InfoKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Info(args...)
}

func (l *StructuredLoggerEntry) DebugKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Debug(args...)
}

func (l *StructuredLoggerEntry) ErrorKVf(ctx context.Context, err error, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	errorFields := getErrorFields(err)
	for k, v := range errorFields {
		fields[k] = v
	}
	l.logger.WithFields(fields).Errorf(message, args...)
}

func (l *StructuredLoggerEntry) WarnKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Warnf(message, args...)
}

func (l *StructuredLoggerEntry) InfoKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Infof(message, args...)
}

func (l *StructuredLoggerEntry) DebugKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	l.logger.WithFields(fields).Debugf(message, args...)
}

func (l *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, data interface{}) {
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
