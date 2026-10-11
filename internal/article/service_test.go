package article

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/kanaru0928/cms/internal/content"
	"github.com/kanaru0928/cms/internal/metadata"
	"github.com/kanaru0928/cms/internal/myerrors"
	"go.uber.org/mock/gomock"
)

func newServiceForTest(t *testing.T) (*Service, *content.MockRepository, *metadata.MockRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	contentMock := content.NewMockRepository(ctrl)
	metadataMock := metadata.NewMockRepository(ctrl)
	return NewService(contentMock, metadataMock), contentMock, metadataMock
}

func TestListArticles(t *testing.T) {
	validProps := ListArticlesDTOProps{
		Tag:    "tag1",
		Status: string(metadata.StatusPublished),
		Limit:  10,
		Order:  string(metadata.SortOrderDesc),
	}

	t.Run("記事一覧を取得できる", func(t *testing.T) {
		s, _, metadataMock := newServiceForTest(t)

		wantDTO, err := metadata.NewListArticlesDTO(metadata.ListArticlesDTOProps(validProps))
		if err != nil {
			t.Fatalf("failed to create DTO: %v", err)
		}
		repoOutput := &metadata.ListArticlesOutput{
			Items: []metadata.ListArticlesOutputItem{
				{Slug: "test-article", Title: "Test Article"},
			},
		}
		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Eq(wantDTO)).Return(repoOutput, nil)

		got, err := s.ListArticles(context.Background(), validProps)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got.Items, repoOutput.Items) {
			t.Fatalf("expected items %v, got %v", repoOutput.Items, got.Items)
		}
	})

	t.Run("入力が不正な場合はリポジトリを呼ばずに ValidationError を返す", func(t *testing.T) {
		s, _, _ := newServiceForTest(t)

		props := validProps
		props.Limit = 0

		_, err := s.ListArticles(context.Background(), props)
		if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("リポジトリのエラーをそのまま返す", func(t *testing.T) {
		s, _, metadataMock := newServiceForTest(t)

		repoErr := errors.New("repository error")
		metadataMock.EXPECT().ListArticles(gomock.Any(), gomock.Any()).Return(nil, repoErr)

		_, err := s.ListArticles(context.Background(), validProps)
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestGetArticle(t *testing.T) {
	t.Run("記事を取得できる", func(t *testing.T) {
		s, _, metadataMock := newServiceForTest(t)

		wantDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps{Slug: "test-article"})
		if err != nil {
			t.Fatalf("failed to create DTO: %v", err)
		}
		repoOutput := &metadata.GetArticleOutput{Slug: "test-article", Title: "Test Article"}
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Eq(wantDTO)).Return(repoOutput, nil)

		got, err := s.GetArticle(context.Background(), GetArticleDTOProps{Slug: "test-article"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !reflect.DeepEqual(got, (*GetArticleOutput)(repoOutput)) {
			t.Fatalf("expected %v, got %v", repoOutput, got)
		}
	})

	t.Run("Slug が空の場合はリポジトリを呼ばずに ValidationError を返す", func(t *testing.T) {
		s, _, _ := newServiceForTest(t)

		_, err := s.GetArticle(context.Background(), GetArticleDTOProps{Slug: ""})
		if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("記事が存在しない場合は NotFoundError を返す", func(t *testing.T) {
		s, _, metadataMock := newServiceForTest(t)

		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(nil, &myerrors.NotFoundError{
			ID:           "test-article",
			ResourceType: "article",
		})

		_, err := s.GetArticle(context.Background(), GetArticleDTOProps{Slug: "test-article"})
		if !myerrors.HasSameTypeAs(err, &myerrors.NotFoundError{}) {
			t.Fatalf("expected NotFoundError, got %v", err)
		}
	})
}

func TestUpsertArticle(t *testing.T) {
	newProps := func() UpsertArticleProps {
		return UpsertArticleProps{
			Slug:         "test-article",
			Title:        "Test Article",
			ContentBody:  "content",
			Source:       "www.kanaru.me",
			Status:       string(metadata.StatusPublished),
			Tags:         []string{"tag1"},
			ThumbnailURL: "https://example.com/thumbnail.jpg",
		}
	}

	// md5("content") = 9a0364b9e99bb480dd25e1f0284c8555
	const wantContentKey = "articles/test-article/9a0364b9e99bb480dd25e1f0284c8555.md"

	for _, tt := range []struct {
		name    string
		created bool
	}{
		{name: "記事が存在しない場合は作成して Created が true になる", created: true},
		{name: "記事が存在する場合は更新して Created が false になる", created: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, contentMock, metadataMock := newServiceForTest(t)
			props := newProps()

			wantPutTextContentDTO, err := content.NewPutTextContentDTO(content.PutTextContentDTOProps{
				Key:     wantContentKey,
				Content: content.ContentType(props.ContentBody),
			})
			if err != nil {
				t.Fatalf("failed to create DTO: %v", err)
			}
			wantPutArticleDTO, err := metadata.NewPutArticleDTO(metadata.PutArticleDTOProps{
				Slug:         props.Slug,
				Title:        props.Title,
				Source:       props.Source,
				ContentKey:   wantContentKey,
				Status:       props.Status,
				Tags:         props.Tags,
				ThumbnailURL: props.ThumbnailURL,
			})
			if err != nil {
				t.Fatalf("failed to create DTO: %v", err)
			}
			wantGetArticleDTO, err := metadata.NewGetArticleDTO(metadata.GetArticleDTOProps{Slug: props.Slug})
			if err != nil {
				t.Fatalf("failed to create DTO: %v", err)
			}

			repoArticle := &metadata.GetArticleOutput{
				Slug:       props.Slug,
				Title:      props.Title,
				ContentKey: wantContentKey,
				CreatedAt:  "2026-01-01T00:00:00Z",
				PV:         42,
			}

			gomock.InOrder(
				contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Eq(*wantPutTextContentDTO)).Return(nil),
				metadataMock.EXPECT().PutArticle(gomock.Any(), gomock.Eq(wantPutArticleDTO)).Return(&metadata.PutArticleOutput{Created: tt.created}, nil),
				metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Eq(wantGetArticleDTO)).Return(repoArticle, nil),
			)

			got, err := s.UpsertArticle(context.Background(), props)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if got.Created != tt.created {
				t.Fatalf("expected Created %v, got %v", tt.created, got.Created)
			}
			if !reflect.DeepEqual(got.Article, (*GetArticleOutput)(repoArticle)) {
				t.Fatalf("expected article %v, got %v", repoArticle, got.Article)
			}
		})
	}

	for _, tt := range []struct {
		name   string
		modify func(p *UpsertArticleProps)
	}{
		{name: "本文が空の場合", modify: func(p *UpsertArticleProps) { p.ContentBody = "" }},
		{name: "Slug が空の場合", modify: func(p *UpsertArticleProps) { p.Slug = "" }},
		{name: "Title が空の場合", modify: func(p *UpsertArticleProps) { p.Title = "" }},
		{name: "Status が不正な場合", modify: func(p *UpsertArticleProps) { p.Status = "invalid" }},
		{name: "Tags が空の場合", modify: func(p *UpsertArticleProps) { p.Tags = []string{} }},
	} {
		t.Run(tt.name+"は書き込まずに ValidationError を返す", func(t *testing.T) {
			s, _, _ := newServiceForTest(t)
			props := newProps()
			tt.modify(&props)

			_, err := s.UpsertArticle(context.Background(), props)
			if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
				t.Fatalf("expected ValidationError, got %v", err)
			}
		})
	}

	t.Run("本文の保存に失敗した場合はメタデータを書き込まずにエラーを返す", func(t *testing.T) {
		s, contentMock, _ := newServiceForTest(t)

		repoErr := errors.New("content error")
		contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Any()).Return(repoErr)

		_, err := s.UpsertArticle(context.Background(), newProps())
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})

	t.Run("メタデータの保存に失敗した場合はエラーを返す", func(t *testing.T) {
		s, contentMock, metadataMock := newServiceForTest(t)

		repoErr := errors.New("metadata error")
		contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().PutArticle(gomock.Any(), gomock.Any()).Return(nil, repoErr)

		_, err := s.UpsertArticle(context.Background(), newProps())
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})

	t.Run("保存後の記事の取得に失敗した場合はエラーを返す", func(t *testing.T) {
		s, contentMock, metadataMock := newServiceForTest(t)

		repoErr := errors.New("get error")
		contentMock.EXPECT().PutTextContent(gomock.Any(), gomock.Any()).Return(nil)
		metadataMock.EXPECT().PutArticle(gomock.Any(), gomock.Any()).Return(&metadata.PutArticleOutput{Created: true}, nil)
		metadataMock.EXPECT().GetArticle(gomock.Any(), gomock.Any()).Return(nil, repoErr)

		_, err := s.UpsertArticle(context.Background(), newProps())
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}

func TestDeleteArticle(t *testing.T) {
	t.Run("メタデータを削除してから本文を全て削除する", func(t *testing.T) {
		s, contentMock, metadataMock := newServiceForTest(t)

		wantDeleteArticleDTO, err := metadata.NewDeleteArticleDTO(metadata.DeleteArticleDTOProps{Slug: "test-article"})
		if err != nil {
			t.Fatalf("failed to create DTO: %v", err)
		}
		wantDeleteContentsDTO, err := content.NewDeleteContentsByPrefixDTO(content.DeleteContentsByPrefixDTOProps{
			Prefix: "articles/test-article/",
		})
		if err != nil {
			t.Fatalf("failed to create DTO: %v", err)
		}

		gomock.InOrder(
			metadataMock.EXPECT().DeleteArticle(gomock.Any(), gomock.Eq(wantDeleteArticleDTO)).Return(nil),
			contentMock.EXPECT().DeleteContentsByPrefix(gomock.Any(), gomock.Eq(*wantDeleteContentsDTO)).Return(nil),
		)

		if err := s.DeleteArticle(context.Background(), DeleteArticleProps{Slug: "test-article"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("Slug が空の場合は削除せずに ValidationError を返す", func(t *testing.T) {
		s, _, _ := newServiceForTest(t)

		err := s.DeleteArticle(context.Background(), DeleteArticleProps{Slug: ""})
		if !myerrors.HasSameTypeAs(err, &myerrors.ValidationError{}) {
			t.Fatalf("expected ValidationError, got %v", err)
		}
	})

	t.Run("記事が存在しない場合は本文を削除せずに NotFoundError を返す", func(t *testing.T) {
		s, _, metadataMock := newServiceForTest(t)

		metadataMock.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(&myerrors.NotFoundError{
			ID:           "test-article",
			ResourceType: "Article",
		})

		err := s.DeleteArticle(context.Background(), DeleteArticleProps{Slug: "test-article"})
		if !myerrors.HasSameTypeAs(err, &myerrors.NotFoundError{}) {
			t.Fatalf("expected NotFoundError, got %v", err)
		}
	})

	t.Run("本文の削除に失敗した場合はエラーを返す", func(t *testing.T) {
		s, contentMock, metadataMock := newServiceForTest(t)

		repoErr := errors.New("content error")
		metadataMock.EXPECT().DeleteArticle(gomock.Any(), gomock.Any()).Return(nil)
		contentMock.EXPECT().DeleteContentsByPrefix(gomock.Any(), gomock.Any()).Return(repoErr)

		err := s.DeleteArticle(context.Background(), DeleteArticleProps{Slug: "test-article"})
		if !errors.Is(err, repoErr) {
			t.Fatalf("expected %v, got %v", repoErr, err)
		}
	})
}
