package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanaru0928/cms/internal/auth"
	"github.com/labstack/echo/v5"
	"github.com/lestrrat-go/jwx/v4/jwt"
	"go.uber.org/mock/gomock"
)

func TestCacheControlMiddleware(t *testing.T) {
	for _, tt := range []struct {
		name          string
		method        string
		authenticated bool
		want          string
	}{
		{name: "未認証の GET は public になる", method: http.MethodGet, authenticated: false, want: publicCacheControl},
		{name: "認証ありの GET は private になる", method: http.MethodGet, authenticated: true, want: privateCacheControl},
		{name: "PUT は private になる", method: http.MethodPut, authenticated: true, want: privateCacheControl},
		{name: "DELETE は private になる", method: http.MethodDelete, authenticated: true, want: privateCacheControl},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			authenticatorMock := auth.NewMockAuthenticator(ctrl)
			if tt.authenticated {
				authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(jwt.New(), true).AnyTimes()
			} else {
				authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false).AnyTimes()
			}
			s := &server{authenticator: authenticatorMock}

			e := echo.New()
			e.Use(privateCacheControlMiddleware, s.publicCacheControlMiddleware)
			e.Add(tt.method, "/api/articles", func(c *echo.Context) error {
				return c.NoContent(http.StatusOK)
			})

			req := httptest.NewRequest(tt.method, "/api/articles", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if got := rec.Header().Get(echo.HeaderCacheControl); got != tt.want {
				t.Fatalf("expected Cache-Control %q, got %q", tt.want, got)
			}
		})
	}

	t.Run("後続でエラーレスポンスを返しても private になる", func(t *testing.T) {
		e := echo.New()
		e.Use(privateCacheControlMiddleware)
		e.GET("/api/articles", func(c *echo.Context) error {
			return echo.ErrUnauthorized
		})

		req := httptest.NewRequest(http.MethodGet, "/api/articles", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		if got := rec.Header().Get(echo.HeaderCacheControl); got != privateCacheControl {
			t.Fatalf("expected Cache-Control %q, got %q", privateCacheControl, got)
		}
	})
}
