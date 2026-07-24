package auth

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

type Authenticator interface {
	CreateAuthenticationFunc() openapi3filter.AuthenticationFunc
	GetTokenFromContext(ctx context.Context) (jwt.Token, bool)
}
