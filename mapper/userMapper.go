package mapper

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"task-golang/model"
)

func BuildUser(ctx context.Context, dto *model.UserRegistrationDto) (*model.User, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	ps, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)

	if err != nil {
		logger.Errorf("ActionLog.BuildUser.error: cannot encrypt password for user id")
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.can't-encrypt-password", model.Exception),
			Message: err.Error(),
			Code:    http.StatusNotFound,
		}
	}

	user := &model.User{
		Email:    dto.Email,
		Password: ps,
		IsActive: false,
	}

	return user, nil
}
