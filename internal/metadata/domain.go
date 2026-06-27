package metadata

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kanaru0928/cms/internal/myerrors"
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

type PutArticleDTOProps struct {
	Slug         string
	Title        string
	Source       string
	ContentKey   string
	Status       string
	Tags         []string
	ThumbnailURL string
}

func NewPutArticleDTO(props PutArticleDTOProps) (*PutArticleDTO, error) {
	slugLength := len(props.Slug)
	if slugLength < 1 || slugLength > 100 {
		return nil, &myerrors.ValidationError{Message: "slug must be between 1 and 100 characters"}
	}

	titleLength := len(props.Title)
	if titleLength < 1 || titleLength > 200 {
		return nil, &myerrors.ValidationError{Message: "title must be between 1 and 200 characters"}
	}

	contentKeyLength := len(props.ContentKey)
	if contentKeyLength < 1 || contentKeyLength > 200 {
		return nil, &myerrors.ValidationError{Message: "contentKey must be between 1 and 200 characters"}
	}

	if props.Status != string(StatusPublished) && props.Status != string(StatusUnpublished) {
		return nil, &myerrors.ValidationError{Message: "status must be either 'published' or 'unpublished'"}
	}

	sourceLength := len(props.Source)
	if sourceLength < 1 || sourceLength > 100 {
		return nil, &myerrors.ValidationError{Message: "source must be between 1 and 100 characters"}
	}

	tagsLength := len(props.Tags)
	if tagsLength < 1 || tagsLength > 20 {
		return nil, &myerrors.ValidationError{Message: "number of tags must be between 1 and 20"}
	}
	for i, tag := range props.Tags {
		tagLength := len(tag)
		if tagLength < 1 || tagLength > 50 {
			return nil, &myerrors.ValidationError{
				Message: fmt.Sprintf("each tag must be between 1 and 50 characters, but tag at index %d is invalid", i),
			}
		}
	}

	return &PutArticleDTO{
		Slug:         props.Slug,
		Title:        props.Title,
		ContentKey:   props.ContentKey,
		Source:       props.Source,
		Status:       StatusType(props.Status),
		Tags:         props.Tags,
		ThumbnailURL: props.ThumbnailURL,
	}, nil
}

type GetArticleTagsDTO struct {
	Slug string
}

type GetArticleTagsDTOProps struct {
	Slug string
}

func NewGetArticleTagsDTO(props GetArticleTagsDTOProps) (*GetArticleTagsDTO, error) {
	slugLength := len(props.Slug)
	if slugLength < 1 || slugLength > 100 {
		return nil, &myerrors.ValidationError{Message: "slug must be between 1 and 100 characters"}
	}

	return &GetArticleTagsDTO{
		Slug: props.Slug,
	}, nil
}

type getArticleTagsOutputItem struct {
	Tags []string `dynamodbav:"tags"`
}

type ListArticlesDTO struct {
	tag    string
	status StatusType
	limit  int32
}

type ListArticlesDTOProps struct {
	Tag    string
	Status string
	Limit  int32
}

func NewListArticlesDTO(props ListArticlesDTOProps) (*ListArticlesDTO, error) {
	if props.Limit < 1 || props.Limit > 100 {
		return nil, &myerrors.ValidationError{Message: "limit must be between 1 and 100"}
	}

	if props.Status != string(StatusPublished) && props.Status != string(StatusUnpublished) {
		return nil, &myerrors.ValidationError{Message: "status must be either 'published' or 'unpublished'"}
	}

	return &ListArticlesDTO{
		tag:    props.Tag,
		status: StatusType(props.Status),
		limit:  props.Limit,
	}, nil
}

func (dto *ListArticlesDTO) GetTag() string {
	return dto.tag
}

func (dto *ListArticlesDTO) GetStatus() StatusType {
	return dto.status
}

func (dto *ListArticlesDTO) GetLimit() int32 {
	return dto.limit
}
