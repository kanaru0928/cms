package metadata

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kanaru0928/cms/internal/myerrors"
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
				IndexName: aws.String(gsiStatusTagUpdatedAt),
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
	repo := newDynamoDBRepositoryForTest(&cfg, "articles", "http://localhost:8000")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      *putArticleDTO
		wantErr    bool
		wantOutput []*articleItem
	}{
		{
			name: "単一タグの公開記事の登録に成功",
			input: &putArticleDTO{
				slug:         "test-article",
				title:        "Test Article",
				source:       "www.kanaru.me",
				contentKey:   "test-article/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1"},
				thumbnailURL: "https://example.com/thumbnail.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article/content.md",
					Tags:         []string{"tag1"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail.jpg",
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
			input: &putArticleDTO{
				slug:         "test-article-multiple-tags",
				title:        "Test Article with Multiple Tags",
				source:       "www.kanaru.me",
				contentKey:   "test-article-multiple-tags/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1", "tag2", "tag3"},
				thumbnailURL: "https://example.com/thumbnail.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article-multiple-tags",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article with Multiple Tags",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article-multiple-tags/content.md",
					Tags:         []string{"tag1", "tag2", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail.jpg",
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
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "test-article",
					title:        "Test Article",
					source:       "www.kanaru.me",
					contentKey:   "test-article/content.md",
					status:       StatusPublished,
					tags:         []string{"tag1", "tag2"},
					thumbnailURL: "https://example.com/thumbnail1.jpg",
				})
			},
			input: &putArticleDTO{
				slug:         "test-article-additional",
				title:        "Test Article Additional",
				source:       "www.kanaru.me",
				contentKey:   "test-article-additional/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1", "tag3"},
				thumbnailURL: "https://example.com/thumbnail2.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article-additional",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article Additional",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article-additional/content.md",
					Tags:         []string{"tag1", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail2.jpg",
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
					Slug:         "test-article",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail1.jpg",
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
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "test-article",
					title:        "Test Article",
					source:       "www.kanaru.me",
					contentKey:   "test-article/content.md",
					status:       StatusPublished,
					tags:         []string{"tag1", "tag2"},
					thumbnailURL: "https://example.com/thumbnail.jpg",
				})
			},
			input: &putArticleDTO{
				slug:         "test-article",
				title:        "Test Article Updated",
				source:       "www.kanaru.me",
				contentKey:   "test-article/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1", "tag2"},
				thumbnailURL: "https://example.com/thumbnail.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article Updated",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail.jpg",
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
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "test-article",
					title:        "Test Article",
					source:       "www.kanaru.me",
					contentKey:   "test-article/content.md",
					status:       StatusPublished,
					tags:         []string{"tag1"},
					thumbnailURL: "https://example.com/thumbnail.jpg",
				})
			},
			input: &putArticleDTO{
				slug:         "test-article",
				title:        "Test Article Updated with Additional Tag",
				source:       "www.kanaru.me",
				contentKey:   "test-article/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1", "tag2"},
				thumbnailURL: "https://example.com/thumbnail.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article Updated with Additional Tag",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail.jpg",
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
					Status:     StatusPublished,
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
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "test-article",
					title:        "Test Article",
					source:       "www.kanaru.me",
					contentKey:   "test-article/content.md",
					status:       StatusPublished,
					tags:         []string{"tag1", "tag2"},
					thumbnailURL: "https://example.com/thumbnail.jpg",
				})
			},
			input: &putArticleDTO{
				slug:         "test-article",
				title:        "Test Article Updated with Removed Tag",
				source:       "www.kanaru.me",
				contentKey:   "test-article/content.md",
				status:       StatusPublished,
				tags:         []string{"tag1", "tag3"},
				thumbnailURL: "https://example.com/thumbnail.jpg",
			},
			wantErr: false,
			wantOutput: []*articleItem{
				{
					Slug:         "test-article",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Test Article Updated with Removed Tag",
					Source:       "www.kanaru.me",
					ContentKey:   "test-article/content.md",
					Tags:         []string{"tag1", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail.jpg",
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
			input: &putArticleDTO{
				slug:         "test-article-no-tags",
				title:        "Test Article with No Tags",
				source:       "www.kanaru.me",
				contentKey:   "test-article-no-tags/content.md",
				status:       StatusPublished,
				tags:         []string{},
				thumbnailURL: "https://example.com/thumbnail.jpg",
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

			items, err := repo.scanAllForTest(context.Background())
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

func TestGetArticleTags(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}
	repo := newDynamoDBRepositoryForTest(&cfg, "articles", "http://localhost:8000")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      *getArticleTagsDTO
		wantErr    error
		want       []string
	}{
		{
			name: "記事が存在する場合、タグを取得できる",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "existing-article",
					title:        "Existing Article",
					source:       "www.kanaru.me",
					contentKey:   "existing-article/content.md",
					status:       StatusPublished,
					tags:         []string{"tag1", "tag2"},
					thumbnailURL: "https://example.com/thumbnail.jpg",
				})
			},
			input: &getArticleTagsDTO{
				slug: "existing-article",
			},
			wantErr: nil,
			want:    []string{"tag1", "tag2"},
		},
		{
			name:       "記事が存在しない場合、NotFound を返す",
			beforeFunc: nil,
			input: &getArticleTagsDTO{
				slug: "non-existing-article",
			},
			wantErr: &myerrors.NotFoundError{},
			want:    nil,
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

			got, err := repo.GetArticleTags(context.Background(), tt.input)
			if (err != nil) && tt.wantErr != nil {
				if !myerrors.HasSameTypeAs(err, tt.wantErr) {
					t.Errorf("GetArticleTags() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("GetArticleTags() got length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("GetArticleTags() got[%d] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestListArticles(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}
	repo := newDynamoDBRepositoryForTest(&cfg, "articles", "http://localhost:8000")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      *listArticlesDTO
		wantErr    bool
		want       []*articleItem
	}{
		{
			name: "複数記事が存在する場合、全ての記事を取得できる",
			beforeFunc: func() error {
				articles := []putArticleDTO{
					{
						slug:         "article-1",
						title:        "Article 1",
						source:       "www.kanaru.me",
						contentKey:   "article-1/content.md",
						status:       StatusPublished,
						tags:         []string{"tag1", "tag2"},
						thumbnailURL: "https://example.com/thumbnail1.jpg",
					},
					{
						slug:         "article-2",
						title:        "Article 2",
						source:       "www.kanaru.me",
						contentKey:   "article-2/content.md",
						status:       StatusPublished,
						tags:         []string{"tag2", "tag3"},
						thumbnailURL: "https://example.com/thumbnail2.jpg",
					},
				}

				for _, article := range articles {
					if err := repo.PutArticle(context.Background(), &article); err != nil {
						return err
					}
				}
				return nil
			},
			input: &listArticlesDTO{
				status: StatusPublished,
				tag:    "",
				limit:  10,
			},
			wantErr: false,
			want: []*articleItem{
				{
					Slug:         "article-1",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Article 1",
					Source:       "www.kanaru.me",
					ContentKey:   "article-1/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail1.jpg",
				},
				{
					Slug:         "article-2",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Article 2",
					Source:       "www.kanaru.me",
					ContentKey:   "article-2/content.md",
					Tags:         []string{"tag2", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail2.jpg",
				},
			},
		},
		{
			name: "タグでフィルタリングして記事を取得できる",
			beforeFunc: func() error {
				articles := []putArticleDTO{
					{
						slug:         "article-1",
						title:        "Article 1",
						source:       "www.kanaru.me",
						contentKey:   "article-1/content.md",
						status:       StatusPublished,
						tags:         []string{"tag1", "tag2"},
						thumbnailURL: "https://example.com/thumbnail1.jpg",
					},
					{
						slug:         "article-2",
						title:        "Article 2",
						source:       "www.kanaru.me",
						contentKey:   "article-2/content.md",
						status:       StatusPublished,
						tags:         []string{"tag2", "tag3"},
						thumbnailURL: "https://example.com/thumbnail2.jpg",
					},
					{
						slug:         "article-3",
						title:        "Article 3",
						source:       "www.kanaru.me",
						contentKey:   "article-3/content.md",
						status:       StatusPublished,
						tags:         []string{"tag3", "tag4"},
						thumbnailURL: "https://example.com/thumbnail3.jpg",
					},
				}

				for _, article := range articles {
					if err := repo.PutArticle(context.Background(), &article); err != nil {
						return err
					}
				}
				return nil
			},
			input: &listArticlesDTO{
				status: StatusPublished,
				tag:    "tag2",
				limit:  10,
			},
			wantErr: false,
			want: []*articleItem{
				{
					Slug:         "article-1",
					ItemType:     ItemTypeTag + "#tag2",
					Status:       StatusPublished,
					FilterTag:    "tag2",
					Title:        "Article 1",
					Source:       "www.kanaru.me",
					ContentKey:   "article-1/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail1.jpg",
				},
				{
					Slug:         "article-2",
					ItemType:     ItemTypeTag + "#tag2",
					Status:       StatusPublished,
					FilterTag:    "tag2",
					Title:        "Article 2",
					Source:       "www.kanaru.me",
					ContentKey:   "article-2/content.md",
					Tags:         []string{"tag2", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail2.jpg",
				},
			},
		},
		{
			name:       "記事が存在しない場合、空のリストを返す",
			beforeFunc: nil,
			input: &listArticlesDTO{
				status: StatusPublished,
				tag:    "",
				limit:  10,
			},
			wantErr: false,
			want:    []*articleItem{},
		},
		{
			name: "公開の記事のみを取得できる",
			beforeFunc: func() error {
				articles := []putArticleDTO{
					{
						slug:         "article-1",
						title:        "Article 1",
						source:       "www.kanaru.me",
						contentKey:   "article-1/content.md",
						status:       StatusPublished,
						tags:         []string{"tag1", "tag2"},
						thumbnailURL: "https://example.com/thumbnail1.jpg",
					},
					{
						slug:         "article-2",
						title:        "Article 2",
						source:       "www.kanaru.me",
						contentKey:   "article-2/content.md",
						status:       StatusUnpublished,
						tags:         []string{"tag2", "tag3"},
						thumbnailURL: "https://example.com/thumbnail2.jpg",
					},
					{
						slug:         "article-3",
						title:        "Article 3",
						source:       "www.kanaru.me",
						contentKey:   "article-3/content.md",
						status:       StatusPublished,
						tags:         []string{"tag3", "tag4"},
						thumbnailURL: "https://example.com/thumbnail3.jpg",
					},
				}

				for _, article := range articles {
					if err := repo.PutArticle(context.Background(), &article); err != nil {
						return err
					}
				}
				return nil
			},
			input: &listArticlesDTO{
				status: StatusPublished,
				tag:    "",
				limit:  10,
			},
			wantErr: false,
			want: []*articleItem{
				{
					Slug:         "article-1",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Article 1",
					Source:       "www.kanaru.me",
					ContentKey:   "article-1/content.md",
					Tags:         []string{"tag1", "tag2"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail1.jpg",
				},
				{
					Slug:         "article-3",
					ItemType:     ItemTypeArticle,
					Status:       StatusPublished,
					FilterTag:    tagAll,
					Title:        "Article 3",
					Source:       "www.kanaru.me",
					ContentKey:   "article-3/content.md",
					Tags:         []string{"tag3", "tag4"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail3.jpg",
				},
			},
		},
		{
			name: "非公開の記事のみを取得できる",
			beforeFunc: func() error {
				articles := []putArticleDTO{
					{
						slug:         "article-1",
						title:        "Article 1",
						source:       "www.kanaru.me",
						contentKey:   "article-1/content.md",
						status:       StatusPublished,
						tags:         []string{"tag1", "tag2"},
						thumbnailURL: "https://example.com/thumbnail1.jpg",
					},
					{
						slug:         "article-2",
						title:        "Article 2",
						source:       "www.kanaru.me",
						contentKey:   "article-2/content.md",
						status:       StatusUnpublished,
						tags:         []string{"tag2", "tag3"},
						thumbnailURL: "https://example.com/thumbnail2.jpg",
					},
					{
						slug:         "article-3",
						title:        "Article 3",
						source:       "www.kanaru.me",
						contentKey:   "article-3/content.md",
						status:       StatusUnpublished,
						tags:         []string{"tag3", "tag4"},
						thumbnailURL: "https://example.com/thumbnail3.jpg",
					},
				}

				for _, article := range articles {
					if err := repo.PutArticle(context.Background(), &article); err != nil {
						return err
					}
				}
				return nil
			},
			input: &listArticlesDTO{
				status: StatusUnpublished,
				tag:    "",
				limit:  10,
			},
			wantErr: false,
			want: []*articleItem{
				{
					Slug:         "article-2",
					ItemType:     ItemTypeArticle,
					Status:       StatusUnpublished,
					FilterTag:    tagAll,
					Title:        "Article 2",
					Source:       "www.kanaru.me",
					ContentKey:   "article-2/content.md",
					Tags:         []string{"tag2", "tag3"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail2.jpg",
				},
				{
					Slug:         "article-3",
					ItemType:     ItemTypeArticle,
					Status:       StatusUnpublished,
					FilterTag:    tagAll,
					Title:        "Article 3",
					Source:       "www.kanaru.me",
					ContentKey:   "article-3/content.md",
					Tags:         []string{"tag3", "tag4"},
					PV:           0,
					ThumbnailURL: "https://example.com/thumbnail3.jpg",
				},
			},
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

			got, _, err := repo.ListArticles(context.Background(), *tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ListArticles() got length = %v, want %v", len(got), len(tt.want))
			}
			for i := range got {
				if !CompareArticleItems(got[i], tt.want[i]) {
					t.Errorf("ListArticles() got[%d] = %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetArticle(t *testing.T) {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("dummy", "dummy", "")),
	)
	if err != nil {
		t.Fatalf("Failed to load AWS config: %v", err)
	}
	repo := newDynamoDBRepositoryForTest(&cfg, "articles", "http://localhost:8000")

	tests := []struct {
		name       string
		beforeFunc func() error
		input      *getArticleDTO
		wantErr    error
		want       *articleItem
	}{
		{
			name: "記事が存在する場合、記事を取得できる",
			beforeFunc: func() error {
				return repo.PutArticle(context.Background(), &putArticleDTO{
					slug:         "article-1",
					status:       StatusPublished,
					title:        "Article 1",
					source:       "www.kanaru.me",
					contentKey:   "article-1/content.md",
					tags:         []string{"tag1", "tag2"},
					thumbnailURL: "https://example.com/thumbnail1.jpg",
				})
			},
			input: &getArticleDTO{
				slug: "article-1",
			},
			wantErr: nil,
			want: &articleItem{
				Slug:         "article-1",
				ItemType:     ItemTypeArticle,
				Status:       StatusPublished,
				FilterTag:    tagAll,
				Title:        "Article 1",
				Source:       "www.kanaru.me",
				ContentKey:   "article-1/content.md",
				Tags:         []string{"tag1", "tag2"},
				PV:           0,
				ThumbnailURL: "https://example.com/thumbnail1.jpg",
			},
		},
		{
			name:       "記事が存在しない場合、NotFound を返す",
			beforeFunc: nil,
			input: &getArticleDTO{
				slug: "non-existing-article",
			},
			wantErr: &myerrors.NotFoundError{},
			want:    nil,
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

			got, err := repo.GetArticle(context.Background(), tt.input)
			if (err != nil) && tt.wantErr != nil {
				if !myerrors.HasSameTypeAs(err, tt.wantErr) {
					t.Errorf("GetArticle() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if !CompareArticleItems(got, tt.want) {
				t.Errorf("GetArticle() got = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func CompareArticleItems(a, b *articleItem) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
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
