package util

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"task-golang/model"
	"task-golang/repo"
)

type IUserUtil interface {
	GetUserFromRequest(ctx context.Context) (*model.User, *model.ErrorResponse)
}

type UserUtil struct {
	UserRepo repo.IUserRepo
}

func (u *UserUtil) GetUserFromRequest(ctx context.Context) (*model.User, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	userIdString, ok := ctx.Value(model.ContextUserID).(string)
	if !ok {
		logger.Warn("ActionLog.CreateBoard.missingUserId")
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_id_is_missing", model.Exception),
			Message: "User id is missing",
			Code:    http.StatusNotFound,
		}
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_id_can_not_be_parsed", model.Exception),
			Message: "Can parse string to int",
			Code:    http.StatusNotFound,
		}
	}

	user, err := u.UserRepo.FindUserById(userId)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.can_not_get_user", model.Exception),
			Message: "Can get user",
			Code:    http.StatusNotFound,
		}
	}

	return user, nil
}
