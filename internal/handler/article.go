package handler

import (
	"github.com/kanaru0928/cms/internal/api"
	"github.com/labstack/echo/v5"
)

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
