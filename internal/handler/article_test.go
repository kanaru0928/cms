package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/article"
	"github.com/kanaru0928/cms/internal/auth"
	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
	"github.com/kanaru0928/cms/internal/myerrors"
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

		if _, err := s.ArticlesList(c.Request().Context(), api.ArticlesListRequestObject{}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("認証なしで非公開記事を取得しようとすると 401 が返る", func(t *testing.T) {
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

		response, err := s.ArticlesList(c.Request().Context(), api.ArticlesListRequestObject{Params: api.ArticlesListParams{Status: &status}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesListResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnauthorized {
			t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, rec.Code)
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

		if _, err := s.ArticlesList(c.Request().Context(), api.ArticlesListRequestObject{}); err != nil {
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

		if _, err := s.ArticlesList(c.Request().Context(), api.ArticlesListRequestObject{Params: api.ArticlesListParams{Status: &status}}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("入力が不正な場合は 422 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles?limit=0", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)

		limit := int32(0)

		response, err := s.ArticlesList(c.Request().Context(), api.ArticlesListRequestObject{Params: api.ArticlesListParams{Limit: &limit}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesListResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

func TestArticlesRead(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	metadataMock := metadata.NewMockRepository(ctrl)
	contentMock := content.NewMockRepository(ctrl)
	authenticatorMock := auth.NewMockAuthenticator(ctrl)

	newServer := func() *server {
		return &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
		}
	}

	t.Run("認証なしで公開済みの記事を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&metadata.GetArticleOutput{
			Slug: "test-article", Status: metadata.StatusPublished, ContentKey: "articles/test-article/hash.md",
		}, nil)
		contentMock.EXPECT().GetTextContent(gomock.Any(), "articles/test-article/hash.md").Return(content.ContentType("content"), nil)

		response, err := newServer().ArticlesRead(c.Request().Context(), api.ArticlesReadRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesReadResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("認証なしで非公開記事を取得しようとすると 404 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&metadata.GetArticleOutput{
			Slug: "test-article", Status: metadata.StatusUnpublished,
		}, nil)

		response, err := newServer().ArticlesRead(c.Request().Context(), api.ArticlesReadRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesReadResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("認証ありで非公開記事を取得できる", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(jwt.New(), true)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&metadata.GetArticleOutput{
			Slug: "test-article", Status: metadata.StatusUnpublished, ContentKey: "articles/test-article/hash.md",
		}, nil)
		contentMock.EXPECT().GetTextContent(gomock.Any(), "articles/test-article/hash.md").Return(content.ContentType("content"), nil)

		response, err := newServer().ArticlesRead(c.Request().Context(), api.ArticlesReadRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesReadResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("記事が存在しない場合は 404 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(nil, &myerrors.NotFoundError{
			ID:           "test-article",
			ResourceType: "Article",
		})

		response, err := newServer().ArticlesRead(c.Request().Context(), api.ArticlesReadRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesReadResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("入力が不正な場合は 422 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/articles/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authenticatorMock.EXPECT().GetTokenFromContext(gomock.Any()).Return(nil, false)

		response, err := newServer().ArticlesRead(c.Request().Context(), api.ArticlesReadRequestObject{Slug: ""})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesReadResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

func TestArticlesUpsert(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	metadataMock := metadata.NewMockRepository(ctrl)
	contentMock := content.NewMockRepository(ctrl)
	authenticatorMock := auth.NewMockAuthenticator(ctrl)

	newRequest := func() api.ArticlesUpsertRequestObject {
		return api.ArticlesUpsertRequestObject{
			Slug: "test-article",
			Body: &api.ArticlesUpsertJSONRequestBody{
				Title:        "Test Article",
				ContentBody:  "content",
				Source:       "www.kanaru.me",
				Status:       api.ArticleStatus(metadata.StatusPublished),
				Tags:         []string{"tag1"},
				ThumbnailUrl: "https://example.com/thumbnail.jpg",
			},
		}
	}

	t.Run("記事が存在しない場合は作成して 201 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
			articleService:     article.NewService(contentMock, metadataMock),
		}

		contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().PutArticle(gomock.Any(), gomock.Any()).Return(&metadata.PutArticleOutput{Created: true}, nil)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&metadata.GetArticleOutput{
			Slug: "test-article", Title: "Test Article",
		}, nil)

		response, err := s.ArticlesUpsert(c.Request().Context(), newRequest())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesUpsertResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected status code %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	t.Run("記事が存在する場合は更新して 200 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
			articleService:     article.NewService(contentMock, metadataMock),
		}

		contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().PutArticle(gomock.Any(), gomock.Any()).Return(&metadata.PutArticleOutput{Created: false}, nil)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(&metadata.GetArticleOutput{
			Slug: "test-article", Title: "Test Article",
		}, nil)

		response, err := s.ArticlesUpsert(c.Request().Context(), newRequest())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesUpsertResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("入力が不正な場合は書き込まずに 422 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
			articleService:     article.NewService(contentMock, metadataMock),
		}

		request := newRequest()
		request.Body.Tags = []string{}

		response, err := s.ArticlesUpsert(c.Request().Context(), request)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesUpsertResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}

func TestArticlesDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	metadataMock := metadata.NewMockRepository(ctrl)
	contentMock := content.NewMockRepository(ctrl)
	authenticatorMock := auth.NewMockAuthenticator(ctrl)

	newServer := func() *server {
		return &server{
			metadataRepository: metadataMock,
			contentRepository:  contentMock,
			authenticator:      authenticatorMock,
			articleService:     article.NewService(contentMock, metadataMock),
		}
	}

	t.Run("記事を削除して 200 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		contentMock.EXPECT().DeleteContentsByPrefix(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(nil)

		response, err := newServer().ArticlesDelete(c.Request().Context(), api.ArticlesDeleteRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesDeleteResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("記事が存在しない場合は 404 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/articles/test-article", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		contentMock.EXPECT().DeleteContentsByPrefix(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(&myerrors.NotFoundError{
			ID:           "test-article",
			ResourceType: "Article",
		})

		response, err := newServer().ArticlesDelete(c.Request().Context(), api.ArticlesDeleteRequestObject{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesDeleteResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})

	t.Run("入力が不正な場合は削除せずに 422 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/articles/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		response, err := newServer().ArticlesDelete(c.Request().Context(), api.ArticlesDeleteRequestObject{Slug: ""})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitArticlesDeleteResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})
}
