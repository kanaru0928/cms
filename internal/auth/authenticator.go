package auth

import (
	"github.com/getkin/kin-openapi/openapi3filter"
)

type Authenticator interface {
	CreateAuthenticationFunc() openapi3filter.AuthenticationFunc
}
