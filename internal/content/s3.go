package content

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
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
