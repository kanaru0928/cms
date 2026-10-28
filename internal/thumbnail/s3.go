package thumbnail

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go/middleware"
)

const uploadURLExpires = 5 * time.Minute

type S3Repository struct {
	presignClient *s3.PresignClient
	bucketName    string
}

func NewS3Repository(cfg *aws.Config, bucketName string) *S3Repository {
	return &S3Repository{
		presignClient: s3.NewPresignClient(s3.NewFromConfig(*cfg)),
		bucketName:    bucketName,
	}
}

func (r *S3Repository) PresignPutThumbnail(ctx context.Context, key, contentType string) (string, error) {
	req, err := r.presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(uploadURLExpires), withSignedContentType)
	if err != nil {
		return "", fmt.Errorf("failed to presign put object: %w", err)
	}
	return req.URL, nil
}

// SDK は本文のない署名付き URL から Content-Type を除外するため、除外処理を外して Content-Type を署名に含める
// 署名に含めることで、指定した Content-Type 以外でのアップロードを拒否する
func withSignedContentType(o *s3.PresignOptions) {
	o.ClientOptions = append(o.ClientOptions, func(o *s3.Options) {
		o.APIOptions = append(o.APIOptions, func(stack *middleware.Stack) error {
			_, err := stack.Build.Remove("RemoveContentTypeHeader")
			return err
		})
	})
}
