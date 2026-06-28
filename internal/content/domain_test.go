package content

import "testing"

func TestNewPutTextContentDTO(t *testing.T) {
	veryLongContent := make([]byte, 100001)

	tests := []struct {
		name    string
		input   PutTextContentDTOProps
		wantErr bool
	}{
		{
			name: "すべてのバリデーションを通過",
			input: PutTextContentDTOProps{
				Key:     "valid-key",
				Content: "Valid content",
			},
			wantErr: false,
		},
		{
			name: "Key が空ならエラー",
			input: PutTextContentDTOProps{
				Key:     "",
				Content: "Valid content",
			},
			wantErr: true,
		},
		{
			name: "Key が200文字を超えるとエラー",
			input: PutTextContentDTOProps{
				Key:     "this-is-a-very-long-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid-this-is-a-very-long-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid",
				Content: "Valid content",
			},
			wantErr: true,
		},
		{
			name: "Content が空ならエラー",
			input: PutTextContentDTOProps{
				Key:     "valid-key",
				Content: "",
			},
			wantErr: true,
		},
		{
			name: "Content が100000文字を超えるとエラー",
			input: PutTextContentDTOProps{
				Key:     "valid-key",
				Content: ContentType(veryLongContent),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPutTextContentDTO(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPutTextContentDTO() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewGetTextContentDTO(t *testing.T) {
	tests := []struct {
		name    string
		input   GetTextContentDTOProps
		wantErr bool
	}{
		{
			name: "すべてのバリデーションを通過",
			input: GetTextContentDTOProps{
				Key: "valid-key",
			},
			wantErr: false,
		},
		{
			name: "Key が空ならエラー",
			input: GetTextContentDTOProps{
				Key: "",
			},
			wantErr: true,
		},
		{
			name: "Key が200文字を超えるとエラー",
			input: GetTextContentDTOProps{
				Key: "this-is-a-very-long-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid-this-is-a-very-long-key-that-exceeds-the-maximum-length-of-two-hundred-characters-which-is-not-valid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewGetTextContentDTO(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGetTextContentDTO() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
