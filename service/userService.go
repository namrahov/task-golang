package service

import (
	"context"
	"net/http"
	"task-golang/model"
	"task-golang/repo"
	"task-golang/util"
)

type IUserService interface {
	Register(ctx context.Context, userRegistrationDto *model.UserRegistrationDto) *model.ErrorResponse
}

type UserService struct {
	UserRepo        repo.IUserRepo
	PasswordChecker util.IPasswordChecker
}

func (us *UserService) Register(ctx context.Context, dto *model.UserRegistrationDto) *model.ErrorResponse {

	if !us.PasswordChecker.IsMiddleStrength(dto.Password) {
		return &model.ErrorResponse{
			Error:   "PASSWORD_CHECK_EXCEPTION",
			Message: "PASSWORD_SHOULD_HAS_MIN_8_SYMBOL_LOWERCASE_UPPERCASE_DIGIT",
			Code:    http.StatusBadRequest,
		}
	}

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
