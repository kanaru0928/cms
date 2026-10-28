package article

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
)

type Service struct {
	contentRepository  content.Repository
	metadataRepository metadata.Repository
}

func NewService(contentRepository content.Repository, metadataRepository metadata.Repository) *Service {
	return &Service{
		contentRepository:  contentRepository,
		metadataRepository: metadataRepository,
	}
}

func (s *Service) ListArticles(ctx context.Context, props ListArticlesDTOProps) (*ListArticlesOutput, error) {
	listArticlesDTO, err := metadata.NewListArticlesDTO(metadata.ListArticlesDTOProps(props))
	if err != nil {
		return nil, err
	}

	articles, err := s.metadataRepository.ListArticles(ctx, listArticlesDTO)
	if err != nil {
		return nil, err
	}

	return (*ListArticlesOutput)(articles), nil
}

func (s *Service) GetArticle(ctx context.Context, props GetArticleDTOProps) (*GetArticleOutput, error) {
	getArticleDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps(props))
	if err != nil {
		return nil, err
	}

	article, err := s.metadataRepository.GetArticle(ctx, getArticleDTO)
	if err != nil {
		return nil, err
	}
	return (*GetArticleOutput)(article), nil
}

func (s *Service) UpsertArticle(ctx context.Context, props UpsertArticleProps) (*UpsertArticleOutput, error) {
	contentHash := md5.Sum([]byte(props.ContentBody))
	contentKey := fmt.Sprintf("articles/%s/%s.md", props.Slug, hex.EncodeToString(contentHash[:]))

	// 書き込み前に全ての入力を検証する
	putTextContentDTO, err := content.NewPutTextContentDTO(content.PutTextContentDTOProps{
		Key:     contentKey,
		Content: content.ContentType(props.ContentBody),
	})
	if err != nil {
		return nil, err
	}

	putArticleDTO, err := metadata.NewPutArticleDTO(metadata.PutArticleDTOProps{
		Slug:         props.Slug,
		Title:        props.Title,
		Source:       props.Source,
		ContentKey:   contentKey,
		Status:       props.Status,
		Tags:         props.Tags,
		ThumbnailURL: props.ThumbnailURL,
	})
	if err != nil {
		return nil, err
	}

	getArticleDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps{
		Slug: props.Slug,
	})
	if err != nil {
		return nil, err
	}

	if err := s.contentRepository.PutTextContent(ctx, *putTextContentDTO); err != nil {
		return nil, fmt.Errorf("failed to put content: %w", err)
	}

	putArticleOutput, err := s.metadataRepository.PutArticle(ctx, putArticleDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to put article: %w", err)
	}

	// 作成日時や PV はリポジトリ側で決まるため、書き込み後に取得し直す
	article, err := s.metadataRepository.GetArticle(ctx, getArticleDTO)
	if err != nil {
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	return &UpsertArticleOutput{
		Created: putArticleOutput.Created,
		Article: (*GetArticleOutput)(article),
	}, nil
}

func (s *Service) DeleteArticle(ctx context.Context, props DeleteArticleProps) error {
	// 削除前に全ての入力を検証する
	deleteArticleDTO, err := metadata.NewDeleteArticleDTO(metadata.DeleteArticleDTOProps{
		Slug: props.Slug,
	})
	if err != nil {
		return err
	}

	deleteContentsDTO, err := content.NewDeleteContentsByPrefixDTO(content.DeleteContentsByPrefixDTOProps{
		Prefix: fmt.Sprintf("articles/%s/", props.Slug),
	})
	if err != nil {
		return err
	}

	// メタデータを先に消すと、本文の削除に失敗した場合に再実行しても 404 となり本文が残るため、本文から削除する
	if err := s.contentRepository.DeleteContentsByPrefix(ctx, *deleteContentsDTO); err != nil {
		return fmt.Errorf("failed to delete contents: %w", err)
	}

	if err := s.metadataRepository.DeleteArticle(ctx, deleteArticleDTO); err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	return nil
}
