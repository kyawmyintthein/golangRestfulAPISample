package logging

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

/*
 * default logger implementation of KVLogger interface
 */
func NewKVLogger(cfg *LogCfg) KVLogger {
	log := logrus.New()
	logLevel, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.Level = logLevel
	log.SetReportCaller(cfg.CallerInfo)

	if cfg.JsonLogFormat {
		log.Formatter = new(JSONFormatter)
	} else {
		log.Formatter = new(logrus.TextFormatter)
	}

	logger := &logger{
		logrus: log,
	}

	if cfg.LogFilePath == "" {
		log.Out = ioutil.Discard
		log.Errorln("Empty log file. Set 'Discard' as default")
	}else{
		logfile, err := os.OpenFile(logger.cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			log.Out = ioutil.Discard
			logger.logrus.Errorln("Failed to set log file : '%v'. Set 'Stdout' as default", err)
		}else{
			logger.logfile = logfile
			logger.logrus.Out = logger.logfile
		}
	}
	log.Info("Initialized Logger successfully")
	return logger
}

func DefaultKVLogger() KVLogger {
	logrus := logrus.New()
	logrus.SetLevel(defaultLogLevel)
	logrus.SetFormatter(&JSONFormatter{})
	logrus.SetOutput(ioutil.Discard)

	return &logger{
		logrus: logrus,
	}
}

func (logger *logger) Error(err error, args ...interface{}) {
	logger.logrus.WithError(err).Errorln(args...)
}

func (logger *logger) Errorf(err error, fmt string, args ...interface{}) {
	logger.logrus.WithError(err).Errorf(fmt, args...)
}

func (logger *logger) ErrorKV(err error, kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).WithError(err).Errorln(args...)
}

func (logger *logger) ErrorKVf(err error, kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).WithError(err).Errorf(fmt, args...)
}

func (logger *logger) Warn(args ...interface{}) {
	logger.logrus.Warnln(args...)
}

func (logger *logger) Warnf(fmt string, args ...interface{}) {
	logger.logrus.Warnf(fmt, args)
}

func (logger *logger) WarnKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Warnln(args...)
}

func (logger *logger) WarnKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Warnf(fmt, args...)
}

func (logger *logger) Info(args ...interface{}) {
	logger.logrus.Infoln(args...)
}

func (logger *logger) Infof(fmt string, args ...interface{}) {
	logger.logrus.Infof(fmt, args...)
}

func (logger *logger) InfoKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Infoln(args...)
}

func (logger *logger) InfoKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Infof(fmt, args...)
}

func (logger *logger) Debug(args ...interface{}) {
	logger.logrus.Debugln(args...)
}

func (logger *logger) Debugf(fmt string, args ...interface{}) {
	logger.logrus.Debugf(fmt, args...)
}

func (logger *logger) DebugKV(kv KV, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Debugln(args...)
}

func (logger *logger) DebugKVf(kv KV, fmt string, args ...interface{}) {
	fields := convertKvsToLogrusFields(kv)
	logger.logrus.WithFields(fields).Debugf(fmt, args...)
}