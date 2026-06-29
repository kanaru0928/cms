package content

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

const bucketName = "bucket"

func TestPutTextContent(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("garage"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"dummy-access-key",
			"dummy-secret-key",
			"",
		)),
	)

	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	repo := newS3RepositoryForTest(&cfg, bucketName, "test-prefix/", "http://localhost:3900", "garage")

	tests := []struct {
		name         string
		input        putTextContentDTO
		wantErr      bool
		wantContents map[string]ContentType
	}{
		{
			name: "正常にアップロードできる",
			input: putTextContentDTO{
				key:     "test-key",
				content: "This is a test content.",
			},
			wantErr: false,
			wantContents: map[string]ContentType{
				"test-prefix/test-key": "This is a test content.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.PutTextContent(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutTextContent() error = %v, wantErr %v", err, tt.wantErr)
			}
			contents, err := repo.getAllContentsForTest(context.Background())
			if err != nil {
				t.Fatalf("failed to get all contents: %v", err)
			}
			if !compareTextContents(tt.wantContents, contents) {
				t.Errorf("PutTextContent() got = %v, want %v", contents, tt.wantContents)
			}
		})
	}
}

func TestGetTextContent(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("garage"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"dummy-access-key",
			"dummy-secret-key",
			"",
		)),
	)

	if err != nil {
		t.Fatalf("failed to load AWS config: %v", err)
	}

	repo := newS3RepositoryForTest(&cfg, bucketName, "test-prefix/", "http://localhost:3900", "garage")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      string
		want       ContentType
		wantErr    bool
	}{
		{
			name: "正常に取得できる",
			beforeFunc: func() error {
				return repo.PutTextContent(context.Background(), putTextContentDTO{
					key:     "test-key",
					content: "This is a test content.",
				})
			},
			input:   "test-key",
			want:    "This is a test content.",
			wantErr: false,
		},
		{
			name:       "存在しないキーを指定するとエラー",
			beforeFunc: func() error { return nil },
			input:      "non-existent-key",
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.beforeFunc(); err != nil {
				t.Fatalf("beforeFunc() error = %v", err)
			}
			got, err := repo.GetTextContent(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTextContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetTextContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareTextContents(a map[string]ContentType, b map[string][]byte) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if string(v) != string(b[k]) {
			return false
		}
	}
	return true
}
