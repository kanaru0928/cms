//go:generate go tool mockgen -source=authenticator.go -destination=mock_authenticator.go -package=auth

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
