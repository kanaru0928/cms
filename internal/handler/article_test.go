package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
	"github.com/labstack/echo/v5"
	"github.com/lestrrat-go/jwx/v4/jwt"
	"go.uber.org/mock/gomock"
)

// mockAuthenticator は認証なしの状態を再現するフェイク実装。
type mockAuthenticator struct{}

func (m *mockAuthenticator) CreateAuthenticationFunc() openapi3filter.AuthenticationFunc {
	return nil
}

func (m *mockAuthenticator) GetTokenFromContext(ctx context.Context) (jwt.Token, bool) {
	return nil, false
}

func TestArticlesList(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	metadataMock := metadata.NewMockRepository(ctrl)
	contentMock := content.NewMockRepository(ctrl)

	t.Run("認証なしで公開済みの記事を取得できる", func(t *testing.T) {
		// Setup
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      &mockAuthenticator{},
		}

		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&metadata.ListArticlesOutput{
			Items: []metadata.ListArticlesOutputItem{
				{Slug: "test-article", Title: "Test Article"},
			},
		}, nil)

		if err := s.ArticlesList(c, api.ArticlesListParams{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
