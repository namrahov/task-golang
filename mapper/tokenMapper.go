package mapper

import (
	"github.com/google/uuid"
	"task-golang/model"
	"time"
)

func BuildActivationToken(activationToken string, userId int64) *model.Token {
	token := &model.Token{
		Token:           uuid.New().String(), // Generate a unique token using UUID
		ActivationToken: activationToken,
		UserID:          userId,
		CreatedAt:       time.Now(),
		TTL:             3600, // 1 hour
	}

	return token
}

// BuildToken creates a Token based on the provided parameters.
func BuildToken(token string, user *model.User, rememberMe bool, tokenLifeTime, tokenExtendedLifeTime int64) *model.Token {
	ttl := tokenLifeTime * 60
	if rememberMe {
		ttl = tokenExtendedLifeTime * 60
	}

	return &model.Token{
		Token:     "Bearer " + token,
		UserID:    user.Id,
		CreatedAt: time.Now(),
		TTL:       ttl,
	}
}
