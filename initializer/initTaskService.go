package initializer

import (
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

func InitTaskService() *service.TaskService {
	return &service.TaskService{
		TaskRepo:  &repo.TaskRepo{},
		BoardRepo: &repo.BoardRepo{},
		FileRepo:  &repo.FileRepo{},
		UserUtil: &util.UserUtil{
			UserRepo: &repo.UserRepo{},
		},
	}
}
