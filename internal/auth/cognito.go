package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

type cognitoAuthenticator struct {
	fetcher JWKSFetcher
}

func NewCognitoAuthenticator(ctx context.Context, jwksURL string) (*cognitoAuthenticator, error) {
	cache := NewJWKSCache(jwksURL)

	return &cognitoAuthenticator{
		fetcher: cache,
	}, nil
}

func (c *cognitoAuthenticator) validate(ctx context.Context, token string) (jwt.Token, error) {
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

func (c *cognitoAuthenticator) CreateAuthenticationFunc() openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		req := input.RequestValidationInput.Request
		authorization := req.Header.Get("Authorization")
		if authorization == "" {
			return fmt.Errorf("missing Authorization header")
		}

		if !strings.HasPrefix(authorization, "Bearer ") {
			return fmt.Errorf("invalid Authorization header format")
		}

		token := strings.TrimPrefix(authorization, "Bearer ")

		parsedToken, err := c.validate(ctx, token)
		if err != nil {
			return fmt.Errorf("failed to validate token: %w", err)
		}

		ctxWithToken := context.WithValue(ctx, "token", parsedToken)
		input.RequestValidationInput.Request = req.WithContext(ctxWithToken)

		return nil
	}
}

func (c *cognitoAuthenticator) GetTokenFromContext(ctx context.Context) (jwt.Token, bool) {
	token, ok := ctx.Value("token").(jwt.Token)
	if !ok {
		return nil, false
	}
	return token, true
}
