package initializer

import (
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

func InitUserService() *service.UserService {
	return &service.UserService{
		UserRepo:        &repo.UserRepo{},
		TokenRepo:       &repo.TokenRepo{},
		PasswordChecker: &util.PasswordChecker{},
		TokenUtil: &util.TokenUtil{
			TokenRepo: &repo.TokenRepo{},
		},
	}
}
