package repo

import (
	"errors"
	"gorm.io/gorm"
	"task-golang/model"
)

type ITaskRepo interface {
	SaveTask(task *model.Task) (*model.Task, error)
	GetTaskById(id int64) (*model.Task, error)
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

func (r TaskRepo) GetTaskById(id int64) (*model.Task, error) {
	var task model.Task
	err := Db.
		Preload("CreatedBy").
		Preload("ChangedBy").
		Preload("AssignedBy").
		Preload("AssignedTo").
		Preload("Board").
		First(&task, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Return nil user and no error if the record is not found
		return nil, nil
	}

	if err != nil {
		// Return any other error
		return nil, err
	}

	return &task, nil
}
