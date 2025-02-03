package repo

import (
	"errors"
	"gorm.io/gorm"
	"task-golang/model"
)

type ITaskRepo interface {
	SaveTask(task *model.Task) (*model.Task, error)
	GetTaskById(id int64) (*model.Task, error)
	GetTasks(name string, priority string, boardId int64, page int, count int) (*model.TaskPageResponseDto, error)
}

type TaskRepo struct {
}

func (r TaskRepo) GetTasks(name string, priority string, boardId int64, page int, count int) (*model.TaskPageResponseDto, error) {
	var tasks []model.Task
	var totalCount int64

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if count < 1 {
		count = 10
	}

	// Initialize query
	query := Db.Model(&model.Task{})

	// Apply filters dynamically
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if priority != "" {
		query = query.Where("priority = ?", priority)
	}
	if boardId > 0 {
		query = query.Where("board_id = ?", boardId)
	}

	// Count total records after filtering
	query.Count(&totalCount)

	// Apply pagination
	offset := (page - 1) * count
	query = query.Offset(offset).Limit(count)

	// Execute query
	if err := query.Preload("CreatedBy").Preload("ChangedBy").
		Preload("AssignedBy").Preload("AssignedTo").
		Preload("Board").
		Find(&tasks).Error; err != nil {
		return &model.TaskPageResponseDto{}, err
	}

	// Convert Task to TaskResponseDto
	var taskDtos []*model.TaskResponseDto
	for _, task := range tasks {
		taskDtos = append(taskDtos, &model.TaskResponseDto{
			Id:                task.Id,
			Name:              task.Name,
			Priority:          task.Priority,
			Status:            task.Status,
			CreatedBy:         task.CreatedBy,
			ChangedBy:         task.ChangedBy,
			AssignedBy:        task.AssignedBy,
			AssignedTo:        task.AssignedTo,
			Board:             task.Board,
			Deadline:          task.Deadline,
			AttachmentFileIds: []int64{}, // Populate accordingly
			TaskVideoId:       []int64{}, // Populate accordingly
			TaskImageUrl:      "",        // Populate accordingly
		})
	}

	// Calculate last page number
	lastPageNumber := (totalCount + int64(count) - 1) / int64(count)

	// Construct response
	response := model.TaskPageResponseDto{
		Tasks:          taskDtos,
		HasNextPage:    int64(page) < lastPageNumber,
		LastPageNumber: lastPageNumber,
		TotalCount:     totalCount,
	}

	return &response, nil
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
