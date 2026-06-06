package metadata

import "testing"

func TestNewPutArticleDTO(t *testing.T) {
	tests := []struct {
		name	string
		input	PutArticleDTOProps
		wantErr	bool
	}{
		{
			name: "すべてのバリデーションを通過",
			input: PutArticleDTOProps{
				Slug:         "valid-slug",
				Title:        "Valid Title",
				ContentKey:   "valid-content-key",
				Source:       "Valid Source",
				Status:       "draft",
				Tags:         []string{"tag1", "tag2"},
				ThumbnailURL: "http://example.com/thumbnail.jpg",
			},
			wantErr: false,
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
				if got.Slug != tt.input.Slug {
					t.Errorf("NewPutArticleDTO() got.Slug = %v, want %v", got.Slug, tt.input.Slug)
				}
				if got.Title != tt.input.Title {
					t.Errorf("NewPutArticleDTO() got.Title = %v, want %v", got.Title, tt.input.Title)
				}
				if got.ContentKey != tt.input.ContentKey {
					t.Errorf("NewPutArticleDTO() got.ContentKey = %v, want %v", got.ContentKey, tt.input.ContentKey)
				}
				if got.Source != tt.input.Source {
					t.Errorf("NewPutArticleDTO() got.Source = %v, want %v", got.Source, tt.input.Source)
				}
				if got.Status != tt.input.Status {
					t.Errorf("NewPutArticleDTO() got.Status = %v, want %v", got.Status, tt.input.Status)
				}
				if len(got.Tags) != len(tt.input.Tags) {
					t.Errorf("NewPutArticleDTO() got.Tags length = %v, want %v", len(got.Tags), len(tt.input.Tags))
				} else {
					for i := range got.Tags {
						if got.Tags[i] != tt.input.Tags[i] {
							t.Errorf("NewPutArticleDTO() got.Tags[%d] = %v, want %v", i, got.Tags[i], tt.input.Tags[i])
						}
					}
				}
				if got.ThumbnailURL != tt.input.ThumbnailURL {
					t.Errorf("NewPutArticleDTO() got.ThumbnailURL = %v, want %v", got.ThumbnailURL, tt.input.ThumbnailURL)
				}
			}
		})
	}
}
