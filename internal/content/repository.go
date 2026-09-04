//go:generate go tool mockgen -source=repository.go -destination=mock_repository.go -package=content

package content

import "context"

type Repository interface {
	PutTextContent(ctx context.Context, props putTextContentDTO) error
	GetTextContent(ctx context.Context, key string) (ContentType, error)
}
