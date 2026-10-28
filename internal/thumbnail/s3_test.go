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

	got, err := repo.PresignPutThumbnail(context.Background(), "thumbnails/test.png", "image/png")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	u, err := url.Parse(got)
	if err != nil {
		t.Fatalf("failed to parse url: %v", err)
	}
	if u.Path != "/thumbnails/test.png" {
		t.Errorf("Path = %v, want %v", u.Path, "/thumbnails/test.png")
	}
	query := u.Query()
	if query.Get("X-Amz-Expires") != "300" {
		t.Errorf("X-Amz-Expires = %v, want %v", query.Get("X-Amz-Expires"), "300")
	}
	if query.Get("X-Amz-SignedHeaders") != "content-type;host" {
		t.Errorf("X-Amz-SignedHeaders = %v, want %v", query.Get("X-Amz-SignedHeaders"), "content-type;host")
	}
}
