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

// ArticlesDelete implements [api.ServerInterface].
func (s *server) ArticlesDelete(ctx *echo.Context, slug string) error {
	panic("unimplemented")
}

// ArticlesList implements [api.ServerInterface].
func (s *server) ArticlesList(ctx *echo.Context, params api.ArticlesListParams) error {
	panic("unimplemented")
}

// ArticlesRead implements [api.ServerInterface].
func (s *server) ArticlesRead(ctx *echo.Context, slug string) error {
	panic("unimplemented")
}

// ArticlesUpsert implements [api.ServerInterface].
func (s *server) ArticlesUpsert(ctx *echo.Context, slug string) error {
	panic("unimplemented")
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
