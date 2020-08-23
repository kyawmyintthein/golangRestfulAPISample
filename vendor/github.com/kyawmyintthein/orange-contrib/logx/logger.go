package logx

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/go-chi/chi/middleware"
	"github.com/kyawmyintthein/orange-contrib/errorx"
	"github.com/sirupsen/logrus"
)

type KV map[string]interface{}

type Logger interface {
	Error(context.Context, error, ...interface{})
	Warn(context.Context, ...interface{})
	Info(context.Context, ...interface{})
	Debug(context.Context, ...interface{})

	Errorf(context.Context, error, string, ...interface{})
	Warnf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
	Debugf(context.Context, string, ...interface{})

	ErrorKV(context.Context, error, KV, ...interface{})
	WarnKV(context.Context, KV, ...interface{})
	InfoKV(context.Context, KV, ...interface{})
	DebugKV(context.Context, KV, ...interface{})

	ErrorKVf(context.Context, error, KV, string, ...interface{})
	WarnKVf(context.Context, KV, string, ...interface{})
	InfoKVf(context.Context, KV, string, ...interface{})
	DebugKVf(context.Context, KV, string, ...interface{})

	SetLogLevel(string) error
	NewRequestLogger() func(next http.Handler) http.Handler
}

type logger struct {
	cfg     *LogCfg
	logrus  *logrus.Logger
	logfile *os.File
	absPath string
}

const (
	jsonLogFormat   = "json"
	textLogFormat   = "text"
	defaultLogLevel = logrus.InfoLevel
)

var (
	_defaultFileFlag = os.O_APPEND | os.O_CREATE | os.O_WRONLY
	_defaultFileMode = os.FileMode(0755)
)

var (
	_stdLogger Logger
)

func Init(cfg *LogCfg) Logger {
	if _stdLogger == nil {
		_stdLogger = new(cfg)
	}
	return _stdLogger
}

func New(cfg *LogCfg) Logger {
	return new(cfg)
}

func new(cfg *LogCfg) Logger {
	logger := &logger{
		cfg: cfg,
	}
	logger.newLogrus()

	go logger.watchLogRotation()

	return logger
}

// create new logrus object
func (logger *logger) newLogrus() {
	logger.logrus = &logrus.Logger{
		Hooks: make(logrus.LevelHooks),
	}

	logLevel, err := logrus.ParseLevel(logger.cfg.LogLevel)
	if err != nil {
		logLevel = defaultLogLevel
	}
	logger.logrus.Level = logLevel

	switch logger.cfg.LogFormat {
	case jsonLogFormat:
		logger.logrus.SetFormatter(&logrus.JSONFormatter{})
	default:
		logger.logrus.SetFormatter(&logrus.TextFormatter{})
	}

	if logger.cfg.LogFilePath == "" {
		logger.logrus.Out = os.Stdout
		logger.logrus.Errorf("[%s]:: empty log file. Set 'Stdout' as default \n", PackageName)
		logger.logrus.Infof("[%s]:: initialized logx successfully \n", PackageName)
		return
	}

	logfile, err := os.OpenFile(logger.cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		logger.logrus.Errorln("[%s]:: failed to set log file. Error : '%v'. Set 'Stdout' as default", PackageName, err)
		return
	}

	logger.logfile = logfile
	logger.logrus.Out = logger.logfile

	logger.logrus.Infof("[%s]:: initialized logx successfully", PackageName)
}

func (logger *logger) watchLogRotation() {
	if !logger.cfg.LogRotation {
		logger.Infof(context.Background(), "[%s]:: disabled log rotation", PackageName)
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
				if ev.Name == logger.absPath && (ev.Op.String() == "REMOVE" || ev.Op.String() == "RENAME") {
					logger.logfile.Close()
					f, err := os.OpenFile(logger.cfg.LogFilePath, _defaultFileFlag, _defaultFileMode)
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
				logger.Errorf(context.Background(), err, "[%s]:: failed to watch log file", PackageName)
			}
		}
	}()

	err = watcher.Add(logger.cfg.LogFilePath)
	if err != nil {
		panic(err)
	}
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

func getLogger() Logger {
	if _stdLogger == nil {
		return new(&LogCfg{})
	}
	return _stdLogger
}

func (l *logger) Error(ctx context.Context, err error, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(getErrorFields(err)).Errorln(args...)
		return
	}
	l.logrus.WithFields(getErrorFields(err)).Errorln(args...)
}

func (l *logger) Warn(ctx context.Context, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Warnln(args...)
		return
	}
	l.logrus.Warnln(args...)
}

func (l *logger) Info(ctx context.Context, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Infoln(args...)
		return
	}
	l.logrus.Infoln(args...)
}

func (l *logger) Debug(ctx context.Context, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Debugln(args...)
		return
	}
	l.logrus.Debugln(args...)
}

func (l *logger) Errorf(ctx context.Context, err error, message string, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(getErrorFields(err)).Errorf(message, args...)
		return
	}
	l.logrus.WithFields(getErrorFields(err)).Errorf(message, args...)
}

func (l *logger) Warnf(ctx context.Context, message string, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Warnf(message, args...)
		return
	}
	l.logrus.Warnf(message, args...)
}

func (l *logger) Infof(ctx context.Context, message string, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Infof(message, args...)
		return
	}
	l.logrus.Infof(message, args...)
}

func (l *logger) Debugf(ctx context.Context, message string, args ...interface{}) {
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.Debugf(message, args...)
		return
	}
	l.logrus.Debugf(message, args...)
}

func (l *logger) ErrorKV(ctx context.Context, err error, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	errorFields := getErrorFields(err)
	for k, v := range errorFields {
		fields[k] = v
	}
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Errorln(args...)
		return
	}
	l.logrus.WithFields(fields).Errorln(args...)
}

func (l *logger) WarnKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Warnln(args...)
	}
	l.logrus.WithFields(fields).Warnln(args...)
}

func (l *logger) InfoKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Info(args...)
	}
	l.logrus.WithFields(fields).Info(args...)
}

func (l *logger) DebugKV(ctx context.Context, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Debug(args...)
		return
	}
	l.logrus.WithFields(fields).Debug(args...)
}

func (l *logger) ErrorKVf(ctx context.Context, err error, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	errorFields := getErrorFields(err)
	for k, v := range errorFields {
		fields[k] = v
	}
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Errorf(message, args...)
		return
	}
	l.logrus.WithFields(fields).Errorf(message, args...)
}

func (l *logger) WarnKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Warnf(message, args...)
		return
	}
	l.logrus.WithFields(fields).Warnf(message, args...)
}

func (l *logger) InfoKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Infof(message, args...)
		return
	}
	l.logrus.WithFields(fields).Infof(message, args...)
}

func (l *logger) DebugKVf(ctx context.Context, kv KV, message string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	log := l.getStructuredLogEntry(ctx)
	if log != nil {
		log.logger.WithFields(fields).Debugf(message, args...)
		return
	}
	l.logrus.WithFields(fields).Debugf(message, args...)
}

// Implementation of ErrorLogger interface{} with clerrors custom error
func getErrorFields(err error) logrus.Fields {
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
	errStacktrace, ok := err.(errorx.StackTracer)
	if ok {
		stacks = errStacktrace.GetStackAsJSON()
	}

	errWithCode, ok := err.(errorx.ErrorCode)
	if ok {
		errCode = errWithCode.Code()
	}

	errWithName, ok := err.(errorx.ErrorID)
	if ok {
		errTitle = errWithName.ID()
	}

	httpError, ok := err.(errorx.HttpError)
	if ok {
		statusCode = httpError.StatusCode()
	}

	errorFormatter, ok := err.(errorx.ErrorFormatter)
	if ok {
		errMsg = errorFormatter.FormattedMessage()
	}

	errorCauser, ok := err.(causer)
	if ok {
		rootCause = errorCauser.Cause()
	}

	messages = errorx.GetErrorMessages(err)
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

	return fields
}

func convertKvsToLogrusFields(kv KV) logrus.Fields {
	fields := make(logrus.Fields)
	for k, v := range kv {
		fields[k] = v
	}
	return fields
}

func (logger *logger) getStructuredLogEntry(ctx context.Context) *StructuredLoggerEntry {
	log := ctx.Value(middleware.LogEntryCtxKey)
	if log == nil {
		log = ctx.Value(middleware.LogEntryCtxKey.String())
	}
	if log != nil {
		logEntry := log.(*StructuredLoggerEntry)
		logEntry.logger = logEntry.logger.WithField("@timestamp", time.Now().Format(time.RFC3339Nano))
		return logEntry
	}

	return nil
}

func (logger *logger) NewRequestLogger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&RequestStructureLogger{logger.logrus})
}
