package hasher

import (
	"devoratio.dev/web-resume/internal/errorx"
	"golang.org/x/crypto/bcrypt"
)

const hashCost = 12

func GenerateFromPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errorx.New(errorx.TypeUnauthorized, errorx.TypeUnauthorized.String(), err)
	}

	return nil
}
