package repo

import "task-golang/model"

type ITaskRepo interface {
	SaveTask(task *model.Task) (*model.Task, error)
}

type TaskRepo struct {
}

func (r TaskRepo) SaveTask(task *model.Task) (*model.Task, error) {
	result := Db.Create(task)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil
}
