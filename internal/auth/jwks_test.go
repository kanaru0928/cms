package auth

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/lestrrat-go/jwx/v4/jwa"
	"github.com/lestrrat-go/jwx/v4/jwk"
	"github.com/lestrrat-go/jwx/v4/jwt"
)

// fakeJWKSFetcher は JWKSFetcher インターフェースを満たすテスト用のフェイクである。
type fakeJWKSFetcher struct {
	jwks jwk.Set
	err  error
}

func (f *fakeJWKSFetcher) GetJWKS(ctx context.Context) (jwk.Set, error) {
	return f.jwks, f.err
}

// newRSAKeyPairForTest は指定した kid を持つ RSA 鍵ペアを生成し、
// jwk.Key 形式の秘密鍵と公開鍵を返す。
func newRSAKeyPairForTest(t *testing.T, kid string) (priv jwk.Key, pub jwk.Key) {
	t.Helper()

	rawPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	privKey, err := jwk.Import[jwk.Key](rawPriv)
	if err != nil {
		t.Fatalf("failed to import RSA private key: %v", err)
	}
	if err := privKey.Set(jwk.KeyIDKey, kid); err != nil {
		t.Fatalf("failed to set kid on private key: %v", err)
	}

	pubKey, err := privKey.PublicKey()
	if err != nil {
		t.Fatalf("failed to derive public key: %v", err)
	}
	if err := pubKey.Set(jwk.AlgorithmKey, jwa.RS256()); err != nil {
		t.Fatalf("failed to set alg on public key: %v", err)
	}

	return privKey, pubKey
}

// signTokenForTest は与えられた秘密鍵で署名した JWT 文字列を生成する。
func signTokenForTest(t *testing.T, priv jwk.Key) string {
	t.Helper()

	token, err := jwt.NewBuilder().Subject("test-user").Build()
	if err != nil {
		t.Fatalf("failed to build token: %v", err)
	}

	signed, err := jwt.Sign(token, jwt.WithKey(jwa.RS256(), priv))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}

	return string(signed)
}

func TestCreateAuthenticationFunc(t *testing.T) {
	privA, pubA := newRSAKeyPairForTest(t, "key-a")
	privB, _ := newRSAKeyPairForTest(t, "key-b")

	jwksWithA := jwk.NewSet()
	if err := jwksWithA.AddKey(pubA); err != nil {
		t.Fatalf("failed to add key to JWKS: %v", err)
	}

	validToken := signTokenForTest(t, privA)
	tokenSignedByUnknownKey := signTokenForTest(t, privB)

	tests := []struct {
		name          string
		authHeader    string
		setAuthHeader bool
		fetcher       Fetcher
		wantErr       bool
		wantToken     bool
	}{
		{
			name:          "Authorizationヘッダーが無いとエラー",
			setAuthHeader: false,
			fetcher:       &fakeJWKSFetcher{jwks: jwksWithA},
			wantErr:       true,
		},
		{
			name:          "Bearerプレフィックスが無いとエラー",
			setAuthHeader: true,
			authHeader:    validToken,
			fetcher:       &fakeJWKSFetcher{jwks: jwksWithA},
			wantErr:       true,
		},
		{
			name:          "正常に検証できる",
			setAuthHeader: true,
			authHeader:    "Bearer " + validToken,
			fetcher:       &fakeJWKSFetcher{jwks: jwksWithA},
			wantErr:       false,
			wantToken:     true,
		},
		{
			name:          "JWKSに無い鍵で署名されたトークンはエラー",
			setAuthHeader: true,
			authHeader:    "Bearer " + tokenSignedByUnknownKey,
			fetcher:       &fakeJWKSFetcher{jwks: jwksWithA},
			wantErr:       true,
		},
		{
			name:          "JWKSの取得に失敗するとエラー",
			setAuthHeader: true,
			authHeader:    "Bearer " + validToken,
			fetcher:       &fakeJWKSFetcher{err: errors.New("failed to fetch JWKS")},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.setAuthHeader {
				req.Header.Set("Authorization", tt.authHeader)
			}

			input := &openapi3filter.AuthenticationInput{
				RequestValidationInput: &openapi3filter.RequestValidationInput{
					Request: req,
				},
			}

			authenticator := &jwksAuthenticator{fetcher: tt.fetcher}
			authFn := authenticator.CreateAuthenticationFunc()

			err := authFn(context.Background(), input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateAuthenticationFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantToken {
				if input.RequestValidationInput.Request.Context().Value("token") == nil {
					t.Errorf("expected token to be set in request context, but it was not")
				}
			}
		})
	}
}
