package initializer

import (
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

func InitFileService() *service.FileService {
	return &service.FileService{
		FileRepo: &repo.FileRepo{},
		UserRepo: &repo.UserRepo{},
		UserUtil: &util.UserUtil{
			UserRepo: &repo.UserRepo{},
		},
	}
}
