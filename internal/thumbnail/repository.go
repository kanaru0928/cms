//go:generate go tool mockgen -source=repository.go -destination=mock_repository.go -package=thumbnail

package thumbnail

import "context"

type Repository interface {
	PresignPutThumbnail(ctx context.Context, props PresignPutThumbnailProps) (string, error)
}
