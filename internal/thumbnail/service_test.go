package thumbnail

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/kanaru0928/cms/internal/myerrors"
	"go.uber.org/mock/gomock"
)

func newServiceForTest(t *testing.T) (*Service, *MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repositoryMock := NewMockRepository(ctrl)
	return NewService(repositoryMock), repositoryMock
}

func TestCreateUploadURL(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		wantKey     *regexp.Regexp
	}{
		{name: "PNG のアップロード URL を発行できる", contentType: "image/png", wantKey: regexp.MustCompile(`^thumbnails/[0-9a-f-]{36}\.png$`)},
		{name: "JPEG のアップロード URL を発行できる", contentType: "image/jpeg", wantKey: regexp.MustCompile(`^thumbnails/[0-9a-f-]{36}\.jpg$`)},
		{name: "WebP のアップロード URL を発行できる", contentType: "image/webp", wantKey: regexp.MustCompile(`^thumbnails/[0-9a-f-]{36}\.webp$`)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, repositoryMock := newServiceForTest(t)

			var gotKey string
			repositoryMock.EXPECT().PresignPutThumbnail(gomock.Any(), gomock.Any(), tt.contentType).DoAndReturn(
				func(_ context.Context, key, _ string) (string, error) {
					gotKey = key
					return "https://example.com/upload", nil
				},
			)

			got, err := s.CreateUploadURL(context.Background(), CreateUploadURLProps{ContentType: tt.contentType})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if !tt.wantKey.MatchString(gotKey) {
				t.Fatalf("expected key to match %v, got %v", tt.wantKey, gotKey)
			}
			if got.UploadURL != "https://example.com/upload" {
				t.Fatalf("expected upload url %v, got %v", "https://example.com/upload", got.UploadURL)
			}
			if got.URL != "/"+gotKey {
				t.Fatalf("expected url %v, got %v", "/"+gotKey, got.URL)
			}
		})
	}

	t.Run("許可していない Content-Type の場合はリポジトリを呼ばずに ValidationError を返す", func(t *testing.T) {
		s, _ := newServiceForTest(t)

		_, err := s.CreateUploadURL(context.Background(), CreateUploadURLProps{ContentType: "image/gif"})
		if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("リポジトリのエラーをそのまま返す", func(t *testing.T) {
		s, repositoryMock := newServiceForTest(t)

		repoErr := errors.New("repository error")
		repositoryMock.EXPECT().PresignPutThumbnail(gomock.Any(), gomock.Any(), gomock.Any()).Return("", repoErr)

		_, err := s.CreateUploadURL(context.Background(), CreateUploadURLProps{ContentType: "image/png"})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
