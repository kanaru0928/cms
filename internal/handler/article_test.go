package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/auth"
	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
	"github.com/labstack/echo/v5"
	"github.com/lestrrat-go/jwx/v4/jwt"
	"go.uber.org/mock/gomock"
)

func TestArticlesList(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	metadataMock := metadata.NewMockRepository(ctrl)
	contentMock := content.NewMockRepository(ctrl)
	authenticatorMock := auth.NewMockAuthenticator(ctrl)

	t.Run("認証なしで公開済みの記事を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}

		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&metadata.ListArticlesOutput{
			Items: []metadata.ListArticlesOutputItem{
				{Slug: "test-article", Title: "Test Article"},
			},
		}, nil)
		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)

		if err := s.ArticlesList(c, api.ArticlesListParams{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("認証なしで非公開記事を取得しようとすると 403 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles?status=unpublished", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)

		status := api.ArticleStatus(metadata.StatusUnpublished)

		if err := s.ArticlesList(c, api.ArticlesListParams{Status: &status}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusForbidden {
			t.Fatalf("expected status code %d, got %d", http.StatusForbidden, rec.Code)
		}
	})

	t.Run("認証ありで公開済みの記事を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}

		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&metadata.ListArticlesOutput{
			Items: []metadata.ListArticlesOutputItem{
				{Slug: "test-article", Title: "Test Article"},
			},
		}, nil)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(jwt.New(), true)

		if err := s.ArticlesList(c, api.ArticlesListParams{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("認証ありで非公開記事を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles?status=unpublished", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}

		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(&metadata.ListArticlesOutput{
			Items: []metadata.ListArticlesOutputItem{
				{Slug: "test-article", Title: "Test Article"},
			},
		}, nil)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(jwt.New(), true)

		status := api.ArticleStatus(metadata.StatusUnpublished)

		if err := s.ArticlesList(c, api.ArticlesListParams{Status: &status}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})
}
