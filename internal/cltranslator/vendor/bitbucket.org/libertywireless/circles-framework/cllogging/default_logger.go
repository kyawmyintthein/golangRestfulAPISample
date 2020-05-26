package cllogging

import (
	clerrors "bitbucket.org/libertywireless/circles-framework/clerrors"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	//once          sync.Once
	defaultLogger LoggerInterface
)

// create global default logger
func Init(cfg *LogCfg) {
	newLogger := &logger{
		cfg: cfg,
	}
	newLogger.newLogrus()
	go newLogger.watchLoggerChanges()
	defaultLogger = newLogger
}

// get existing logger or initialed with default config
func GetLogger() *logrus.Logger {
	if defaultLogger == nil {
		Init(&LogCfg{})
	}
	return defaultLogger.GetLogrus()
}

// get existing logger or initialed with default config
func GetStructuredLogger(ctx context.Context) *StructuredLoggerEntry {
	if defaultLogger == nil {
		Init(&LogCfg{})
	}
	return defaultLogger.GetStructuredLogger(ctx)
}

func GetStructuredLoggerIfExist(ctx context.Context) (*StructuredLoggerEntry, bool) {
	isExist := false
	if defaultLogger == nil {
		Init(&LogCfg{})
		isExist = false
	}
	isExist = true
	return defaultLogger.GetStructuredLogger(ctx), isExist
}

func SetReportCaller(reportCaller bool) {
	logrus := GetLogger()
	logrus.SetReportCaller(reportCaller)
}

func NewRequestStructuredLogger() func(next http.Handler) http.Handler {
	logrus := GetLogger()
	return middleware.RequestLogger(&RequestStructureLogger{logrus})
}

func NewGinRequestLogger() gin.HandlerFunc {
	logrus := GetLogger()
	return GinRequestLogger(&RequestStructureLogger{logrus})
}

func NewChiRequestsLogger() func(next http.Handler) http.Handler {
	logrus := GetLogger()
	return middleware.RequestLogger(&RequestStructureLogger{logrus})
}

func SetLogLevel(level string) error {
	if defaultLogger != nil {
		Init(&LogCfg{})
	}
	defaultLogger.SetLogLevel(level)
	return nil
}

// Implementation of ErrorLogger interface{} with clerrors custom error
func WithError(err error) *logrus.Entry {
	logrusLogger := GetLogger()
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

	errWithName, ok := err.(clerrors.ErrorName)
	if ok {
		errTitle = errWithName.Name()
	}

	httpError, ok := err.(clerrors.HttpError)
	if ok {
		statusCode = httpError.StatusCode()
	}

	errorFormatter, ok := err.(clerrors.ErrorFormatter)
	if ok {
		errMsg = errorFormatter.FormattedMessage()
	}

	errorCauser, ok := err.(causer)
	if ok {
		rootCause = errorCauser.Cause()
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

	return logrusLogger.WithFields(fields)
}
