package mapper

import (
	"task-golang/model"
	"time"
)

func BuildTask(dto *model.TaskRequestDto, user *model.User, board *model.Board) *model.Task {

	task := &model.Task{
		Name:       dto.Name,
		Priority:   dto.Priority,
		Status:     model.NotStarted,
		CreatedBy:  user,
		AssignedBy: user,
		Board:      board,
		CreatedAt:  time.Now(),
	}

	return task
}
