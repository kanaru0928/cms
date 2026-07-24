package metadata

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/kanaru0928/cms/internal/myerrors"
)

const gsiStatusTagUpdatedAt = "GSI_StatusTag_UpdatedAt"

type DynamoDBRepository struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDBRepository(cfg *aws.Config, tableName string) *DynamoDBRepository {
	client := dynamodb.NewFromConfig(*cfg)
	return &DynamoDBRepository{client: client, tableName: tableName}
}

// テストでは DynamoDB Local を使用するため、エンドポイントを指定できるようにする
func newDynamoDBRepositoryForTest(cfg *aws.Config, tableName string, endpoint string) *DynamoDBRepository {
	client := dynamodb.NewFromConfig(*cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	return &DynamoDBRepository{client: client, tableName: tableName}
}

func (r *DynamoDBRepository) GetArticleTags(ctx context.Context, getArticleTagsDTO *getArticleTagsDTO) ([]string, error) {
	key := articlePKMap{
		articlePKSlug:     &types.AttributeValueMemberS{Value: getArticleTagsDTO.slug},
		articlePKItemType: &types.AttributeValueMemberS{Value: string(ItemTypeArticle)},
	}

	items, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName:            aws.String(r.tableName),
		Key:                  key.ToStringMap(),
		ProjectionExpression: aws.String(string(articleKeyTags)),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get article tags: %w", err)
	}

	if items.Item == nil {
		return nil, &myerrors.NotFoundError{
			ID:           getArticleTagsDTO.slug,
			ResourceType: "Article",
		}
	}

	var outputItem getArticleTagsOutputItem
	err = attributevalue.UnmarshalMap(items.Item, &outputItem)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal article tags: %w", err)
	}

	return outputItem.Tags, nil
}

func (r *DynamoDBRepository) PutArticle(ctx context.Context, putArticleDTO *putArticleDTO) error {
	// 既存のタグを取得
	getArticleTagsDTO, err := NewGetArticleTagsDTO(GetArticleTagsDTOProps{
		Slug: putArticleDTO.slug,
	})
	if err != nil {
		return fmt.Errorf("failed to create GetArticleTagsDTO: %w", err)
	}

	existingTags, err := r.GetArticleTags(ctx, getArticleTagsDTO)
	var notFoundErr *myerrors.NotFoundError
	if err != nil && !errors.As(err, &notFoundErr) {
		return fmt.Errorf("failed to get article tags: %w", err)
	}

	// 孤立するタグアイテムは削除する必要がある
	var tagsToRemove []string
	if notFoundErr == nil {
		tagsToRemove = subtractTags(existingTags, putArticleDTO.tags)
	}

	// トランザクションを構築
	var transactoinItems []types.TransactWriteItem

	now := time.Now().Format(time.RFC3339)
	pk := articlePKMap{
		articlePKItemType: &types.AttributeValueMemberS{Value: string(ItemTypeArticle)},
		articlePKSlug:     &types.AttributeValueMemberS{Value: putArticleDTO.slug},
	}
	updateCommand := types.Update{
		TableName: aws.String(r.tableName),
		Key:       pk.ToStringMap(),
		UpdateExpression: aws.String(`
			SET #status = :status,
				#filter_tag = :filter_tag,
				#title = :title,
				#source = :source,
				#content_key = :content_key,
				#tags = :tags,
				#updated_at = :updated_at,
				#created_at = if_not_exists(#created_at, :created_at),
				#pv = if_not_exists(#pv, :pv),
				#thumbnail_url = :thumbnail_url
		`),
		ExpressionAttributeNames: map[string]string{
			"#status":        string(articleKeyStatus),
			"#filter_tag":    string(articleKeyFilterTag),
			"#title":         string(articleKeyTitle),
			"#source":        string(articleKeySource),
			"#content_key":   string(articleKeyContentKey),
			"#tags":          string(articleKeyTags),
			"#updated_at":    string(articleKeyUpdatedAt),
			"#created_at":    string(articleKeyCreatedAt),
			"#pv":            string(articleKeyPV),
			"#thumbnail_url": string(articleKeyThumbnailURL),
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":        &types.AttributeValueMemberS{Value: string(putArticleDTO.status)},
			":filter_tag":    &types.AttributeValueMemberS{Value: tagAll},
			":title":         &types.AttributeValueMemberS{Value: putArticleDTO.title},
			":source":        &types.AttributeValueMemberS{Value: putArticleDTO.source},
			":content_key":   &types.AttributeValueMemberS{Value: putArticleDTO.contentKey},
			":tags":          &types.AttributeValueMemberSS{Value: putArticleDTO.tags},
			":updated_at":    &types.AttributeValueMemberS{Value: now},
			":created_at":    &types.AttributeValueMemberS{Value: now},
			":pv":            &types.AttributeValueMemberN{Value: "0"},
			":thumbnail_url": &types.AttributeValueMemberS{Value: putArticleDTO.thumbnailURL},
		},
	}
	transactoinItems = append(transactoinItems, types.TransactWriteItem{
		Update: &updateCommand,
	})

	for _, tag := range putArticleDTO.tags {
		itemType := fmt.Sprintf("%s#%s", ItemTypeTag, tag)

		updateCommand := types.Update{
			TableName: aws.String(r.tableName),
			Key: map[string]types.AttributeValue{
				string(articlePKItemType): &types.AttributeValueMemberS{Value: itemType},
				string(articlePKSlug):     &types.AttributeValueMemberS{Value: putArticleDTO.slug},
			},
			UpdateExpression: aws.String(`
				SET #status = :status,
					#filter_tag = :filter_tag,
					#title = :title,
					#source = :source,
					#content_key = :content_key,
					#tags = :tags,
					#updated_at = :updated_at,
					#created_at = if_not_exists(#created_at, :created_at)
			`),
			ExpressionAttributeNames: map[string]string{
				"#status":      string(articleKeyStatus),
				"#filter_tag":  string(articleKeyFilterTag),
				"#title":       string(articleKeyTitle),
				"#source":      string(articleKeySource),
				"#content_key": string(articleKeyContentKey),
				"#tags":        string(articleKeyTags),
				"#updated_at":  string(articleKeyUpdatedAt),
				"#created_at":  string(articleKeyCreatedAt),
			},
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":status":      &types.AttributeValueMemberS{Value: string(putArticleDTO.status)},
				":filter_tag":  &types.AttributeValueMemberS{Value: tag},
				":title":       &types.AttributeValueMemberS{Value: putArticleDTO.title},
				":source":      &types.AttributeValueMemberS{Value: putArticleDTO.source},
				":content_key": &types.AttributeValueMemberS{Value: putArticleDTO.contentKey},
				":tags":        &types.AttributeValueMemberSS{Value: putArticleDTO.tags},
				":updated_at":  &types.AttributeValueMemberS{Value: now},
				":created_at":  &types.AttributeValueMemberS{Value: now},
			},
		}
		transactoinItems = append(transactoinItems, types.TransactWriteItem{
			Update: &updateCommand,
		})
	}

	for _, tag := range tagsToRemove {
		itemType := fmt.Sprintf("%s#%s", ItemTypeTag, tag)
		pk := articlePKMap{
			articlePKItemType: &types.AttributeValueMemberS{Value: itemType},
			articlePKSlug:     &types.AttributeValueMemberS{Value: putArticleDTO.slug},
		}
		deleteCommand := types.Delete{
			TableName: aws.String(r.tableName),
			Key:       pk.ToStringMap(),
		}
		transactoinItems = append(transactoinItems, types.TransactWriteItem{
			Delete: &deleteCommand,
		})
	}

	_, err = r.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: transactoinItems,
	})
	if err != nil {
		return fmt.Errorf("failed to put article: %w", err)
	}

	return nil
}

func subtractTags(existingTags, newTags []string) []string {
	newTagsSet := make(map[string]struct{})
	for _, tag := range newTags {
		newTagsSet[tag] = struct{}{}
	}
	tagsToRemove := make([]string, 0)
	for _, tag := range existingTags {
		if _, exists := newTagsSet[tag]; !exists {
			tagsToRemove = append(tagsToRemove, tag)
		}
	}

	return tagsToRemove
}

func (r *DynamoDBRepository) ListArticles(ctx context.Context, props *listArticlesDTO) (*listArticlesOutput, error) {
	var filterTag string
	if props.tag == "" {
		filterTag = tagAll
	} else {
		filterTag = props.tag
	}

	sortOrder := props.order == SortOrderAsc

	queryInput := &dynamodb.QueryInput{
		TableName:              aws.String(r.tableName),
		IndexName:              aws.String(gsiStatusTagUpdatedAt),
		KeyConditionExpression: aws.String("#status = :status AND #filter_tag = :filter_tag"),
		ExpressionAttributeNames: map[string]string{
			"#status":     string(articleKeyStatus),
			"#filter_tag": string(articleKeyFilterTag),
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":     &types.AttributeValueMemberS{Value: string(props.status)},
			":filter_tag": &types.AttributeValueMemberS{Value: filterTag},
		},
		ScanIndexForward: aws.Bool(sortOrder),
		Limit:            aws.Int32(props.limit),
	}

	queryOutput, err := r.client.Query(ctx, queryInput)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}

	var articles []*articleItem
	err = attributevalue.UnmarshalListOfMaps(queryOutput.Items, &articles)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal articles: %w", err)
	}

	lastKey := make(map[articlePK]types.AttributeValue)
	if queryOutput.LastEvaluatedKey != nil {
		lastKey = articlePKMap{
			articlePKSlug:     queryOutput.LastEvaluatedKey[string(articlePKSlug)],
			articlePKItemType: queryOutput.LastEvaluatedKey[string(articlePKItemType)],
		}
	}

	articleItems := make([]listArticlesOutputItem, len(articles))
	for i, article := range articles {
		articleItems[i] = listArticlesOutputItem{
			Slug:         article.Slug,
			Status:       article.Status,
			Title:        article.Title,
			Source:       article.Source,
			ContentKey:   article.ContentKey,
			Tags:         article.Tags,
			UpdatedAt:    article.UpdatedAt,
			CreatedAt:    article.CreatedAt,
			PV:           article.PV,
			ThumbnailURL: article.ThumbnailURL,
		}
	}

	return &listArticlesOutput{
		Items:            articleItems,
		lastEvaluatedKey: lastKey,
	}, nil
}

func (r *DynamoDBRepository) GetArticle(ctx context.Context, getArticleDTO *getArticleDTO) (*getArticleOutput, error) {
	key := articlePKMap{
		articlePKSlug:     &types.AttributeValueMemberS{Value: getArticleDTO.slug},
		articlePKItemType: &types.AttributeValueMemberS{Value: string(ItemTypeArticle)},
	}

	items, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(r.tableName),
		Key:       key.ToStringMap(),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	if items.Item == nil {
		return nil, &myerrors.NotFoundError{
			ID:           getArticleDTO.slug,
			ResourceType: "Article",
		}
	}

	var article articleItem
	err = attributevalue.UnmarshalMap(items.Item, &article)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal article: %w", err)
	}

	item := getArticleOutput{
		Slug:         article.Slug,
		Title:        article.Title,
		Source:       article.Source,
		ContentKey:   article.ContentKey,
		Status:       article.Status,
		Tags:         article.Tags,
		UpdatedAt:    article.UpdatedAt,
		CreatedAt:    article.CreatedAt,
		PV:           article.PV,
		ThumbnailURL: article.ThumbnailURL,
	}

	return &item, nil
}

// テスト用の関数であるため、実装での使用は不可
func (r *DynamoDBRepository) scanAllForTest(ctx context.Context) ([]*articleItem, error) {
	var articles []*articleItem
	scanPagenator := dynamodb.NewScanPaginator(r.client, &dynamodb.ScanInput{
		TableName: aws.String(r.tableName),
	})

	for scanPagenator.HasMorePages() {
		page, err := scanPagenator.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to scan items: %w", err)
		}

		var pageItems []*articleItem
		err = attributevalue.UnmarshalListOfMaps(page.Items, &pageItems)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal items: %w", err)
		}

		articles = append(articles, pageItems...)
	}

	return articles, nil
}
