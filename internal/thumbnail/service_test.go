package thumbnail

import (
	"context"
	"errors"
	"testing"

	"github.com/kanaru0928/cms/internal/myerrors"
	"go.uber.org/mock/gomock"
)

// "test" の MD5
const (
	testMD5Hex    = "098f6bcd4621d373cade4e832627b4f6"
	testMD5Base64 = "CY9rzUYh03PK3k6DJie09g=="
)

func newServiceForTest(t *testing.T) (*Service, *MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repositoryMock := NewMockRepository(ctrl)
	return NewService(repositoryMock), repositoryMock
}

func TestCreateUploadURL(t *testing.T) {
	t.Run("MD5 をキーとしたアップロード URL を発行できる", func(t *testing.T) {
		s, repositoryMock := newServiceForTest(t)

		repositoryMock.EXPECT().PresignPutThumbnail(gomock.Any(), PresignPutThumbnailProps{
			Key:         "thumbnails/" + testMD5Hex + ".webp",
			ContentType: "image/webp",
			ContentMD5:  testMD5Base64,
		}).Return("https://example.com/upload", nil)

		got, err := s.CreateUploadURL(context.Background(), CreateUploadURLProps{ContentType: "image/webp", MD5: testMD5Hex})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if got.UploadURL != "https://example.com/upload" {
			t.Fatalf("expected upload url %v, got %v", "https://example.com/upload", got.UploadURL)
		}
		if got.URL != "/thumbnails/"+testMD5Hex+".webp" {
			t.Fatalf("expected url %v, got %v", "/thumbnails/"+testMD5Hex+".webp", got.URL)
		}
	})

	invalidTests := []struct {
		name  string
		props CreateUploadURLProps
	}{
		{name: "許可していない Content-Type", props: CreateUploadURLProps{ContentType: "image/png", MD5: testMD5Hex}},
		{name: "MD5 が 16 進数でない", props: CreateUploadURLProps{ContentType: "image/webp", MD5: "zz8f6bcd4621d373cade4e832627b4f6"}},
		{name: "MD5 の長さが不正", props: CreateUploadURLProps{ContentType: "image/webp", MD5: "098f6bcd"}},
	}
	for _, tt := range invalidTests {
		t.Run(tt.name+"の場合はリポジトリを呼ばずに ValidationError を返す", func(t *testing.T) {
			s, _ := newServiceForTest(t)

			_, err := s.CreateUploadURL(context.Background(), tt.props)
			if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
				t.Fatalf("expected ValidationError, got %v", err)
			}
		})
	}

	t.Run("リポジトリのエラーをそのまま返す", func(t *testing.T) {
		s, repositoryMock := newServiceForTest(t)

		repoErr := errors.New("repository error")
		repositoryMock.EXPECT().PresignPutThumbnail(gomock.Any(), gomock.Any()).Return("", repoErr)

		_, err := s.CreateUploadURL(context.Background(), CreateUploadURLProps{ContentType: "image/webp", MD5: testMD5Hex})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
