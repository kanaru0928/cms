package handler

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/auth"
	"github.com/labstack/echo/v5"
	middleware "github.com/oapi-codegen/echo-v5-middleware"
)

type server struct {
	authenticator auth.Authenticator
}

func NewServer(authenticator auth.Authenticator) *server {
	return &server{
		authenticator: authenticator,
	}
}

func (s *server) CreateMiddleware() ([]echo.MiddlewareFunc, error) {
	spec, err := api.GetSpec()
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAPI spec: %w", err)
	}

	validator := middleware.OapiRequestValidatorWithOptions(spec, &middleware.Options{
		Options: openapi3filter.Options{
			AuthenticationFunc: s.authenticator.CreateAuthenticationFunc(),
		},
	})

	return []echo.MiddlewareFunc{validator}, nil
}
