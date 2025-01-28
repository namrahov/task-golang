package service

import (
	"task-golang/repo"
)

type ITaskService interface {
}

type TaskService struct {
	TaskRepo repo.ITaskRepo
}
