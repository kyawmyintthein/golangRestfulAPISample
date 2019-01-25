package logging

import (
	"context"
	"github.com/kyawmyintthein/golangRestfulAPISample/internal/logjsonformatter"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

type DiscardLoggingKey struct{}

type Logger interface{
	NewStructuredLogger() func(next http.Handler) http.Handler
	GetLogger(ctx context.Context) *StructuredLoggerEntry
	Logrus() *logrus.Logger
	SetLogLevel(level string) error
}

type logger struct{
	logRotation     bool
	isJsonLogFormat bool
	logLevel    string
	logfilepath string
	logrusLogger *logrus.Logger
	logfile *os.File
}

func InitializeLogger(logLevel string, logFilepath string, isJsonLogFormat bool, logRotation bool) Logger {

	logger := &logger{
		logLevel: logLevel,
		logfilepath: logFilepath,
		logRotation: logRotation,
		isJsonLogFormat: isJsonLogFormat,
	}

	logger.initLogrus()

	go logger.watchLoggerChanges()
	return logger
}

func (logger *logger) initLogrus(){
	logger.logrusLogger = &logrus.Logger{
		Hooks:     make(logrus.LevelHooks),
	}

	logLevel, err := logrus.ParseLevel(logger.logLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logger.logrusLogger.Level = logLevel

	if logger.isJsonLogFormat{
		logger.logrusLogger.Formatter = new(logjsonformatter.JSONFormatter)
	} else {
		logger.logrusLogger.Formatter = new(logrus.TextFormatter)
	}

	if logger.logfilepath == "" {
		logger.logrusLogger.Out = os.Stdout
		logger.logrusLogger.Errorln("Empty log file. Set 'Stdout' as default")
		logger.logrusLogger.Info("Initialized Logger successfully")
		return
	}

	logfile, err := os.OpenFile(logger.logfilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.logrusLogger.Errorln("Failed to set log file : '%v'. Set 'Stdout' as default", err)
		return
	}

	logger.logfile = logfile
	logger.logrusLogger.Out = logger.logfile

	logger.logrusLogger.Info("Initialized Logger successfully")
}

func (logger *logger) Logrus() *logrus.Logger{
	return logger.logrusLogger
}

func (logger *logger) SetLogLevel(level string) error{
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
		return err
	}
	logger.logrusLogger.Level = logLevel
	return nil
}

func (logger *logger) GetLogger(ctx context.Context) *StructuredLoggerEntry {
	log := ctx.Value(middleware.LogEntryCtxKey)

	var isDiscard bool
	i := ctx.Value(DiscardLoggingKey{})
	if i != nil{
		isDiscard = i.(bool)
	}

	if log != nil {
		logEntry := log.(*StructuredLoggerEntry)
		logEntry.logger = logEntry.logger.WithField("@timestamp", time.Now().Format(time.RFC3339Nano))
		return logEntry
	}

	if logger.logrusLogger == nil{
		logger.initLogrus()
	}

	if isDiscard{
		logger.logrusLogger.Out = ioutil.Discard
	}

	return &StructuredLoggerEntry{logger: logger.logrusLogger}
}

func (logger *logger) watchLoggerChanges() {
	if !logger.logRotation {
		logger.logrusLogger.Info("Disable log rotation")
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
				if ev.Name == logger.logLevel && (ev.Op.String() == "REMOVE" || ev.Op.String() == "RENAME") {
					logger.logfile.Close()
					f, err := os.OpenFile(logger.logLevel, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
					if err != nil{
						panic(err)
					}

					logger.logrusLogger.Out = f
					logger.logfile = f

					err = watcher.Add(logger.logfile.Name())
					if err != nil {
						panic(err)
					}
				}
			case err := <-watcher.Errors:
				logger.logrusLogger.Error("error:", err)
			}
		}
	}()

	err = watcher.Add(logger.logfilepath)
	if err != nil {
		panic(err)
	}
}