package logger

import (
	"context"
	"sync"
)

const (
	LevelCritical = "crit"
	LevelError    = "error"
	LevelInfo     = "info"
	LevelDebug    = "debug"
)

const (
	ModeDevelopment = "dev"
	ModeProduction  = "prod"
)

const (
	FieldComponent = "component"
)

type Sourceable interface {
	LogSource() string
}

type Traceable interface {
	LogTrace() string
}

var ctxKey struct{}

type message struct {
	ctx    context.Context
	fields map[string]interface{}
}

func Enrich(ctx context.Context, key string, value interface{}) context.Context {
	contextData := ctx.Value(ctxKey)

	if contextData == nil {
		values := &sync.Map{}
		values.Store(key, value)
		return context.WithValue(ctx, ctxKey, values)

	} else if values, ok := contextData.(*sync.Map); ok {
		values.Store(key, value)
		return context.WithValue(ctx, ctxKey, values)

	} else {
		return context.WithValue(ctx, ctxKey, contextData)
	}
}

func New(ctx context.Context) *message {
	fields := make(map[string]interface{})
	contextData := ctx.Value(ctxKey)
	if contextData != nil {
		if values, ok := contextData.(*sync.Map); ok {
			values.Range(func(key, value any) bool {
				if k, ok := key.(string); ok {
					fields[k] = value
				}
				return true
			})
		}
	}
	return &message{ctx: ctx, fields: fields}
}

func (m *message) Field(key string, value interface{}) *message {
	m.fields[key] = value
	return m
}

func (m *message) Src(source Sourceable) *message {
	m.fields["src"] = source.LogSource()
	return m
}

func (m *message) Err(err error) *message {
	if traceable, ok := err.(Traceable); ok {
		m.fields["trace"] = traceable.LogTrace()
	}
	m.fields["err"] = err.Error()
	return m
}
