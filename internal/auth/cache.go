package auth

import (
	"context"
	"fmt"

	"github.com/jwx-go/jwkfetch/v4"
	"github.com/lestrrat-go/jwx/v4/jwk"
)

type jwksCache struct {
	jwksURL string
	jwks    jwk.Set
}

func NewJWKSCache(jwksURL string) *jwksCache {

	return &jwksCache{
		jwksURL: jwksURL,
	}
}

func (c *jwksCache) GetJWKSURL(ctx context.Context) (jwk.Set, error) {
	if c.jwks == nil {
		jwks, err := fetchJWKS(ctx, c.jwksURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
		}
		c.jwks = jwks
	}
	return c.jwks, nil
}

func fetchJWKS(ctx context.Context, jwksURL string) (jwk.Set, error) {
	client := jwkfetch.NewClient()
	jwks, err := client.Fetch(ctx, jwksURL)
	if err != nil {
		return nil, err
	}
	return jwks, nil
}
