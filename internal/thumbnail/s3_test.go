package thumbnail

import (
	"context"
	"net/url"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const bucketName = "bucket"

func TestPresignPutThumbnail(t *testing.T) {
	cfg := aws.Config{
		Region:      "ap-northeast-1",
		Credentials: credentials.NewStaticCredentialsProvider("dummy-access-key", "dummy-secret-key", ""),
	}
	repo := NewS3Repository(&cfg, bucketName)

	got, err := repo.PresignPutThumbnail(context.Background(), PresignPutThumbnailProps{
		Key:         "thumbnails/test.webp",
		ContentType: "image/webp",
		ContentMD5:  testMD5Base64,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	if u.Path != "/thumbnails/test.webp" {
		t.Errorf("Path = %v, want %v", u.Path, "/thumbnails/test.webp")
	}
	query := u.Query()
	if query.Get("X-Amz-Expires") != "300" {
		t.Errorf("X-Amz-Expires = %v, want %v", query.Get("X-Amz-Expires"), "300")
	}
	if query.Get("X-Amz-SignedHeaders") != "content-md5;content-type;host" {
		t.Errorf("X-Amz-SignedHeaders = %v, want %v", query.Get("X-Amz-SignedHeaders"), "content-md5;content-type;host")
	}
}
