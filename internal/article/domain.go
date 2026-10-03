package article

import "github.com/kanaru0928/cms/internal/metadata"

type ListArticlesDTOProps metadata.ListArticlesDTOProps

type ListArticlesOutput metadata.ListArticlesOutput

type GetArticleDTOProps metadata.GetArticleDTOProps

type GetArticleOutput metadata.GetArticleOutput

type UpsertArticleProps struct {
	Slug         string
	Title        string
	ContentBody  string
	Source       string
	Status       string
	Tags         []string
	ThumbnailURL string
}

type UpsertArticleOutput struct {
	// 記事が新規作成された場合は true、既存の記事が更新された場合は false
	Created bool
	Article *GetArticleOutput
}
