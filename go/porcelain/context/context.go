package context

import (
	"context"

	"github.com/go-openapi/runtime"
	"github.com/sirupsen/logrus"
)

type Context interface {
	context.Context
}

type Fields map[interface{}]interface{}

func WithAuthInfo(ctx Context, authInfo runtime.ClientAuthInfoWriter) Context {
	return context.WithValue(ctx, "netlify.auth_info", authInfo)
}

func GetAuthInfo(ctx Context) runtime.ClientAuthInfoWriter {
	return ctx.Value("netlify.auth_info").(runtime.ClientAuthInfoWriter)
}

func WithLogger(ctx Context, entry *logrus.Entry) Context {
	return context.WithValue(ctx, "netlify.logger", entry)
}

func GetLogger(ctx Context) *logrus.Entry {
	logger := ctx.Value("netlify.logger")
	if logger == nil {
		return logrus.NewEntry(logrus.StandardLogger())
	}
	return logger.(*logrus.Entry)
}
