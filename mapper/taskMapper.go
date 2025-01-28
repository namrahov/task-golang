package mapper

import (
	"task-golang/model"
	"time"
)

func BuildTask(dto *model.TaskRequestDto, user *model.User, board *model.Board) *model.Task {
	task := &model.Task{
		Name:         dto.Name,
		Priority:     dto.Priority,
		Status:       model.NotStarted,
		CreatedByID:  &user.Id,  // Setting the foreign key
		AssignedByID: &user.Id,  // e.g. assigned by same user
		BoardID:      &board.Id, // referencing the board
		CreatedAt:    time.Now(),
	}
	return task
}
