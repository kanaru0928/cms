package auth

import (
	"context"

	"github.com/lestrrat-go/jwx/v4/jwk"
)

type Fetcher interface {
	GetJWKS(ctx context.Context) (jwk.Set, error)
}
