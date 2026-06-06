package metadata

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type StatusType string

const (
	StatusPublished   StatusType = "published"
	StatusUnpublished StatusType = "unpublished"
)

type ItemType string

const (
	ItemTypeArticle ItemType = "ARTICLE"
	ItemTypeTag     ItemType = "TAG"
)

const tagAll = "#ALL"

type articleKey string

const (
	articleKeySlug     articleKey = "slug"
	articleKeyItemType articleKey = "item_type"

	articleKeyStatus     articleKey = "status"
	articleKeyFilterTag  articleKey = "filter_tag"
	articleKeySource     articleKey = "source"
	articleKeyTitle      articleKey = "title"
	articleKeyTags       articleKey = "tags"
	articleKeyContentKey articleKey = "content_key"
	articleKeyUpdatedAt  articleKey = "updated_at"
	articleKeyCreatedAt  articleKey = "created_at"

	articleKeyPV           articleKey = "pv"
	articleKeyThumbnailURL articleKey = "thumbnail_url"
)

type articlePK string

const (
	articlePKSlug     articlePK = articlePK(articleKeySlug)
	articlePKItemType articlePK = articlePK(articleKeyItemType)
)

type articlePKMap map[articlePK]types.AttributeValue

type articleItem struct {
	Slug     string   `dynamodbav:"slug"`
	ItemType ItemType `dynamodbav:"item_type"`

	// 共通のフィールド
	Status     StatusType `dynamodbav:"status"`
	FilterTag  string     `dynamodbav:"filter_tag"`
	Title      string     `dynamodbav:"title"`
	Source     string     `dynamodbav:"source"`
	ContentKey string     `dynamodbav:"content_key"`
	Tags       []string   `dynamodbav:"tags"`
	UpdatedAt  string     `dynamodbav:"updated_at"`
	CreatedAt  string     `dynamodbav:"created_at"`

	// ItemTypeArticle 固有のフィールド
	PV           int    `dynamodbav:"pv"`
	ThumbnailURL string `dynamodbav:"thumbnail_url"`
}

type tagItem struct {
	Slug     string   `dynamodbav:"slug"`
	ItemType ItemType `dynamodbav:"item_type"`

	// 共通のフィールド
	Status     StatusType `dynamodbav:"status"`
	FilterTag  string     `dynamodbav:"filter_tag"`
	Title      string     `dynamodbav:"title"`
	Source     string     `dynamodbav:"source"`
	ContentKey string     `dynamodbav:"content_key"`
	Tags       []string   `dynamodbav:"tags"`
	UpdatedAt  string     `dynamodbav:"updated_at"`
	CreatedAt  string     `dynamodbav:"created_at"`
}

type PutArticleDTO struct {
	Slug         string
	Title        string
	Source       string
	ContentKey   string
	Status       StatusType
	Tags         []string
	ThumbnailURL string
}

func NewPutArticleDTO(slug, title, contentKey, source string, status StatusType, tags []string, thumbnailURL string) (*PutArticleDTO, error) {
	slugLength := len(slug)
	if slugLength < 1 || slugLength > 100 {
		return nil, fmt.Errorf("slug must be between 1 and 100 characters")
	}

	titleLength := len(title)
	if titleLength < 1 || titleLength > 200 {
		return nil, fmt.Errorf("title must be between 1 and 200 characters")
	}

	if contentKey == "" {
		return nil, fmt.Errorf("contentKey cannot be empty")
	}

	sourceLength := len(source)
	if sourceLength < 1 || sourceLength > 100 {
		return nil, fmt.Errorf("source must be between 1 and 100 characters")
	}

	for i, tag := range tags {
		tagLength := len(tag)
		if tagLength < 1 || tagLength > 50 {
			return nil, fmt.Errorf("each tag must be between 1 and 50 characters, but tag at index %d is invalid", i)
		}
	}

	return &PutArticleDTO{
		Slug:       slug,
		Title:      title,
		ContentKey: contentKey,
		Source:     source,
		Status:     status,
		Tags:       tags,
		ThumbnailURL: thumbnailURL,
	}, nil
}

type GetArticleTagsDTO struct {
	Slug string
}

func NewGetArticleTagsDTO(slug string) (*GetArticleTagsDTO, error) {
	slugLength := len(slug)
	if slugLength < 1 || slugLength > 100 {
		return nil, fmt.Errorf("slug must be between 1 and 100 characters")
	}

	return &GetArticleTagsDTO{
		Slug: slug,
	}, nil
}

type getArticleTagsOutputItem struct {
	Tags []string `dynamodbav:"tags"`
}
