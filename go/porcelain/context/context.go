package context

import (
	"github.com/docker/distribution/context"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

const defaultUserAgent = "netlify-go-client"

type Context interface {
	context.Context
}

type Fields map[interface{}]interface{}

func WithAuthToken(ctx Context, authToken string) Context {
	res := context.WithValue(ctx, "netlify.auth_token", authToken)
	return res
}

func GetAuthWriter(ctx Context) runtime.ClientAuthInfoWriter {
	writer := runtime.ClientAuthInfoWriterFunc(func(r runtime.ClientRequest, _ strfmt.Registry) error {
		r.SetHeaderParam("Authorization", "Bearer "+GetAuthToken(ctx))
		r.SetHeaderParam("User-Agent", GetUserAgent(ctx))
		return nil
	})

	return writer
}

func GetAuthToken(ctx Context) string {
	return ctx.Value("netlify.auth_token").(string)
}

func WithUserAgent(ctx Context, userAgent string) Context {
	return context.WithValue(ctx, "netlify.user_agent", userAgent)
}

func GetUserAgent(ctx Context) string {
	userAgent, found := ctx.Value("netlify.user_agent").(string)
	if !found || userAgent == "" {
		userAgent = defaultUserAgent
	}
	return userAgent
}
