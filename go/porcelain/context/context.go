package context

import (
	"github.com/docker/distribution/context"
	"github.com/go-openapi/runtime"
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
