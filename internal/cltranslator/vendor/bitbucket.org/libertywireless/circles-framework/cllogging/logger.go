package cllogging

import (
	"bitbucket.org/libertywireless/circles-framework/clerrors"
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

/*
 * default logger implementation
 */
type logger struct {
	cfg     *LogCfg
	logrus  *logrus.Logger
	logfile *os.File
}

// create logger object
func New(cfg *LogCfg) LoggerInterface {
	logger := &logger{
		cfg: cfg,
	}

	logger.newLogrus()

	go logger.watchLoggerChanges()

	return logger
}

func (logger *logger) GetLogrus() *logrus.Logger {
	return logger.logrus
}

func (logger *logger) SetReportCaller(reportCaller bool) {
	logger.logrus.SetReportCaller(reportCaller)
}

func (logger *logger) NewRequestStructuredLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&RequestStructureLogger{logger.logrus})
}

// create new logrus object
func (logger *logger) newLogrus() {
	logger.logrus = &logrus.Logger{
		Hooks: make(logrus.LevelHooks),
	}

	logLevel, err := logrus.ParseLevel(logger.cfg.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.logrus.Level = logLevel
	logger.logrus.SetReportCaller(logger.cfg.CallerInfo)

	if logger.cfg.JsonLogFormat {
		logger.logrus.Formatter = new(JSONFormatter)
	} else {
		logger.logrus.Formatter = new(logrus.TextFormatter)
	}

	if logger.cfg.LogFilePath == "" {
		logger.logrus.Out = os.Stdout
		logger.logrus.Errorln("Empty log file. Set 'Stdout' as default")
		logger.logrus.Info("Initialized Logger successfully")
		return
	}

	logfile, err := os.OpenFile(logger.cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.logrus.Errorln("Failed to set log file : '%v'. Set 'Stdout' as default", err)
		return
	}

	logger.logfile = logfile
	logger.logrus.Out = logger.logfile

	logger.logrus.Info("Initialized Logger successfully")
}

func (logger *logger) SetLogLevel(level string) error {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
		return err
	}
	logger.logrus.Level = logLevel
	return nil
}

func (logger *logger) GetStructuredLogger(ctx context.Context) *StructuredLoggerEntry {
	log := ctx.Value(middleware.LogEntryCtxKey)
	if log == nil {
		log = ctx.Value(middleware.LogEntryCtxKey.String())
	}

	isDiscard, _ := ctx.Value(DiscardLoggingKey{}).(bool)
	if log != nil {
		logEntry := log.(*StructuredLoggerEntry)
		logEntry.logger = logEntry.logger.WithField("@timestamp", time.Now().Format(time.RFC3339Nano))
		return logEntry
	}

	if logger.logrus == nil {
		logger.newLogrus()
	}

	if isDiscard {
		logger.logrus.Out = ioutil.Discard
	}

	return &StructuredLoggerEntry{logger: logger.logrus}
}

func (logger *logger) watchLoggerChanges() {
	if !logger.cfg.LogRotation {
		logger.logrus.Info("disabled log rotation")
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Events:
				if ev.Name == logger.cfg.LogFilePath && (ev.Op.String() == "REMOVE" || ev.Op.String() == "RENAME") {
					logger.logfile.Close()
					f, err := os.OpenFile(logger.cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
					if err != nil {
						panic(err)
					}

					logger.logrus.Out = f
					logger.logfile = f

					err = watcher.Add(logger.logfile.Name())
					if err != nil {
						panic(err)
					}
				}
			case err := <-watcher.Errors:
				logger.logrus.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(logger.cfg.LogFilePath)
	if err != nil {
		panic(err)
	}
}

// Implementation of ErrorLogger interface{} with clerrors custom error
func (l *logger) WithError(err error) *logrus.Entry {
	type causer interface {
		Cause() error
	}

	var stacks interface{}
	var rootCause interface{}
	errCode := 0
	errTitle := ""
	statusCode := 0
	errMsg := err.Error()
	messages := ""
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

	return l.logrus.WithFields(fields)
}
