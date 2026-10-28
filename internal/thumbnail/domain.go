package thumbnail

var extensions = map[string]string{
	"image/webp": "webp",
}

type CreateUploadURLProps struct {
	ContentType string
	// 画像の MD5 ハッシュ（16 進数）
	MD5 string
}

type CreateUploadURLOutput struct {
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
