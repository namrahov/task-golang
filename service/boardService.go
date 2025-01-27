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

	tx := repo.BeginTransaction()
	if tx.Error != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.transaction-begin-failed", model.Exception),
			Message: tx.Error.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	// Handle transaction rollback/commit with deferred function
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback() // Rollback on panic
			panic(p)
		} else if tx.Error != nil {
			_ = tx.Rollback() // Rollback on error
		} else {
			err := tx.Commit() // Commit if no error
			if err != nil {
				//logger.WithError(err).Error("Transaction commit failed")
				fmt.Println("Transaction commit failed")
			}
		}
	}()

	_, err := bs.BoardRepo.SaveBoard(tx, mapper.BuildBoard(dto.Name, user.UserName))

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
