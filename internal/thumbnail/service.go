package thumbnail

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"

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

	sum, err := hex.DecodeString(props.MD5)
	if err != nil || len(sum) != md5.Size {
		return nil, &myerrors.ValidationError{Message: "md5 must be a 32-character hex string"}
	}

	key := fmt.Sprintf("%s%s.%s", keyPrefix, hex.EncodeToString(sum), extension)

	uploadURL, err := s.repository.PresignPutThumbnail(ctx, PresignPutThumbnailProps{
		Key:         key,
		ContentType: props.ContentType,
		// S3 に本文と MD5 の一致を検証させ、キーと異なる画像のアップロードを拒否する
		ContentMD5: base64.StdEncoding.EncodeToString(sum),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to presign thumbnail: %w", err)
	}

	return &CreateUploadURLOutput{
		UploadURL: uploadURL,
		URL:       "/" + key,
	}, nil
}
