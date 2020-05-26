package cllogging

import (
	"github.com/sirupsen/logrus"
)

const (
	defaultLogLevel logrus.Level = logrus.InfoLevel
)

type KV map[string]interface{}

/*
 * KVLogger: Logger interface to pass logrus fields as key and value
 *           This is useful when integrating with third parties libraries
 */
type KVLogger interface {
	Error(err error, args ...interface{})
	Errorf(err error, fmt string, args ...interface{})
	ErrorKV(err error, kv KV, args ...interface{})
	ErrorKVf(err error, kv KV, fmt string, args ...interface{})

	Warn(args ...interface{})
	Warnf(fmt string, args ...interface{})
	WarnKV(kv KV, args ...interface{})
	WarnKVf(kv KV, fmt string, args ...interface{})

	Info(args ...interface{})
	Infof(fmt string, args ...interface{})
	InfoKV(kv KV, args ...interface{})
	InfoKVf(kv KV, fmt string, args ...interface{})

	Debug(args ...interface{})
	Debugf(fmt string, args ...interface{})
	DebugKV(kv KV, args ...interface{})
	DebugKVf(kv KV, fmt string, args ...interface{})
}
