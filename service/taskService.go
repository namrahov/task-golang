package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	mapper "task-golang/mapper"
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

	board, errGetBoard := ts.BoardRepo.GetBoardById(boardId)
	if errGetBoard != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-take-board", model.Exception),
			Message: errGetBoard.Error(),
			Code:    http.StatusNotFound,
		}
	}

	user, err := ts.UserUtil.GetUserFromRequest(ctx)
	if err != nil {
		return err
	}

	var task = mapper.BuildTask(dto, user, board)
	_, errSaveTask := ts.TaskRepo.SaveTask(task)
	if errSaveTask != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-save-task", model.Exception),
			Message: errSaveTask.Error(),
			Code:    http.StatusNotFound,
		}
	}
	logger.Info("ActionLog.CreateTask.end")
	return nil
}
