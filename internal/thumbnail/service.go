package thumbnail

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kanaru0928/cms/internal/myerrors"
)

// CloudFront の /thumbnails/* と揃える
const keyPrefix = "thumbnails/"

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateUploadURL(ctx context.Context, props CreateUploadURLProps) (*CreateUploadURLOutput, error) {
	extension, ok := extensions[props.ContentType]
	if !ok {
		return nil, &myerrors.ValidationError{Message: fmt.Sprintf("unsupported content type: %s", props.ContentType)}
	}

	key := fmt.Sprintf("%s%s.%s", keyPrefix, uuid.NewString(), extension)

	uploadURL, err := s.repository.PresignPutThumbnail(ctx, key, props.ContentType)
	if err != nil {
		return nil, fmt.Errorf("failed to presign thumbnail: %w", err)
	}

	return &CreateUploadURLOutput{
		UploadURL: uploadURL,
		URL:       "/" + key,
	}, nil
}
