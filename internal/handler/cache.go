package handler

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

const (
	privateCacheControl = "private, no-store"
	publicCacheControl  = "public, max-age=0, s-maxage=30, stale-while-revalidate=60"
)

// 認証付きのレスポンスを共有キャッシュに保存しないよう、全てのレスポンスを既定で private にする
// 認証の失敗などでバリデーターがレスポンスを返す場合にも適用するため、バリデーターより前に配置する
func privateCacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		c.Response().Header().Set(echo.HeaderCacheControl, privateCacheControl)
		return next(c)
	}
}

// 未認証の読み取りは公開記事のみを返すため、共有キャッシュに保存させる
// 認証結果を参照するため、バリデーターより後に配置する
func (s *server) publicCacheControlMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		method := c.Request().Method
		_, authenticated := s.authenticator.GetTokenFromContext(c.Request().Context())
		if (method == http.MethodGet || method == http.MethodHead) && !authenticated {
			c.Response().Header().Set(echo.HeaderCacheControl, publicCacheControl)
		}
		return next(c)
	}
}
