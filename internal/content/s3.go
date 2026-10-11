package content

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/kanaru0928/cms/internal/myerrors"
)

type S3Repository struct {
	client     *s3.Client
	bucketName string
	prefix     string
}

func NewS3Repository(cfg *aws.Config, bucketName, prefix string) *S3Repository {
	client := s3.NewFromConfig(*cfg)
	return &S3Repository{
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}
}

func newS3RepositoryForTest(cfg *aws.Config, bucketName, prefix, endpoint, region string) *S3Repository {
	client := s3.NewFromConfig(*cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
		o.Region = region
	})
	return &S3Repository{
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}
}

func (r *S3Repository) getFullKey(key string) string {
	return r.prefix + key
}

func (r *S3Repository) PutTextContent(ctx context.Context, props putTextContentDTO) error {
	r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(r.getFullKey(props.key)),
		Body:   strings.NewReader(string(props.content)),
	})
	return nil
}

func (r *S3Repository) GetTextContent(ctx context.Context, key string) (ContentType, error) {
	resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(r.getFullKey(key)),
	})
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", &myerrors.NotFoundError{
			ID:           key,
			ResourceType: "content",
		}
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return ContentType(buf), nil
}

func (r *S3Repository) DeleteContentsByPrefix(ctx context.Context, props deleteContentsByPrefixDTO) error {
	// ListObjectsV2 は 1 ページ最大 1000 件で、DeleteObjects の上限と一致する
	paginator := s3.NewListObjectsV2Paginator(r.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucketName),
		Prefix: aws.String(r.getFullKey(props.prefix)),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list contents: %w", err)
		}
		if len(page.Contents) == 0 {
			continue
		}

		objects := make([]types.ObjectIdentifier, len(page.Contents))
		for i, obj := range page.Contents {
			objects[i] = types.ObjectIdentifier{Key: obj.Key}
		}

		output, err := r.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
			Bucket: aws.String(r.bucketName),
			Delete: &types.Delete{
				Objects: objects,
				Quiet:   aws.Bool(true),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to delete contents: %w", err)
		}
		if len(output.Errors) > 0 {
			return fmt.Errorf("failed to delete %d contents: %s", len(output.Errors), aws.ToString(output.Errors[0].Message))
		}
	}

	return nil
}

func (r *S3Repository) getAllContentsForTest(ctx context.Context) (map[string][]byte, error) {
	keys, err := r.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(r.bucketName),
		Prefix: aws.String(r.prefix),
	})
	if err != nil {
		return nil, err
	}

	contents := make(map[string][]byte)
	for _, obj := range keys.Contents {
		resp, err := r.client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(r.bucketName),
			Key:    obj.Key,
		})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		buf, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		contents[*obj.Key] = buf
	}
	return contents, nil
}
