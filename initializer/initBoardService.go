package initializer

import (
	"task-golang/repo"
	"task-golang/service"
	"task-golang/util"
)

func InitBoardService() *service.BoardService {
	return &service.BoardService{
		BoardRepo: &repo.BoardRepo{},
		UserRepo:  &repo.UserRepo{},
		UserUtil: &util.UserUtil{
			UserRepo: &repo.UserRepo{},
		},
	}
}
