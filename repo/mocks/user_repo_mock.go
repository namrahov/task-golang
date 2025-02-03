package mocks

import (
	"gorm.io/gorm"
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (r *UserRepoMock) FindUserById(id int64) (*model.User, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserRepoMock) GetUserByEmail(email string) (*model.User, error) {
	args := r.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserRepoMock) SaveUser(tx *gorm.DB, user *model.User) (*model.User, error) {
	args := r.Called(tx, user)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserRepoMock) UpdateUser(user *model.User) (*model.User, error) {
	args := r.Called(user)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserRepoMock) AddRolesToUser(tx *gorm.DB, userId int64, roles []*model.Role) error {
	args := r.Called(tx, userId, roles)
	return args.Error(0)
}

func (r *UserRepoMock) FindActiveUserByEmailOrUsername(emailOrNickname string) (*model.User, error) {
	args := r.Called(emailOrNickname)
	if args.Get(0) != nil {
		return args.Get(0).(*model.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *UserRepoMock) GetPermissions(roles []string) ([]model.Permission, error) {
	args := r.Called(roles)
	if args.Get(0) != nil {
		return args.Get(0).([]model.Permission), args.Error(1)
	}
	return nil, args.Error(1)
}
