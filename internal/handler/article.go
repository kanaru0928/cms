package handler

import (
	"fmt"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/metadata"
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
	getArticleDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps{
		Slug: slug,
	})
	if err != nil {
		return fmt.Errorf("failed to create GetArticleDTO: %w", err)
	}

	article, err := s.metadataRepository.GetArticle(ctx.Request().Context(), getArticleDTO)
	if err != nil {
		return fmt.Errorf("failed to get article: %w", err)
	}

	content, err := s.contentRepository.GetTextContent(ctx.Request().Context(), article.ContentKey)
	if err != nil {
		return fmt.Errorf("failed to get content: %w", err)
	}

	response := api.ArticleDetail{
		Content: string(content),
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		Slug: article.Slug,
		Tags: article.Tags,
		Pv: float32(article.PV),
		Source: article.Source,
		Title: article.Title,
		Status: api.ArticleStatus(article.Status),
		ThumbnailUrl: article.ThumbnailURL,
	}

	return ctx.JSON(200, response)
}

// ArticlesUpsert implements [api.ServerInterface].
func (s *server) ArticlesUpsert(ctx *echo.Context, slug string) error {
	panic("unimplemented")
}
