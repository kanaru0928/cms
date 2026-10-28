package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v5"
	"github.com/lestrrat-go/jwx/v4/jwt"
	echomiddleware "github.com/oapi-codegen/echo-v5-middleware"
)

type tokenContextKey struct{}

// CloudFront の OAC が Authorization ヘッダーを署名で上書きするため、独自ヘッダーでトークンを受け取る
const authorizationHeader = "KCMS-Authorization"

type jwksAuthenticator struct {
	fetcher Fetcher
}

func NewCognitoAuthenticator(ctx context.Context, jwksURL string) (*jwksAuthenticator, error) {
	cache := NewJWKSCache(jwksURL)

	return &jwksAuthenticator{
		fetcher: cache,
	}, nil
}

func (c *jwksAuthenticator) validate(ctx context.Context, token string) (jwt.Token, error) {
	jwks, err := c.fetcher.GetJWKS(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get JWKS: %w", err)
	}

	parsedToken, err := jwt.ParseString(token, jwt.WithKeySet(jwks))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return parsedToken, nil
}

func (c *jwksAuthenticator) CreateAuthenticationFunc() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		req := input.RequestValidationInput.Request
		authorization := req.Header.Get(authorizationHeader)
		if authorization == "" {
			return echo.ErrUnauthorized
		}

		if !strings.HasPrefix(authorization, "Bearer ") {
			return echo.ErrUnauthorized.Wrap(fmt.Errorf("invalid %s header format", authorizationHeader))
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		parsedToken, err := c.validate(ctx, token)
		if err != nil {
			return echo.ErrUnauthorized.Wrap(fmt.Errorf("failed to validate token: %w", err))
		}

		ctxWithToken := context.WithValue(req.Context(), tokenContextKey{}, parsedToken)
		authenticatedRequest := req.WithContext(ctxWithToken)
		input.RequestValidationInput.Request = authenticatedRequest

		if echoContext := echomiddleware.GetEchoContext(ctx); echoContext != nil {
			echoContext.SetRequest(authenticatedRequest)
		}

		return nil
	}
}

func (c *jwksAuthenticator) GetTokenFromContext(ctx context.Context) (jwt.Token, bool) {
	token, ok := ctx.Value(tokenContextKey{}).(jwt.Token)
	if !ok {
		return nil, false
	}
	return token, true
}
