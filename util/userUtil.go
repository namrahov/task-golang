package util

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
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
	userId, ok := ctx.Value(model.ContextUserID).(int64)
	if !ok {
		logger.Info("userIdString:", userId)
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
