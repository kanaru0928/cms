package handler

import (
	"context"
	"fmt"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/metadata"
)

// ArticlesDelete implements [api.StrictServerInterface].
func (s *server) ArticlesDelete(ctx context.Context, request api.ArticlesDeleteRequestObject) (api.ArticlesDeleteResponseObject, error) {
	panic("unimplemented")
}

// ArticlesList implements [api.StrictServerInterface].
func (s *server) ArticlesList(ctx context.Context, request api.ArticlesListRequestObject) (api.ArticlesListResponseObject, error) {
	_, found := s.authenticator.GetTokenFromContext(ctx)

	tag := ""
	if request.Params.Tag != nil {
		tag = *request.Params.Tag
	}

	limit := 10
	if request.Params.Limit != nil {
		limit = int(*request.Params.Limit)
	}

	status := string(metadata.StatusPublished)
	if request.Params.Status != nil {
		if !found && string(*request.Params.Status) == string(metadata.StatusUnpublished) {
			return api.ArticlesList401JSONResponse{Error: "unauthorized to access unpublished articles"}, nil
		}
		status = string(*request.Params.Status)
	}

	order := string(metadata.SortOrderDesc)
	if request.Params.Order != nil {
		order = string(*request.Params.Order)
	}

	listArticlesDTO, err := metadata.NewListArticlesDTO(metadata.ListArticlesDTOProps{
		Tag:    tag,
		Limit:  int32(limit),
		Status: status,
		Order:  order,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ListArticlesDTO: %w", err)
	}

	articles, err := s.metadataRepository.ListArticles(ctx, listArticlesDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to list articles: %w", err)
	}

	response := make([]api.ArticleMetadata, len(articles.Items))
	for i, article := range articles.Items {
		response[i] = api.ArticleMetadata{
			Slug:         article.Slug,
			Title:        article.Title,
			Tags:         article.Tags,
			Status:       api.ArticleStatus(article.Status),
			CreatedAt:    article.CreatedAt,
			UpdatedAt:    article.UpdatedAt,
			Pv:           float32(article.PV),
			Source:       article.Source,
			ThumbnailUrl: article.ThumbnailURL,
		}
	}

	return api.ArticlesList200JSONResponse(response), nil
}

// ArticlesRead implements [api.StrictServerInterface].
func (s *server) ArticlesRead(ctx context.Context, request api.ArticlesReadRequestObject) (api.ArticlesReadResponseObject, error) {
	getArticleDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps{
		Slug: request.Slug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GetArticleDTO: %w", err)
	}

	article, err := s.metadataRepository.GetArticle(ctx, getArticleDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	content, err := s.contentRepository.GetTextContent(ctx, article.ContentKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	response := api.ArticleDetail{
		Content:      string(content),
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
		Slug:         article.Slug,
		Tags:         article.Tags,
		Pv:           float32(article.PV),
		Source:       article.Source,
		Title:        article.Title,
		Status:       api.ArticleStatus(article.Status),
		ThumbnailUrl: article.ThumbnailURL,
	}

	return api.ArticlesRead200JSONResponse(response), nil
}

// ArticlesUpsert implements [api.StrictServerInterface].
func (s *server) ArticlesUpsert(ctx context.Context, request api.ArticlesUpsertRequestObject) (api.ArticlesUpsertResponseObject, error) {
	panic("unimplemented")
}
