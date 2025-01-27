package service

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"task-golang/mapper"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
)

type IBoardService interface {
	CreateBoard(ctx context.Context, userRegistrationDto *model.BoardRequestDto) *model.ErrorResponse
	GiveAccessToBoard(ctx context.Context, userId int64, boardId int64) *model.ErrorResponse
}

type BoardService struct {
	BoardRepo repo.IBoardRepo
	UserRepo  repo.IUserRepo
	UserUtil  util.IUserUtil
}

func (bs *BoardService) CreateBoard(ctx context.Context, dto *model.BoardRequestDto) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.CreateBoard.start")

	user, errGetUser := bs.UserUtil.GetUserFromRequest(ctx)
	if errGetUser != nil {
		return errGetUser
	}

	_, err := bs.BoardRepo.SaveBoard(mapper.BuildBoard(dto.Name, user.UserName))

	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-take-board", model.Exception),
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}

	logger.Info("ActionLog.CreateBoard.end")
	return nil
}

func (bs *BoardService) GiveAccessToBoard(ctx context.Context, userId int64, boardId int64) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GiveAccessToBoard.start")

	err := bs.BoardRepo.SaveUserBoard(ctx, userId, boardId)
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant-save-user-board", model.Exception),
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}

	logger.Info("ActionLog.GiveAccessToBoard.end")
	return nil
}
