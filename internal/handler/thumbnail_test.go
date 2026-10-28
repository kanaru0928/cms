package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kanaru0928/cms/internal/api"
	"github.com/kanaru0928/cms/internal/thumbnail"
	"github.com/labstack/echo/v5"
	"go.uber.org/mock/gomock"
)

func TestThumbnailsCreateUploadUrl(t *testing.T) {
	ctrl := gomock.NewController(t)
	e := echo.New()
	thumbnailMock := thumbnail.NewMockRepository(ctrl)

	newRequest := func(contentType api.ThumbnailContentType) api.ThumbnailsCreateUploadUrlRequestObject {
		return api.ThumbnailsCreateUploadUrlRequestObject{
			Body: &api.ThumbnailsCreateUploadUrlJSONRequestBody{
				ContentType: contentType,
				Md5:         "098f6bcd4621d373cade4e832627b4f6",
			},
		}
	}

	t.Run("アップロード URL を発行して 200 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/thumbnails/upload-url", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			thumbnailService: thumbnail.NewService(thumbnailMock),
		}

		thumbnailMock.EXPECT().PresignPutThumbnail(gomock.Any(), gomock.Any()).Return("https://example.com/upload", nil)

		response, err := s.ThumbnailsCreateUploadUrl(c.Request().Context(), newRequest(api.Imagewebp))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitThumbnailsCreateUploadUrlResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("許可していない Content-Type の場合は 422 が返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/thumbnails/upload-url", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			thumbnailService: thumbnail.NewService(thumbnailMock),
		}

		response, err := s.ThumbnailsCreateUploadUrl(c.Request().Context(), newRequest("image/png"))
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if err := response.VisitThumbnailsCreateUploadUrlResponse(rec); err != nil {
			t.Fatalf("expected response to be written, got %v", err)
		}

		if rec.Code != http.StatusUnprocessableEntity {
			t.Fatalf("expected status code %d, got %d", http.StatusUnprocessableEntity, rec.Code)
		}
	})

	t.Run("署名に失敗した場合はエラーが返る", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/thumbnails/upload-url", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		s := &server{
			thumbnailService: thumbnail.NewService(thumbnailMock),
		}

		thumbnailMock.EXPECT().PresignPutThumbnail(gomock.Any(), gomock.Any()).Return("", errors.New("presign error"))

		if _, err := s.ThumbnailsCreateUploadUrl(c.Request().Context(), newRequest(api.Imagewebp)); err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
