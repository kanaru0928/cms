package metadata

import "testing"

func TestNewPutArticleDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   PutArticleDTOProps
		wantErr bool
	}{
		{
			name: "すべてのバリデーションを通過",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: false,
		},
		{
			name: "Slug が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Slug が100文字を超えるとエラー",
			input: PutArticleDTOProps{
				Slug:         "this-is-a-very-long-slug-that-exceeds-the-maximum-length-of-one-hundred-characters-which-is-not-valid",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Title が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Title が200文字を超えるとエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "this-is-a-very-long-title-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid-this-is-a-very-long-title-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "ContentKey が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "ContentKey が200文字を超えるとエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "this-is-a-very-long-content-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid-this-is-a-very-long-content-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Source が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Status が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       "",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Tags が空ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Tags が20個を超えるとエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8", "tag9", "tag10", "tag11", "tag12", "tag13", "tag14", "tag15", "tag16", "tag17", "tag18", "tag19", "tag20", "tag21"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Tags の各タグが50文字を超えるとエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       string(StatusPublished),
				Tags:         []string{"this-is-a-very-long-tag-that-exceeds-the-maximum-length-of-fifty-characters-and-is-not-valid"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
		{
			name: "Status が published か unpublished 以外ならエラー",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       "invalid-status",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPutArticleDTO(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPutArticleDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.slug != tt.input.Slug {
					t.Errorf("NewPutArticleDTO() got.Slug = %v, want %v", got.slug, tt.input.Slug)
				}
				if got.title != tt.input.Title {
					t.Errorf("NewPutArticleDTO() got.Title = %v, want %v", got.title, tt.input.Title)
				}
				if got.contentKey != tt.input.ContentKey {
					t.Errorf("NewPutArticleDTO() got.ContentKey = %v, want %v", got.contentKey, tt.input.ContentKey)
				}
				if got.source != tt.input.Source {
					t.Errorf("NewPutArticleDTO() got.Source = %v, want %v", got.source, tt.input.Source)
				}
				if string(got.status) != tt.input.Status {
					t.Errorf("NewPutArticleDTO() got.Status = %v, want %v", got.status, tt.input.Status)
				}
				if len(got.tags) != len(tt.input.Tags) {
					t.Errorf("NewPutArticleDTO() got.Tags length = %v, want %v", len(got.tags), len(tt.input.Tags))
				} else {
					for i := range got.tags {
						if got.tags[i] != tt.input.Tags[i] {
							t.Errorf("NewPutArticleDTO() got.Tags[%d] = %v, want %v", i, got.tags[i], tt.input.Tags[i])
						}
					}
				}
				if got.thumbnailURL != tt.input.ThumbnailURL {
					t.Errorf("NewPutArticleDTO() got.ThumbnailURL = %v, want %v", got.thumbnailURL, tt.input.ThumbnailURL)
				}
			}
		})
	}
}

func TestNewListArticlesDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   ListArticlesDTOProps
		wantErr bool
	}{
		{
			name: "すべてのバリデーションを通過",
			input: ListArticlesDTOProps{
				Tag:    "valid-tag",
				Status: string(StatusPublished),
				Limit:  10,
			},
			wantErr: false,
		},
		{
			name: "Limit が1未満ならエラー",
			input: ListArticlesDTOProps{
				Tag:    "valid-tag",
				Status: string(StatusPublished),
				Limit:  0,
			},
			wantErr: true,
		},
		{
			name: "Limit が100を超えるとエラー",
			input: ListArticlesDTOProps{
				Tag:    "valid-tag",
				Status: string(StatusPublished),
				Limit:  101,
			},
			wantErr: true,
		},
		{
			name: "Status が published か unpublished 以外ならエラー",
			input: ListArticlesDTOProps{
				Tag:    "valid-tag",
				Status: "invalid-status",
				Limit:  10,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewListArticlesDTO(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewListArticlesDTO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.tag != tt.input.Tag {
					t.Errorf("NewListArticlesDTO() got.tag = %v, want %v", got.tag, tt.input.Tag)
				}
				if string(got.status) != tt.input.Status {
					t.Errorf("NewListArticlesDTO() got.status = %v, want %v", got.status, tt.input.Status)
				}
				if got.limit != tt.input.Limit {
					t.Errorf("NewListArticlesDTO() got.limit = %v, want %v", got.limit, tt.input.Limit)
				}
			}
		})
	}
}
