package metadata

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func createTable(ctx context.Context, repository *DynamoDBRepository) error {
	_, err := repository.client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(repository.tableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String(string(articleKeySlug)), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(string(articleKeyItemType)), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(string(articleKeyStatus)), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(string(articleKeyFilterTag)), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String(string(articleKeyUpdatedAt)), AttributeType: types.ScalarAttributeTypeS},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String(string(articleKeySlug)), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String(string(articleKeyItemType)), KeyType: types.KeyTypeRange},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String("GSI_StatusTag_UpdatedAt"),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String(string(articleKeyStatus)), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String(string(articleKeyFilterTag)), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String(string(articleKeyUpdatedAt)), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{
					ProjectionType: types.ProjectionTypeInclude,
					NonKeyAttributes: []string{
						string(articleKeyTitle),
						string(articleKeyTags),
						string(articleKeySource),
						string(articleKeyContentKey),
						string(articleKeyCreatedAt),
					},
				},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	return err
}

func deleteTable(ctx context.Context, repository *DynamoDBRepository) error {
	_, err := repository.client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(repository.tableName),
	})
	return err
}

func TestPutArticle(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}
	repo := NewDynamoDBRepositoryForTest(&cfg, "articles", "http://localhost:8000")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      *PutArticleDTO
		wantErr    bool
		wantOutput []*articleItem
	}{
		{
			name: "単一タグの公開記事の登録に成功",
			input: &PutArticleDTO{
				Slug:       "test-article",
				Title:      "Test Article",
				Source:     "www.kanaru.me",
				ContentKey: "test-article/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1"},
					PV:         0,
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1"},
				},
			},
		},
		{
			name: "複数タグの記事の登録に成功",
			input: &PutArticleDTO{
				Slug:       "test-article-multiple-tags",
				Title:      "Test Article with Multiple Tags",
				Source:     "www.kanaru.me",
				ContentKey: "test-article-multiple-tags/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1", "tag2", "tag3"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article-multiple-tags",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article with Multiple Tags",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-multiple-tags/content.md",
					Tags:       []string{"tag1", "tag2", "tag3"},
					PV:         0,
				},
				{
					Slug:       "test-article-multiple-tags",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article with Multiple Tags",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-multiple-tags/content.md",
					Tags:       []string{"tag1", "tag2", "tag3"},
				},
				{
					Slug:       "test-article-multiple-tags",
					ItemType:   ItemTypeTag + "#tag2",
					Status:     StatusPublished,
					FilterTag:  "tag2",
					Title:      "Test Article with Multiple Tags",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-multiple-tags/content.md",
					Tags:       []string{"tag1", "tag2", "tag3"},
				},
				{
					Slug:       "test-article-multiple-tags",
					ItemType:   ItemTypeTag + "#tag3",
					Status:     StatusPublished,
					FilterTag:  "tag3",
					Title:      "Test Article with Multiple Tags",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-multiple-tags/content.md",
					Tags:       []string{"tag1", "tag2", "tag3"},
				},
			},
		},
		{
			name: "記事の追加登録に成功",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &PutArticleDTO{
					Slug:       "test-article",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Status:     StatusPublished,
					Tags:       []string{"tag1", "tag2"},
				})
			},
			input: &PutArticleDTO{
				Slug:       "test-article-additional",
				Title:      "Test Article Additional",
				Source:     "www.kanaru.me",
				ContentKey: "test-article-additional/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1", "tag3"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article-additional",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article Additional",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-additional/content.md",
					Tags:       []string{"tag1", "tag3"},
					PV:         0,
				},
				{
					Slug:       "test-article-additional",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article Additional",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-additional/content.md",
					Tags:       []string{"tag1", "tag3"},
				},
				{
					Slug:       "test-article-additional",
					ItemType:   ItemTypeTag + "#tag3",
					Status:     StatusPublished,
					FilterTag:  "tag3",
					Title:      "Test Article Additional",
					Source:     "www.kanaru.me",
					ContentKey: "test-article-additional/content.md",
					Tags:       []string{"tag1", "tag3"},
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
					PV:         0,
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag2",
					Status:     StatusPublished,
					FilterTag:  "tag2",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
			},
		},
		{
			name: "タグを変更しない記事の更新に成功",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &PutArticleDTO{
					Slug:       "test-article",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Status:     StatusPublished,
					Tags:       []string{"tag1", "tag2"},
				})
			},
			input: &PutArticleDTO{
				Slug:       "test-article",
				Title:      "Test Article Updated",
				Source:     "www.kanaru.me",
				ContentKey: "test-article/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1", "tag2"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article Updated",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
					PV:         0,
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article Updated",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag2",
					Status:     StatusPublished,
					FilterTag:  "tag2",
					Title:      "Test Article Updated",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
			},
		},
		{
			name: "タグを追加する記事の更新に成功",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &PutArticleDTO{
					Slug:       "test-article",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Status:     StatusPublished,
					Tags:       []string{"tag1"},
				})
			},
			input: &PutArticleDTO{
				Slug:       "test-article",
				Title:      "Test Article Updated with Additional Tag",
				Source:     "www.kanaru.me",
				ContentKey: "test-article/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1", "tag2"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article Updated with Additional Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
					PV:         0,
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article Updated with Additional Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag2",
					Status:	 StatusPublished,
					FilterTag:  "tag2",
					Title:      "Test Article Updated with Additional Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag2"},
				},
			},
		},
		{
			name: "タグを削除する記事の更新に成功",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &PutArticleDTO{
					Slug:       "test-article",
					Title:      "Test Article",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Status:     StatusPublished,
					Tags:       []string{"tag1", "tag2"},
				})
			},
			input: &PutArticleDTO{
				Slug:       "test-article",
				Title:      "Test Article Updated with Removed Tag",
				Source:     "www.kanaru.me",
				ContentKey: "test-article/content.md",
				Status:     StatusPublished,
				Tags:       []string{"tag1", "tag3"},
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:       "test-article",
					ItemType:   ItemTypeArticle,
					Status:     StatusPublished,
					FilterTag:  tagAll,
					Title:      "Test Article Updated with Removed Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag3"},
					PV:         0,
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag1",
					Status:     StatusPublished,
					FilterTag:  "tag1",
					Title:      "Test Article Updated with Removed Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag3"},
				},
				{
					Slug:       "test-article",
					ItemType:   ItemTypeTag + "#tag3",
					Status:     StatusPublished,
					FilterTag:  "tag3",
					Title:      "Test Article Updated with Removed Tag",
					Source:     "www.kanaru.me",
					ContentKey: "test-article/content.md",
					Tags:       []string{"tag1", "tag3"},
				},
			},
		},
		{
			name: "タグなしの記事の登録に失敗",
			input: &PutArticleDTO{
				Slug:       "test-article-no-tags",
				Title:      "Test Article with No Tags",
				Source:     "www.kanaru.me",
				ContentKey: "test-article-no-tags/content.md",
				Status:     StatusPublished,
				Tags:       []string{},
			},
			wantErr:    true,
			wantOutput: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := createTable(context.Background(), repo)
			if err != nil {
				t.Fatalf("Failed to create table: %v", err)
			}
			defer func() {
				if err := deleteTable(context.Background(), repo); err != nil {
					t.Fatalf("Failed to delete table: %v", err)
				}
			}()

			if tt.beforeFunc != nil {
				if err := tt.beforeFunc(); err != nil {
					t.Fatalf("Failed to run before function: %v", err)
				}
			}

			err = repo.PutArticle(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("PutArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			items, err := repo.ScanAllForTest(context.Background())
			if err != nil {
				t.Fatalf("Failed to scan all items: %v", err)
			}

			if tt.wantOutput != nil {
				if len(items) != len(tt.wantOutput) {
					t.Errorf("Expected %d items, but got %d", len(tt.wantOutput), len(items))
				} else {
					for i, wantItem := range tt.wantOutput {
						if !CompareArticleItems(items[i], wantItem) {
							t.Errorf("Expected item %d to be %+v, but got %+v", i, wantItem, items[i])
						}
					}
				}
			} else {
				if len(items) != 0 {
					t.Errorf("Expected no items, but got %d", len(items))
				}
			}
		})
	}
}

func CompareArticleItems(a, b *articleItem) bool {
	if a.Slug != b.Slug {
		return false
	}
	if a.ItemType != b.ItemType {
		return false
	}
	if a.Status != b.Status {
		return false
	}
	if a.FilterTag != b.FilterTag {
		return false
	}
	if a.Title != b.Title {
		return false
	}
	if a.Source != b.Source {
		return false
	}
	if a.ContentKey != b.ContentKey {
		return false
	}
	if len(a.Tags) != len(b.Tags) {
		return false
	}
	for i := range a.Tags {
		if a.Tags[i] != b.Tags[i] {
			return false
		}
	}
	if a.PV != b.PV {
		return false
	}

	return true
}
