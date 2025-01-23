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

type IUserService interface {
	Register(ctx context.Context, userRegistrationDto *model.UserRegistrationDto) *model.ErrorResponse
}

type UserService struct {
	UserRepo        repo.IUserRepo
	TokenRepo       repo.ITokenRepo
	PasswordChecker util.IPasswordChecker
	TokenUtil       util.ITokenUtil
}

func (us *UserService) Register(ctx context.Context, dto *model.UserRegistrationDto) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.Register.start")

	if !us.PasswordChecker.IsMiddleStrength(dto.Password) {
		logger.Errorf("ActionLog.Register.error: password is weak")
		return &model.ErrorResponse{
			Error:   "PASSWORD_CHECK_EXCEPTION",
			Message: "PASSWORD_SHOULD_HAS_MIN_8_SYMBOL_LOWERCASE_UPPERCASE_DIGIT",
			Code:    http.StatusBadRequest,
		}
	}

	user, errGetUser := us.UserRepo.GetUserByEmail(dto.Email)
	if errGetUser != nil {
		logger.Errorf("ActionLog.Register.error: cannot get user by email %v", dto.Email)
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.can't-get-user", model.Exception),
			Message: errGetUser.Error(),
			Code:    http.StatusNotFound,
		}
	}
	activationToken := us.TokenUtil.GenerateToken()

	tx, errBeginTransaction := us.UserRepo.BeginTransaction()
	if errBeginTransaction != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.transaction-begin-failed", model.Exception),
			Message: errBeginTransaction.Error(),
			Code:    http.StatusInternalServerError,
		}
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback() // Rollback on panic
			panic(p)
		} else if errBeginTransaction != nil {
			tx.Rollback() // Rollback on error
		} else {
			errBeginTransaction = tx.Commit() // Commit if no error
		}
	}()

	if user == nil {
		buildUser, errBuildUser := mapper.BuildUser(ctx, dto)

		if errBuildUser != nil {
			logger.Errorf("ActionLog.Register.error: cannot build user")
			return errBuildUser
		}

		savedUser, errSaveUser := us.UserRepo.SaveUser(tx, buildUser)

		if errSaveUser != nil {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.can't-save-user", model.Exception),
				Message: errSaveUser.Error(),
				Code:    http.StatusForbidden,
			}
		}
		errAddUserRole := us.UserRepo.AddRolesToUser(tx, savedUser.Id, []*model.Role{
			{
				Id:   1,
				Name: "user",
			},
		})

		if errAddUserRole != nil {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.can't-add-user-role", model.Exception),
				Message: errAddUserRole.Error(),
				Code:    http.StatusForbidden,
			}
		}

		errSaveToken := us.TokenRepo.SaveToken(ctx, mapper.BuildActivationToken(activationToken, savedUser.Id))
		if errSaveToken != nil {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.can't-save-activation-token", model.Exception),
				Message: errSaveToken.Error(),
				Code:    http.StatusForbidden,
			}
		}

		emailDto := util.GenerateActivationEmail(activationToken, model.Registration)
		util.SendEmailAsync(emailDto.From, dto.Email, emailDto.Subject, emailDto.Body)
	} else {
		if user.IsActive == true {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.user_exist", model.Exception),
				Message: "User exist",
				Code:    http.StatusForbidden,
			}
		}

		if user.InactivatedDate != "" {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.user_is_inactive", model.Exception),
				Message: "User is inactive",
				Code:    http.StatusForbidden,
			}
		}

		us.TokenUtil.ReSetActivationToken(ctx, user, activationToken)

		emailDto := util.GenerateActivationEmail(activationToken, model.Registration)
		util.SendEmailAsync(emailDto.From, dto.Email, emailDto.Subject, emailDto.Body)

		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_is_inactive", model.Exception),
			Message: "User is inactive",
			Code:    http.StatusForbidden,
		}
	}

	logger.Info("ActionLog.Register.success")
	return nil
}
