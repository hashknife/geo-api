package middleware

import (
	"context"
	"net/http"
	"strings"

	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/hashknife/geo-api/endpoints"
)

const (
	// AuthHeader holds the header field name
	AuthHeader = "X-Raven-Token"
	// Path
	Path = "Request-Path"
)

// HashknifeRequestAuthenticator
type HashknifeRequestAuthenticator struct {
	authToken    *string
	skippedPaths []string
}

// NewHashknifeRequestAuthenticator
func NewHashknifeRequestAuthenticator(authToken *string) *HashknifeRequestAuthenticator {
	return &HashknifeRequestAuthenticator{
		authToken: authToken,
		skippedPaths: []string{
			"/geo-api/healthcheck",
		},
	}
}

// skipAuth
func (a *HashknifeRequestAuthenticator) skipAuth(ctx context.Context) bool {
	if a.authToken == nil {
		return true
	}

	path, ok := ctx.Value(Path).(string)
	if !ok {
		return false
	}

	for _, s := range a.skippedPaths {
		if strings.HasPrefix(path, s) {
			return true
		}
	}

	return false
}

// verifyTokenHeader
func (a *HashknifeRequestAuthenticator) verifyTokenHeader(ctx context.Context) error {
	if a.skipAuth(ctx) {
		return nil
	}
	authToken, ok := ctx.Value(AuthHeader).(string)
	if !ok {
		return endpoints.NewForbiddenError()
	} else if authToken != *a.authToken {
		return endpoints.NewForbiddenError()
	}
	return nil
}

// EndpointAuthenticate provides an endpoint middleware
func (a *HashknifeRequestAuthenticator) EndpointAuthenticate() kitendpoint.Middleware {
	return func(next kitendpoint.Endpoint) kitendpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (i interface{}, err error) {
			authErr := a.verifyTokenHeader(ctx)
			if authErr != nil {
				return nil, authErr
			}
			return next(ctx, request)
		}
	}
}

// KitServerBefore
func KitServerBefore(ctx context.Context, r *http.Request) context.Context {
	ctx = context.WithValue(ctx, AuthHeader, r.Header.Get(AuthHeader))
	ctx = context.WithValue(ctx, Path, r.URL.Path)
	return ctx
}
