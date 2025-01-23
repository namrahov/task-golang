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
	activationToken := util.GenerateToken()

	if user == nil {
		buildUser, errBuildUser := mapper.BuildUser(ctx, dto)

		if errBuildUser != nil {
			logger.Errorf("ActionLog.Register.error: cannot build user")
			return errBuildUser
		}

		savedUser, errSaveUser := us.UserRepo.SaveUser(buildUser)

		if errSaveUser != nil {
			return &model.ErrorResponse{
				Error:   fmt.Sprintf("%s.can't-save-user", model.Exception),
				Message: errSaveUser.Error(),
				Code:    http.StatusForbidden,
			}
		}
		errAddUserRole := us.UserRepo.AddRolesToUser(savedUser.Id, []*model.Role{
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

	}

	logger.Info("ActionLog.Register.success")
	/*
	  if (!PasswordStrengthChecker.isMiddleStrength(dto.getPassword())) {
	            throw new UnavailableException("PASSWORD_SHOULD_HAS_MIN_8_SYMBOL_LOWERCASE_UPPERCASE_DIGIT");
	        }

	        List<User> userEntities = userUtil.findUserByEmail(dto.getEmail());
	        String activationToken = tokenUtil.generateToken();
	        User user;

	        if (userEntities.isEmpty()) {
	            user = userRepository.save(userUtil.buildUser(dto, false));

	            tokenService.saveToken(tokenMapper.toToken(activationToken, user.getId()));

	            EmailDto emailDto = emailUtil.generateActivationEmail(activationToken, REGISTRATION);
	            emailUtil.send(emailDto.getFrom(), dto.getEmail(), emailDto.getSubject(), emailDto.getBody());
	        } else {
	            user = userEntities.getFirst();

	            if (Boolean.TRUE.equals(user.getIsActive())) throw new UserRegisterException("USER_ALREADY_EXIST");
	            if (user.getInactivatedDate() != null) throw new UserRegisterException("USER_IS_INACTIVE");

	            tokenService.reSetActivationToken(user, activationToken);

	            EmailDto emailDto = emailUtil.generateActivationEmail(activationToken, REGISTRATION);
	            emailUtil.send(emailDto.getFrom(), dto.getEmail(), emailDto.getSubject(), emailDto.getBody());

	            throw new UserRegisterException("ACTIVATION_EMAIL_HAS_SENT");
	        }
	*/
	return nil
}
