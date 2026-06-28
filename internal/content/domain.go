package content

import "github.com/kanaru0928/cms/internal/myerrors"

type ContentType string

type PutTextContentDTOProps struct {
	Key     string
	Content ContentType
}

type PutTextContentDTO struct {
	key     string
	content ContentType
}

const (
	keyMinLength         = 1
	keyMaxLength         = 200
	textContentMinLength = 1
	textContentMaxLength = 100000
)

func NewPutTextContentDTO(props PutTextContentDTOProps) (*PutTextContentDTO, error) {
	keyLength := len(props.Key)
	if keyLength < keyMinLength || keyLength > keyMaxLength {
		return nil, &myerrors.ValidationError{Message: "key must be between 1 and 200 characters"}
	}

	contentLength := len(props.Content)
	if contentLength < textContentMinLength || contentLength > textContentMaxLength {
		return nil, &myerrors.ValidationError{Message: "content must be between 1 and 100000 characters"}
	}

	return &PutTextContentDTO{
		key:     props.Key,
		content: props.Content,
	}, nil
}

type GetTextContentDTO struct {
	key string
}

type GetTextContentDTOProps struct {
	Key string
}

func NewGetTextContentDTO(props GetTextContentDTOProps) (*GetTextContentDTO, error) {
	keyLength := len(props.Key)
	if keyLength < keyMinLength || keyLength > keyMaxLength {
		return nil, &myerrors.ValidationError{Message: "key must be between 1 and 200 characters"}
	}

	return &GetTextContentDTO{
		key: props.Key,
	}, nil
}
