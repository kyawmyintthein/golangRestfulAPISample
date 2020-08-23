package logx

import "context"

func Error(ctx context.Context, err error, args ...interface{}) {
	getLogger().Error(ctx, err, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	getLogger().Warn(ctx, args...)
}
func Info(ctx context.Context, args ...interface{}) {
	getLogger().Info(ctx, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	getLogger().Debug(ctx, args...)
}

func Errorf(ctx context.Context, err error, msg string, args ...interface{}) {
	getLogger().Errorf(ctx, err, msg, args...)
}

func Warnf(ctx context.Context, msg string, args ...interface{}) {
	getLogger().Warnf(ctx, msg, args...)
}

func Infof(ctx context.Context, msg string, args ...interface{}) {
	getLogger().Infof(ctx, msg, args...)
}

func Debugf(ctx context.Context, msg string, args ...interface{}) {
	getLogger().Debugf(ctx, msg, args...)
}

func ErrorKV(ctx context.Context, err error, kv KV, args ...interface{}) {
	getLogger().ErrorKV(ctx, err, kv, args...)
}

func WarnKV(ctx context.Context, kv KV, args ...interface{}) {
	getLogger().WarnKV(ctx, kv, args...)
}

func InfoKV(ctx context.Context, kv KV, args ...interface{}) {
	getLogger().InfoKV(ctx, kv, args...)
}

func DebugKV(ctx context.Context, kv KV, args ...interface{}) {
	getLogger().DebugKV(ctx, kv, args...)
}

func ErrorKVf(ctx context.Context, err error, kv KV, msg string, args ...interface{}) {
	getLogger().ErrorKVf(ctx, err, kv, msg, args...)
}

func WarnKVf(ctx context.Context, kv KV, msg string, args ...interface{}) {
	getLogger().WarnKVf(ctx, kv, msg, args...)
}

func InfoKVf(ctx context.Context, kv KV, msg string, args ...interface{}) {
	getLogger().InfoKVf(ctx, kv, msg, args...)
}

func DebugKVf(ctx context.Context, kv KV, msg string, args ...interface{}) {
	getLogger().DebugKVf(ctx, kv, msg, args...)
}

func SetLogLevel(level string) error {
	return getLogger().SetLogLevel(level)
}
