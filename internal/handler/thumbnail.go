package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/myerrors"
	"github.com/kanaru0928/cms/internal/thumbnail"
)

// ThumbnailsCreateUploadUrl implements [api.StrictServerInterface].
func (s *server) ThumbnailsCreateUploadUrl(ctx context.Context, request api.ThumbnailsCreateUploadUrlRequestObject) (api.ThumbnailsCreateUploadUrlResponseObject, error) {
	output, err := s.thumbnailService.CreateUploadURL(ctx, thumbnail.CreateUploadURLProps{
		ContentType: string(request.Body.ContentType),
		MD5:         request.Body.Md5,
	})
	var validationErr *myerrors.ValidationError
	if errors.As(err, &validationErr) {
		return api.ThumbnailsCreateUploadUrl422JSONResponse{Error: validationErr.Message}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create thumbnail upload url: %w", err)
	}

	return api.ThumbnailsCreateUploadUrl200JSONResponse{
		UploadUrl: output.UploadURL,
		Url:       output.URL,
	}, nil
}
