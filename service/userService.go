package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"task-golang/mapper"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
	"time"
)

type IUserService interface {
	Register(ctx context.Context, userRegistrationDto *model.UserRegistrationDto) *model.ErrorResponse
	Active(ctx context.Context, token string) *model.ErrorResponse
	Authenticate(ctx context.Context, dto *model.AuthRequestDto) (*model.JwtToken, *model.ErrorResponse)
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

		if user.InactivatedDate != nil {
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

func (us *UserService) Active(ctx context.Context, token string) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.Active.start")

	existingToken, err := us.TokenRepo.FindTokenByActivationToken(ctx, token)
	if err != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.can_not_get_token", model.Exception),
			Message: "Can not get token",
			Code:    http.StatusForbidden,
		}
	}

	user, err := us.UserRepo.FindUserById(existingToken.UserID)
	if user == nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_not_found", model.Exception),
			Message: "Can not find user",
			Code:    http.StatusForbidden,
		}
	}

	user.IsActive = true
	_, errSaveUser := us.UserRepo.UpdateUser(user)
	if errSaveUser != nil {
		return nil
	}

	errDeleteToken := us.TokenRepo.DeleteToken(ctx, existingToken)
	if errDeleteToken != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_not_found", model.Exception),
			Message: "Can not find user",
			Code:    http.StatusForbidden,
		}
	}

	logger.Info("ActionLog.Active.end")
	return nil
}

func (us *UserService) Authenticate(ctx context.Context, dto *model.AuthRequestDto) (*model.JwtToken, *model.ErrorResponse) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.Authenticate.start")

	// Find user by email or username
	user, errFindUser := us.UserRepo.FindActiveUserByEmailOrUsername(dto.EmailOrNickname)
	if errFindUser != nil || user == nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.user_not_found", model.Exception),
			Message: "Cannot find user",
			Code:    http.StatusForbidden,
		}
	}

	// Check password
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(dto.Password))
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.invalid_credentials", model.Exception),
			Message: "Invalid email/username or password",
			Code:    http.StatusForbidden,
		}
	}

	// Activate user if inactive less than 30 days
	err = us.activateIfInactiveLess30Days(user)
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.activation_failed", model.Exception),
			Message: "Failed to activate user",
			Code:    http.StatusInternalServerError,
		}
	}

	// Generate JWT Token
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"roles":   user.Roles,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	secretKey := "your_secret_key_here" // Replace with your secure key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.token_generation_failed", model.Exception),
			Message: "Failed to generate token",
			Code:    http.StatusInternalServerError,
		}
	}

	// Prepare the response
	jwtToken := &model.JwtToken{
		Token: signedToken,
	}

	logger.Info("ActionLog.Authenticate.end")
	return jwtToken, nil
}

// ActivateIfInactiveLess30Days activates a user if inactive for less than 30 days
func (us *UserService) activateIfInactiveLess30Days(user *model.User) error {
	if !user.IsActive {
		if user.InactivatedDate != nil {
			// Parse the inactivated date
			inactivatedTime, err := time.Parse("2006-01-02", *user.InactivatedDate)
			if err != nil {
				return errors.New("invalid inactivated date format")
			}

			// Check if the inactivated date is within the last 30 days
			if inactivatedTime.After(time.Now().AddDate(0, 0, -30)) {
				user.IsActive = true
				user.InactivatedDate = nil
				_, err := us.UserRepo.UpdateUser(user)
				if err != nil {
					return err
				}
				return err
			} else {
				return errors.New("USER_IS_INACTIVE")
			}
		} else {
			return errors.New("USER_IS_INACTIVE")
		}
	}
	return nil
}
