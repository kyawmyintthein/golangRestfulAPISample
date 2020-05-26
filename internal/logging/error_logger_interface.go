package logging

import "github.com/sirupsen/logrus"

/*
 * ErrorLogger: Logger interface to support custom error type
 */
type ErrorLogger interface {
	WithError(error) *logrus.Entry
}
