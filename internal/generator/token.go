package generator

import (
	"time"

	"devoratio.dev/web-resume/internal/errorx"
	"devoratio.dev/web-resume/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	issuer   = "web-resume"
	audience = "web-resume"
)

type CustomClaims struct {
	Data model.Claim `json:"data"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(claim model.Claim, signingKey []byte) (string, error) {
	currentTime := time.Now()
	claims := CustomClaims{
		claim,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			NotBefore: jwt.NewNumericDate(currentTime),
			Issuer:    issuer,
			ID:        uuid.New().String(),
			Audience:  []string{audience},
		},
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenString.SignedString(signingKey)
}

func VerifyAccessToken(accessToken string, signingKey []byte) (*model.Claim, error) {
	token, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if !token.Valid || err != nil {
		return nil, errorx.ErrUnauthorized
	}

	customClaim := token.Claims.(*CustomClaims)
	return &customClaim.Data, nil
}
