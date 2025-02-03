package mocks

import (
	"context"
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type UserUtilMock struct {
	mock.Mock
}

func (u *UserUtilMock) GetUserFromRequest(ctx context.Context) (*model.User, *model.ErrorResponse) {
	args := u.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Get(1).(*model.ErrorResponse)
	}
	return nil, args.Get(1).(*model.ErrorResponse)
}
