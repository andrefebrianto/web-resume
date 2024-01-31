package generator

import (
	"strings"
	"testing"
	"time"

	"devoratio.dev/web-resume/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestGenerateAccessToken(t *testing.T) {
	signature := []byte("random_sign_key")

	type args struct {
		claim      model.Claim
		signingKey []byte
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate with custom claim",
			args: args{
				claim: model.Claim{
					UserID:   168,
					Username: "devoratio",
				},
				signingKey: signature,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr: false,
		},
		{
			name: "generate with empty custom claim",
			args: args{
				claim:      model.Claim{},
				signingKey: signature,
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateAccessToken(tt.args.claim, tt.args.signingKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !strings.Contains(got, tt.want) {
				t.Errorf("GenerateAccessToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVerifyAccessToken(t *testing.T) {
	signature := []byte("random_sign_key")
	currentTime := time.Now()

	customClaimAccessToken, _ := GenerateAccessToken(model.Claim{UserID: 123, Username: "devoratio"}, signature)
	accessTokenWithInvalidSignature, _ := GenerateAccessToken(model.Claim{}, []byte("invalid_signature"))
	expiredAccessToken, _ := generateAccessTokenWithCustomClaim(CustomClaims{
		Data: model.Claim{},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(-24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(currentTime.Add(-48 * time.Hour)),
			NotBefore: jwt.NewNumericDate(currentTime.Add(-48 * time.Hour)),
			Issuer:    issuer,
			ID:        uuid.New().String(),
			Audience:  []string{audience},
		},
	}, signature)
	inactiveAccessToken, _ := generateAccessTokenWithCustomClaim(CustomClaims{
		Data: model.Claim{},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime.Add(1 * time.Hour)),
			Issuer:    issuer,
			ID:        uuid.New().String(),
			Audience:  []string{audience},
		},
	}, signature)

	type args struct {
		accessToken string
		signingKey  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "verify malformed token",
			args: args{
				accessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOjF9.8mAIDUfZNQT3TGm1QFIQp91OCpJpQpbB1-m9pA2mkHc",
				signingKey:  signature,
			},
			wantErr: true,
		},
		{
			name: "verify invalid signature token",
			args: args{
				accessToken: accessTokenWithInvalidSignature,
				signingKey:  signature,
			},
			wantErr: true,
		},
		{
			name: "verify expired token",
			args: args{
				accessToken: expiredAccessToken,
				signingKey:  signature,
			},
			wantErr: true,
		},
		{
			name: "verify inactive token",
			args: args{
				accessToken: inactiveAccessToken,
				signingKey:  signature,
			},
			wantErr: true,
		},
		{
			name: "verify valid custom claim token",
			args: args{
				accessToken: customClaimAccessToken,
				signingKey:  signature,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := VerifyAccessToken(tt.args.accessToken, tt.args.signingKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyAccessToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func generateAccessTokenWithCustomClaim(claim jwt.Claims, signingKey []byte) (string, error) {
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return tokenString.SignedString(signingKey)
}
