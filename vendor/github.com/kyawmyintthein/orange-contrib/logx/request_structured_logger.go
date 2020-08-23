package logx

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

/*
 * RequestStructureLogger: implementation of Chi LogFormatter interface
 */
type RequestStructureLogger struct {
	logger *logrus.Logger
}

func (l *RequestStructureLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{logger: logrus.NewEntry(l.logger)}
	logFields := logrus.Fields{}

	logFields["@timestamp"] = time.Now().Format(time.RFC3339Nano)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.logger = entry.logger.WithFields(logFields)

	entry.logger.Infoln("request started")

	return entry
}

func (l *RequestStructureLogger) NewStructuredEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{logger: logrus.NewEntry(l.logger)}

	logFields := logrus.Fields{}

	logFields["@timestamp"] = time.Now().Format(time.RFC3339Nano)

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	entry.logger = entry.logger.WithFields(logFields)

	entry.logger.Infoln("request started")
	return entry
}
