package helpers

import (
	"github.com/pkg/errors"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type Error struct {
	err     error
	source  logger.Sourceable
	message string
}

func (e *Error) Error() string {
	if e.message == "" {
		return e.err.Error()
	}
	return e.message + ": " + e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) LogTrace() string {
	c, ok := e.err.(logger.Traceable)
	if !ok {
		return e.source.LogSource()
	}
	return e.source.LogSource() + "->" + c.LogTrace()
}

func NewErr(source logger.Sourceable, message string) error {
	return &Error{
		err:    errors.New(message),
		source: source,
	}
}

func WrapErr(err error, source logger.Sourceable, message string) error {
	return &Error{
		err:     err,
		source:  source,
		message: message,
	}
}
