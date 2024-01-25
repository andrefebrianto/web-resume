package usecase

import (
	"context"
	"errors"

	"devoratio.dev/web-resume/internal/errorx"
	"devoratio.dev/web-resume/internal/hasher"
	"devoratio.dev/web-resume/model"
)

const invalidInputMessage = "username or email or password is invalid"

//go:generate mockgen -destination=repositorymock/postgresqlmock.go -package=repositorymock . AuthenticationRepository
type AuthenticationRepository interface {
	GetOwnerByUsernameOrEmail(ctx context.Context, identifier string) (*model.OwnerAccount, error)
}

type Authentication struct {
	authRepo AuthenticationRepository
}

func NewUsecase(authRepo AuthenticationRepository) *Authentication {
	return &Authentication{
		authRepo: authRepo,
	}
}

func (a *Authentication) Authenticate(ctx context.Context, identifier, password string) (*model.Owner, error) {
	ownerAccount, err := a.authRepo.GetOwnerByUsernameOrEmail(ctx, identifier)
	if err != nil {
		if errors.Is(err, errorx.ErrNotFound) {
			return nil, errorx.New(errorx.TypeInvalidParameter, invalidInputMessage, err)
		}
		return nil, err
	}

	err = hasher.VerifyPassword(ownerAccount.Password, password)
	if err != nil {
		if errors.Is(err, errorx.ErrNotMatch) {
			return nil, errorx.New(errorx.TypeInvalidParameter, invalidInputMessage, err)
		}
		return nil, err
	}

	return &ownerAccount.Owner, nil
}
