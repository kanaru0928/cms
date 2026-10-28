package thumbnail

// アップロードを許可する Content-Type と、オブジェクトキーに付与する拡張子の対応
var extensions = map[string]string{
	"image/png":  "png",
	"image/jpeg": "jpg",
	"image/webp": "webp",
}

type CreateUploadURLProps struct {
	ContentType string
}

type CreateUploadURLOutput struct {
	// アップロード先の署名付き URL
	UploadURL string
	// アップロード後に配信されるパス
	URL string
}
