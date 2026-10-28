package thumbnail

// アップロードを許可する Content-Type と、オブジェクトキーに付与する拡張子の対応
var extensions = map[string]string{
	"image/webp": "webp",
}

type CreateUploadURLProps struct {
	ContentType string
	// 画像の MD5 ハッシュ（16 進数）
	MD5 string
}

type CreateUploadURLOutput struct {
	// アップロード先の署名付き URL
	UploadURL string
	// アップロード後に配信されるパス
	URL string
}

type PresignPutThumbnailProps struct {
	Key         string
	ContentType string
	// Content-MD5 ヘッダーの値（Base64）
	ContentMD5 string
}
