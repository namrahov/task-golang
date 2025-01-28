package service

import (
	"context"
	log "github.com/sirupsen/logrus"
	"task-golang/mapper"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
)

type ITaskService interface {
	CreateTask(ctx context.Context, dto *model.TaskRequestDto, boardId int64) *model.ErrorResponse
}

type TaskService struct {
	TaskRepo  repo.ITaskRepo
	BoardRepo repo.IBoardRepo
	UserUtil  util.IUserUtil
}

func (ts *TaskService) CreateTask(ctx context.Context, dto *model.TaskRequestDto, boardId int64) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.CreateTask.start")

	user, err := ts.UserUtil.GetUserFromRequest(ctx)
	if err != nil {
		return err
	}

	mapper.BuildTask(dto, user)
	logger.Info("ActionLog.CreateTask.end")
	return nil
}
