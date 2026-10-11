//go:generate go tool mockgen -source=repository.go -destination=mock_repository.go -package=metadata

package metadata

import "context"

type Repository interface {
	GetArticleTags(ctx context.Context, getArticleTagsDTO *getArticleTagsDTO) ([]string, error)
	PutArticle(ctx context.Context, props *putArticleDTO) (*PutArticleOutput, error)
	ListArticles(ctx context.Context, listArticlesDTO *listArticlesDTO) (*ListArticlesOutput, error)
	GetArticle(ctx context.Context, getArticleDTO *getArticleDTO) (*GetArticleOutput, error)
	DeleteArticle(ctx context.Context, deleteArticleDTO *deleteArticleDTO) error
}
