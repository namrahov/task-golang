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
	GetTask(ctx context.Context, id int64) (*model.TaskResponseDto, *model.ErrorResponse)
}

type TaskService struct {
	TaskRepo  repo.ITaskRepo
	BoardRepo repo.IBoardRepo
	FileRepo  repo.IFileRepo
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

func (ts *TaskService) GetTask(ctx context.Context, id int64) (*model.TaskResponseDto, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetTask.start")

	task, err := ts.TaskRepo.GetTaskById(id)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-get-task", model.Exception),
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}

	taskAttachmentFiles, errFindTaskAttachments := ts.FileRepo.FindTaskAttachmentsFileByTaskId(id)
	if errFindTaskAttachments != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-get-task", model.Exception),
			Message: errFindTaskAttachments.Error(),
			Code:    http.StatusNotFound,
		}
	}
	taskTaskVideos, errFindTaskTask := ts.FileRepo.FindTaskTaskVideosByTaskId(id)
	if errFindTaskTask != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-get-task", model.Exception),
			Message: errFindTaskTask.Error(),
			Code:    http.StatusNotFound,
		}
	}

	taskResponseDto := mapper.BuildTaskResponse(task, taskAttachmentFiles, taskTaskVideos)

	logger.Info("ActionLog.GetTask.end")
	return taskResponseDto, nil
}
