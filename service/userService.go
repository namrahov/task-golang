package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"task-golang/config"
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
	Logout(ctx context.Context) *model.ErrorResponse
	CheckPermission(roles []string, requestURI, httpMethod string) bool
	ExistByToken(ctx context.Context, token string) bool
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

	// Begin GORM transaction
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
	for _, userRole := range user.Roles {
		fmt.Println(userRole)
	}
	claims := jwt.MapClaims{
		"user_id": user.Id,
		"roles":   user.Roles,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	signedToken, err := token.SignedString([]byte(config.Props.JwtSecret))
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

	errSaveToken := us.TokenRepo.SaveToken(ctx,
		mapper.BuildToken(signedToken,
			user,
			dto.RememberMe,
			stringToInt64(config.Props.TokenLifetime),
			stringToInt64(config.Props.TokenExtendedLifetime)))
	if errSaveToken != nil {
		return nil, nil
	}

	logger.Info("ActionLog.Authenticate.end")
	return jwtToken, nil
}

// ActivateIfInactiveLess30Days activates a user if inactive for less than 30 days
func (us *UserService) activateIfInactiveLess30Days(user *model.User) error {
	// Check if the user is inactive
	if !user.IsActive {
		if user.InactivatedDate != nil {
			// Check if the inactivated date is within the last 30 days
			if user.InactivatedDate.After(time.Now().AddDate(0, 0, -30)) {
				user.IsActive = true
				user.InactivatedDate = nil

				// Update the user in the repository
				_, err := us.UserRepo.UpdateUser(user)
				if err != nil {
					return fmt.Errorf("failed to update user: %w", err)
				}
				return nil
			} else {
				return errors.New("USER_IS_INACTIVE: inactive for more than 30 days")
			}
		} else {
			return errors.New("USER_IS_INACTIVE: inactivation date is missing")
		}
	}

	// User is already active
	return nil
}

func (us *UserService) Logout(ctx context.Context) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.Logout.start")

	// Retrieve the auth header from the context
	authHeader, ok := ctx.Value(model.ContextAuthHeader).(string)
	if !ok {
		logger.Warn("ActionLog.Logout.missingAuthHeader")
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.authorization_header_is_missing", model.Exception),
			Message: "Authorization header is missing",
			Code:    http.StatusNotFound,
		}
	}

	token, errFindTokenByToken := us.TokenRepo.FindTokenByToken(ctx, authHeader)

	if errFindTokenByToken != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant_find_token", model.Exception),
			Message: "Token is not found",
			Code:    http.StatusNotFound,
		}
	}

	errDeleteToken := us.TokenRepo.DeleteToken(ctx, token)
	if errDeleteToken != nil {
		return &model.ErrorResponse{
			Error:   fmt.Sprintf("%s.cant_delete_token", model.Exception),
			Message: "Token is not deleted",
			Code:    http.StatusForbidden,
		}
	}

	logger.Info("ActionLog.Logout.end")
	return nil
}

func (us *UserService) CheckPermission(roles []string, requestURI, httpMethod string) bool {
	if len(roles) == 0 {
		log.Errorf("Roles slice is empty, cannot check permissions.")
		return false
	}

	permissions, err := us.UserRepo.GetPermissions(roles)
	if err != nil {
		return false
	}

	for _, permission := range permissions {
		if matchPattern(permission.URL, requestURI) && strings.EqualFold(permission.HTTPMethod, httpMethod) {
			return true
		}
	}

	return false
}

func (us *UserService) ExistByToken(ctx context.Context, token string) bool {
	isExist := us.TokenRepo.ExistByToken(ctx, token)
	return isExist
}

func stringToInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// Handle the error, e.g., log it or return it
		log.Fatalf("Failed to parse TokenLifetime: %v", err)
	}

	return i
}

// matchPattern checks if a request URI matches a permission pattern
func matchPattern(pattern, requestURI string) bool {
	// Remove query parameters (everything after '?')
	cleanedRequestURI := strings.Split(requestURI, "?")[0]

	fmt.Println("requestURI=", cleanedRequestURI)
	// Replace all placeholders in the pattern with a generic regex for non-slash values
	regexPattern := "^" + regexp.MustCompile(`\{[^/}]+\}`).ReplaceAllString(pattern, `[^/]+`) + "$"

	// Log the generated regex pattern for debugging
	log.Printf("Generated regex pattern: %s", regexPattern)
	log.Printf("Cleaned request URI: %s", cleanedRequestURI)

	// Match the cleaned request URI against the compiled regex
	matched, err := regexp.MatchString(regexPattern, cleanedRequestURI)
	if err != nil {
		log.Printf("Error matching pattern: %v", err)
		return false
	}

	return matched
}
