package usecase

import (
	"context"

	"devoratio.dev/web-resume/config"
	"devoratio.dev/web-resume/internal/errorx"
	"devoratio.dev/web-resume/internal/generator"
	"devoratio.dev/web-resume/model"
)

//go:generate mockgen -destination=usecasemock/authenticationmock.go -package=usecasemock . AuthenticationUsecase
type AuthenticationUsecase interface {
	Authenticate(ctx context.Context, identifier, password string) (*model.Owner, error)
}

type Login struct {
	authUsecase AuthenticationUsecase
	appConfig   *config.Application
}

func NewUsecase(authUsecase AuthenticationUsecase, appConfig *config.Application) *Login {
	return &Login{
		authUsecase: authUsecase,
		appConfig:   appConfig,
	}
}

func (l *Login) Login(ctx context.Context, identifier, password string) (string, error) {
	ownerAccount, err := l.authUsecase.Authenticate(ctx, identifier, password)
	if err != nil {
		return "", err
	}

	accessToken, err := generator.GenerateAccessToken(model.Claim{
		UserID:   ownerAccount.ID,
		Username: ownerAccount.Username,
	}, l.appConfig.Authentication.SigningKey)
	if err != nil {
		return "", errorx.ErrInternal
	}

	return accessToken, nil
}
