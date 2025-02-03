package mocks

import (
	"context"
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type TokenRepoMock struct {
	mock.Mock
}

func (r *TokenRepoMock) SaveToken(ctx context.Context, token *model.Token) error {
	args := r.Called(ctx, token)
	return args.Error(0)
}

func (r *TokenRepoMock) FindTokenByActivationToken(ctx context.Context, activationToken string) (*model.Token, error) {
	args := r.Called(ctx, activationToken)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TokenRepoMock) FindTokenByUserId(ctx context.Context, userId int64) (*model.Token, error) {
	args := r.Called(ctx, userId)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TokenRepoMock) DeleteToken(ctx context.Context, token *model.Token) error {
	args := r.Called(ctx, token)
	return args.Error(0)
}

func (r *TokenRepoMock) FindTokenByID(ctx context.Context, tokenID string) (*model.Token, error) {
	args := r.Called(ctx, tokenID)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TokenRepoMock) FindTokenByToken(ctx context.Context, token string) (*model.Token, error) {
	args := r.Called(ctx, token)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Token), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TokenRepoMock) ExistByToken(ctx context.Context, token string) bool {
	args := r.Called(ctx, token)
	return args.Bool(0)
}
