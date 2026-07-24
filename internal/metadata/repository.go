package metadata

import "context"

type Repository interface {
	GetArticleTags(ctx context.Context, getArticleTagsDTO *getArticleTagsDTO) ([]string, error)
	PutArticle(ctx context.Context, props *putArticleDTO) error
	ListArticles(ctx context.Context, listArticlesDTO *listArticlesDTO) (*listArticlesOutput, error)
	GetArticle(ctx context.Context, getArticleDTO *getArticleDTO) (*getArticleOutput, error)
}
