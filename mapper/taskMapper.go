package mapper

import "task-golang/model"

func BuildTask(dto *model.TaskRequestDto, user *model.User) model.Task {

	task := model.Task{
		Name:       dto.Name,
		Priority:   dto.Priority,
		Status:     model.NotStarted,
		CreatedBy:  user,
		AssignedBy: user,
	}
}
