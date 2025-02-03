package mocks

import (
	"task-golang/model"

	"github.com/stretchr/testify/mock"
)

type TaskRepoMock struct {
	mock.Mock
}

func (r *TaskRepoMock) SaveTask(task *model.Task) (*model.Task, error) {
	args := r.Called(task)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TaskRepoMock) GetTaskById(id int64) (*model.Task, error) {
	args := r.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Task), args.Error(1)
	}
	return nil, args.Error(1)
}

func (r *TaskRepoMock) GetTasks(name string, priority string, boardId int64, page int, count int) (*model.TaskPageResponseDto, error) {
	args := r.Called(name, priority, boardId, page, count)
	if args.Get(0) != nil {
		return args.Get(0).(*model.TaskPageResponseDto), args.Error(1)
	}
	return nil, args.Error(1)
}
