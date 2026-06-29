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

const (
	slugMinLength         = 1
	slugMaxLength         = 100
	titleMinLength        = 1
	titleMaxLength        = 200
	contentKeyMinLength   = 1
	contentKeyMaxLength   = 200
	sourceMinLength       = 1
	sourceMaxLength       = 100
	tagsMinCount          = 1
	tagsMaxCount          = 20
	tagMinLength          = 1
	tagMaxLength          = 50
	thumbnailURLMinLength = 1
	thumbnailURLMaxLength = 2000
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

type putArticleDTO struct {
	slug         string
	title        string
	source       string
	contentKey   string
	status       StatusType
	tags         []string
	thumbnailURL string
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

func NewPutArticleDTO(props PutArticleDTOProps) (*putArticleDTO, error) {
	slugLength := len(props.Slug)
	if slugLength < slugMinLength || slugLength > slugMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("slug must be between %d and %d characters", slugMinLength, slugMaxLength)}
	}

	titleLength := len(props.Title)
	if titleLength < titleMinLength || titleLength > titleMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("title must be between %d and %d characters", titleMinLength, titleMaxLength)}
	}

	contentKeyLength := len(props.ContentKey)
	if contentKeyLength < contentKeyMinLength || contentKeyLength > contentKeyMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("contentKey must be between %d and %d characters", contentKeyMinLength, contentKeyMaxLength)}
	}

	if props.Status != string(StatusPublished) && props.Status != string(StatusUnpublished) {
		return nil, &myerrors.ValidationError{Message: "status must be either 'published' or 'unpublished'"}
	}

	sourceLength := len(props.Source)
	if sourceLength < sourceMinLength || sourceLength > sourceMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("source must be between %d and %d characters", sourceMinLength, sourceMaxLength)}
	}

	tagsLength := len(props.Tags)
	if tagsLength < tagsMinCount || tagsLength > tagsMaxCount {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("number of tags must be between %d and %d", tagsMinCount, tagsMaxCount)}
	}
	for i, tag := range props.Tags {
		tagLength := len(tag)
		if tagLength < tagMinLength || tagLength > tagMaxLength {
			return nil, &myerrors.ValidationError{
				Message: fmt.Sprintf("each tag must be between %d and %d characters, but tag at index %d is invalid", tagMinLength, tagMaxLength, i),
			}
		}
	}

	return &putArticleDTO{
		slug:         props.Slug,
		title:        props.Title,
		contentKey:   props.ContentKey,
		source:       props.Source,
		status:       StatusType(props.Status),
		tags:         props.Tags,
		thumbnailURL: props.ThumbnailURL,
	}, nil
}

type getArticleTagsDTO struct {
	slug string
}

type GetArticleTagsDTOProps struct {
	Slug string
}

func NewGetArticleTagsDTO(props GetArticleTagsDTOProps) (*getArticleTagsDTO, error) {
	slugLength := len(props.Slug)
	if slugLength < slugMinLength || slugLength > slugMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("slug must be between %d and %d characters", slugMinLength, slugMaxLength)}
	}

	return &getArticleTagsDTO{
		slug: props.Slug,
	}, nil
}

type getArticleTagsOutputItem struct {
	Tags []string `dynamodbav:"tags"`
}

type listArticlesDTO struct {
	tag    string
	status StatusType
	limit  int32
}

type ListArticlesDTOProps struct {
	Tag    string
	Status string
	Limit  int32
}

func NewListArticlesDTO(props ListArticlesDTOProps) (*listArticlesDTO, error) {
	if props.Limit < 1 || props.Limit > 100 {
		return nil, &myerrors.ValidationError{Message: "limit must be between 1 and 100"}
	}

	if props.Status != string(StatusPublished) && props.Status != string(StatusUnpublished) {
		return nil, &myerrors.ValidationError{Message: "status must be either 'published' or 'unpublished'"}
	}

	return &listArticlesDTO{
		tag:    props.Tag,
		status: StatusType(props.Status),
		limit:  props.Limit,
	}, nil
}

type getArticleDTO struct {
	slug string
}

type GetArticleDTOProps struct {
	Slug string
}

func NewGetArticleDTO(props GetArticleDTOProps) (*getArticleDTO, error) {
	slugLength := len(props.Slug)
	if slugLength < slugMinLength || slugLength > slugMaxLength {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("slug must be between %d and %d characters", slugMinLength, slugMaxLength)}
	}

	return &getArticleDTO{
		slug: props.Slug,
	}, nil
}
