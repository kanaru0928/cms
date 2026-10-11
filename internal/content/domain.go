package content

import "github.com/kanaru0928/cms/internal/myerrors"

type ContentType string

type PutTextContentDTOProps struct {
	Key     string
	Content ContentType
}

type putTextContentDTO struct {
	key     string
	content ContentType
}

const (
	keyMinLength         = 1
	keyMaxLength         = 200
	textContentMinLength = 1
	textContentMaxLength = 100000
)

func NewPutTextContentDTO(props PutTextContentDTOProps) (*putTextContentDTO, error) {
	keyLength := len(props.Key)
	if keyLength < keyMinLength || keyLength > keyMaxLength {
		return nil, &myerrors.ValidationError{Message: "key must be between 1 and 200 characters"}
	}

	contentLength := len(props.Content)
	if contentLength < textContentMinLength || contentLength > textContentMaxLength {
		return nil, &myerrors.ValidationError{Message: "content must be between 1 and 100000 characters"}
	}

	return &putTextContentDTO{
		key:     props.Key,
		content: props.Content,
	}, nil
}

type getTextContentDTO struct {
	key string
}

type GetTextContentDTOProps struct {
	Key string
}

func NewGetTextContentDTO(props GetTextContentDTOProps) (*getTextContentDTO, error) {
	keyLength := len(props.Key)
	if keyLength < keyMinLength || keyLength > keyMaxLength {
		return nil, &myerrors.ValidationError{Message: "key must be between 1 and 200 characters"}
	}

	return &getTextContentDTO{
		key: props.Key,
	}, nil
}

type deleteContentsByPrefixDTO struct {
	prefix string
}

type DeleteContentsByPrefixDTOProps struct {
	Prefix string
}

func NewDeleteContentsByPrefixDTO(props DeleteContentsByPrefixDTOProps) (*deleteContentsByPrefixDTO, error) {
	// 空の prefix はリポジトリ配下の全てのコンテンツに一致するため許可しない
	prefixLength := len(props.Prefix)
	if prefixLength < keyMinLength || prefixLength > keyMaxLength {
		return nil, &myerrors.ValidationError{Message: "prefix must be between 1 and 200 characters"}
	}

	return &deleteContentsByPrefixDTO{
		prefix: props.Prefix,
	}, nil
}
