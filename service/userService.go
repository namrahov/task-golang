package service

import (
	"context"
	"fmt"
	"task-golang/model"
	"task-golang/repo"
)

type IUserService interface {
	Register(ctx context.Context, userRegistrationDto *model.UserRegistrationDto) error
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func (us *UserService) Register(ctx context.Context, userRegistrationDto *model.UserRegistrationDto) error {

	fmt.Println(userRegistrationDto)
	return nil
}
