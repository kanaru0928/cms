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
	_, found := s.authenticator.GetTokenFromContext(ctx.Request().Context())

	tag := ""
	if params.Tag != nil {
		tag = *params.Tag
	}

	limit := 10
	if params.Limit != nil {
		limit = int(*params.Limit)
	}

	status := string(metadata.StatusPublished)
	if params.Status != nil {
		if !found && string(*params.Status) == string(metadata.StatusUnpublished) {
			return ctx.JSON(403, api.Error{Error: "unauthorized to access unpublished articles"})
		}
		status = string(*params.Status)
	}

	order := string(metadata.SortOrderDesc)
	if params.Order != nil {
		order = string(*params.Order)
	}

	listArticlesDTO, err := metadata.NewListArticlesDTO(metadata.ListArticlesDTOProps{
		Tag: tag,
		Limit: limit,
		Status: status,
		Order: order,
	})
	if err != nil {
		return fmt.Errorf("failed to create ListArticlesDTO: %w", err)
	}

	articles, err := s.metadataRepository.ListArticles(ctx.Request().Context(), listArticlesDTO)
	if err != nil {
		return fmt.Errorf("failed to list articles: %w", err)
	}

	response := make([]api.ArticleMetadata, len(articles.Items))
	for i, article := range articles.Items {
		response[i] = api.ArticleMetadata{
			Slug: article.Slug,
			Title: article.Title,
			Tags: article.Tags,
			Status: api.ArticleStatus(article.Status),
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
			Pv: float32(article.PV),
			Source: article.Source,
			ThumbnailUrl: article.ThumbnailURL,
		}
	}

	return ctx.JSON(200, response)
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
